package version

import (
	"fmt"
	"strings"
)

var (
	// Raw is the string representation of the version. This will be replaced
	// with the calculated version at build time.
	// set in the Makefile.
	Version = "was not built with version info"

	// Commit is the commit hash from which the software was built.
	// Set via LDFLAGS in Makefile.
	Commit = "unknown"

	// BuildTime is the string representation of build time.
	// Set via LDFLAGS in Makefile.
	BuildTime = "unknown"
)

func GetShortVersion() string {
	if strings.HasSuffix(Version, "dirty") {
		return fmt.Sprintf("Build (%s)", Commit)
	}
	return Version
}
