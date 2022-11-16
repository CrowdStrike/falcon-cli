package version

import (
	"fmt"
	"runtime"
)

var (
	Version    = "unknown"
	GitVersion = "unknown"
	GitCommit  = "unknown"
)

func VersionString() string {
	version := GitVersion
	if version == "unknown" {
		version = Version
	}

	return fmt.Sprintf("falcon version: %q, commit: %q, go version: %q, GOOS: %q, GOARCH: %q",
		version, GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
