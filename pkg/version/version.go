package version

import (
	"fmt"
	"runtime/debug"
)

type Version struct {
	ModuleVersion string
	CommitTime    string
	GitCommit     string
	GoVersion     string
	Dependencies  []*debug.Module
}

func (v Version) String() string {
	return fmt.Sprintf("module: %s, go: %s, git: %s", v.ModuleVersion, v.GoVersion, v.GitCommit)
}

func GetVersion() *Version {

	v := &Version{ModuleVersion: "unknown", CommitTime: "unknown", GitCommit: "unknown", GoVersion: "unknown"}

	info, _ := debug.ReadBuildInfo()

	if info.Main.Version != "" {
		v.ModuleVersion = info.Main.Version
	}
	if info.GoVersion != "" {
		v.GoVersion = info.GoVersion
	}

	v.Dependencies = info.Deps

	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			v.GitCommit = kv.Value
		case "vcs.time":
			v.CommitTime = kv.Value
		}
	}

	return v
}
