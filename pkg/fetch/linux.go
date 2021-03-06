package fetch

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
)

type linuxReleaseDate struct {
	IsoDate string `json:"isodate"`
}

type linuxRelease struct {
	Moniker  string           `json:"moniker"`
	Version  string           `json:"version"`
	Released linuxReleaseDate `json:"released"`
}

type linuxReleases struct {
	Releases []linuxRelease `json:"releases"`
}

func Linux() ([]Series, error) {
	uri := "https://www.kernel.org/releases.json"

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

	data := linuxReleases{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	linuxss, err := parseLinux(data.Releases)
	if err != nil {
		return nil, err
	}

	ss := []Series{}
	for _, s := range linuxss {
		ss = append(ss, s)
	}

	return ss, nil
}

func parseLinux(rs []linuxRelease) (map[string]Series, error) {
	zeroVersion, _ := version.NewVersion("0.0.0")
	ms := []string{"mainline", "stable", "longterm"}
	ss := map[string]Series{}
	for _, m := range ms {
		ss[m] = Series{m, map[string]Release{"": Release{zeroVersion, time.Unix(0, 0)}}}
	}

	for _, r := range rs {
		v, err := version.NewVersion(r.Version)
		if err != nil {
			continue // Ignore weird versions like "next-20201011". Could skip monikers (like "next") known to have them.
		}

		if v.GreaterThan(ss[r.Moniker].Releases[""].Version) {
			d, err := time.Parse("2006-01-02", r.Released.IsoDate)
			if err != nil {
				return nil, err
			}

			ss[r.Moniker].Releases[""] = Release{v, d}
		}
	}

	return ss, nil
}
