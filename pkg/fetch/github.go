package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-version"
)

type ghTag struct {
	Name string `json:"name"`
}
type ghRelease struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Prerelease  bool   `json:"prerelease"`
}

type comparibleSeries struct {
	prefix *version.Version
	series Series
}

func Github(log logr.Logger, org, repo string, tags bool, depth, count int32) ([]Series, error) {
	// TODO: if no releases, auto-fall back to tags
	var uri string
	if tags {
		//FIXME: seems the most recent are first, but because it's reverse-sorted, the old pre-release "weekly" tags are sorting before "go", so the real stuff starts at page 4.
		//What we should do is support the pagination, start at 1 and follow the Link headers until we've got $count's worth of non-pre entries (or we hit the end of the list)
		uri = fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?page=4", org, repo)
	} else {
		uri = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", org, repo)
	}

	// TODO: use gapic for this... learn how to describe this service and generate a client for it that retries. Just get byte array from that client, or give it something that'll let it JSON unmarshall?
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "versions")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := []ghRelease{}
	if tags {
		tgs := []ghTag{}
		err = json.Unmarshal(body, &tgs)
		if err != nil {
			return nil, err
		}
		for _, t := range tgs {
			log.V(2).Info("Got tag", "name", t.Name)
			data = append(data, ghRelease{t.Name, time.Now().Format(time.RFC3339), false})
		}
	} else {
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
	}

	ghss, err := parseGithub(log, data, depth, count)
	if err != nil {
		return nil, err
	}

	ss := make([]Series, len(ghss))
	for i, s := range ghss {
		ss[i] = s.series
	}
	return ss, nil
}

/* Algo:
* []release
* [](release + meaningful version)
* sort
* groupby(version[:depth])
* [][]relase
* []release[:count]
 */
/* since there's no map etc in go, do a linear one */
func parseGithub(log logr.Logger, rs []ghRelease, depth, count int32) ([]comparibleSeries, error) {
	ss := make([]comparibleSeries, count)

	for _, r := range rs {
		rFixed := &r
		log.V(2).Info("Processing tag/release", "name", rFixed.TagName)
		// FIXME: total hack. Use of nil is horrible, and makes ordering important
		// need to indicate
		// - i got this, here's the real version
		// - this is trash, but it's my trash; ignore
		// - this is trash I don't recognise, someone else have a go
		// - if its still marked as trash at the end, drop it.
		// - maybe the alternative - a few functions try to extract it (loop over array of extractors?) If none say yes, drop it. First attempt is just "does it parse as vX.Y.Z?"
		rFixed = fixupLinkerd(rFixed)
		rFixed = fixupZfs(rFixed)
		rFixed = fixupNginx(rFixed)
		rFixed = fixupGo(rFixed)
		if rFixed == nil {
			continue
		}

		v, err := version.NewVersion(rFixed.TagName)
		if err != nil {
			return nil, err
		}
		d, err := time.Parse(time.RFC3339, rFixed.PublishedAt)
		if err != nil {
			return nil, err
		}

		inject(ss, v, d, rFixed.Prerelease, depth)
	}

	var trunc bool
	var i int
	for i = range ss {
		if ss[i].prefix == nil {
			trunc = true
			break
		}
	}
	if trunc {
		ss = ss[:i]
	}

	// TODO: remove all pres that are older than GAs. Has to happen here.

	return ss, nil
}

/* Mutates in place to avoid hammering the allocator.
* Building a new structure would honestly be easier */
func inject(ss []comparibleSeries, v *version.Version, d time.Time, pre bool, depth int32) {
	s := truncate(v, depth)

	var i int
	inj := false
	for i = 0; i < len(ss); i++ {
		if ss[i].prefix == nil || s.GreaterThan(ss[i].prefix) {
			inj = true
			break
		}
		if s.Equal(ss[i].prefix) {
			trump(&ss[i].series, v, d, pre)
			break
		}
	}
	if inj {
		for j := len(ss) - 1; j > i; j-- {
			ss[j] = ss[j-1]
		}
		if pre {
			ss[i] = comparibleSeries{s, Series{s.String(), map[string]Release{"pre": Release{v, d}}}}
		} else {
			ss[i] = comparibleSeries{s, Series{s.String(), map[string]Release{"ga": Release{v, d}}}}
		}
	}
}

func fixupLinkerd(r *ghRelease) *ghRelease {
	rNew := *r

	if strings.HasPrefix(rNew.TagName, "edge-") {
		rNew.TagName = strings.TrimPrefix(rNew.TagName, "edge-")
		rNew.Prerelease = true
	} else if strings.HasPrefix(rNew.TagName, "stable-") {
		rNew.TagName = strings.TrimPrefix(rNew.TagName, "stable-")
		rNew.Prerelease = false
	}

	return &rNew
}

func fixupGo(r *ghRelease) *ghRelease {
	rNew := *r

	if strings.HasPrefix(rNew.TagName, "weekly") {
		return nil
	} else if strings.HasPrefix(rNew.TagName, "release.") {
		return nil
	} else if strings.HasPrefix(rNew.TagName, "go") {
		rNew.TagName = strings.TrimPrefix(rNew.TagName, "go")
	}

	return &rNew
}

func fixupZfs(r *ghRelease) *ghRelease {
	rNew := *r

	if strings.HasPrefix(rNew.TagName, "zfs-") {
		rNew.TagName = strings.TrimPrefix(rNew.TagName, "zfs-")
	}

	return &rNew
}

func fixupNginx(r *ghRelease) *ghRelease {
	rNew := *r

	if strings.HasPrefix(rNew.TagName, "release-") {
		rNew.TagName = "v" + strings.TrimPrefix(rNew.TagName, "release-")
		rNew.Prerelease = false
	}

	return &rNew
}

func trump(s *Series, v *version.Version, d time.Time, isPre bool) {
	pre, preFound := s.Releases["pre"]
	ga, gaFound := s.Releases["ga"]

	if isPre && (!preFound || v.GreaterThan(pre.Version)) {
		s.Releases["pre"] = Release{v, d}
	} else if !isPre && (!gaFound || v.GreaterThan(ga.Version)) {
		s.Releases["ga"] = Release{v, d}
	}
}

func truncate(v *version.Version, depth int32) *version.Version {
	/* ffs golang */
	segs := v.Segments()[:depth]
	strs := make([]string, depth)
	for i, s := range segs {
		strs[i] = strconv.Itoa(s)
	}
	t, _ := version.NewVersion(strings.Join(strs, "."))
	return t
}
