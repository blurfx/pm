package version

import (
	"fmt"
)

var (
	Version = ""
	Commit  = ""
)

func GetVersion() string {
	if Version == "" {
		Version = "unknown"
	}
	if Commit == "" {
		Commit = "unknown"
	}
	return fmt.Sprintf("%s (%s)", Version, Commit)
}
