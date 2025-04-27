package buildinfo

import (
	"runtime/debug"
)

const keyRevision = "vcs.revision"
const keyTime = "vcs.time"

func GetGitRevision() (string, string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", "", false
	}
	var revision, revtime string
	for _, setting := range info.Settings {
		if setting.Key == keyRevision {
			revision = setting.Value
		}
		if setting.Key == keyTime {
			revtime = setting.Value
		}
	}
	return revision, revtime, true
}
