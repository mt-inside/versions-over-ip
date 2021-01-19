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

	"github.com/hashicorp/go-version"
)

type ghRelease struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Prerelease  bool   `json:"prerelease"`
}

type comparibleSeries struct {
	prefix *version.Version
	series Series
}

func Github(org, repo string, depth, count int32) ([]Series, error) {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", org, repo)

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
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	ghss, err := parseGithub(data, depth, count)
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
func parseGithub(rs []ghRelease, depth, count int32) ([]comparibleSeries, error) {
	ss := make([]comparibleSeries, count)

	for _, r := range rs {
		v, err := version.NewVersion(r.TagName)
		if err != nil {
			return nil, err
		}
		d, err := time.Parse(time.RFC3339, r.PublishedAt)
		if err != nil {
			return nil, err
		}

		inject(ss, v, d, r.Prerelease, depth)
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
