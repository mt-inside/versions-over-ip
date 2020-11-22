package fetch

import "github.com/hashicorp/go-version"

type Series struct {
	Prefix     *version.Version
	Stable     *version.Version
	Prerelease *version.Version
}
