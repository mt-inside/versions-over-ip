package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

type ghrelease struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Prerelease  bool   `json:"prerelease"`
}

func Github(org, repo string, depth, count int32) ([]Series, error) {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", org, repo)

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "versions")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := []ghrelease{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	ss := parseGithub(data, depth, count)
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
func parseGithub(rs []ghrelease, depth, count int32) []Series {
	ss := make([]Series, count)

	for _, r := range rs {
		v, _ := version.NewVersion(r.TagName)

		inject(ss, v, r.Prerelease, depth)
	}

	return ss
}

/* Mutates in place to avoid hammering the allocator.
* Building a new structure would honestly be easier */
func inject(ss []Series, v *version.Version, pre bool, depth int32) {
	s := truncate(v, depth)

	var i int
	inj := false
	for i = 0; i < len(ss); i++ {
		if ss[i].Prefix == nil || s.GreaterThan(ss[i].Prefix) {
			inj = true
			break
		}
		if s.Equal(ss[i].Prefix) {
			trump(&ss[i], v, pre)
			break
		}
	}
	if inj {
		for j := len(ss) - 1; j > i; j-- {
			ss[j] = ss[j-1]
		}
		if pre {
			ss[i] = Series{s, nil, v}
		} else {
			ss[i] = Series{s, v, nil}
		}
	}
}

func trump(s *Series, v *version.Version, pre bool) {
	if pre && (s.Prerelease == nil || v.GreaterThan(s.Prerelease)) {
		s.Prerelease = v
	} else if !pre && (s.Stable == nil || v.GreaterThan(s.Stable)) {
		s.Stable = v
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
