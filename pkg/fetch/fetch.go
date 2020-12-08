package fetch

import (
	"time"

	"github.com/hashicorp/go-version"
)

type Series struct {
	Name     string
	Releases map[string]Release
}

type Release struct {
	Version *version.Version
	Date    time.Time
}
