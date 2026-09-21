package main

import (
	"fmt"
	"runtime/debug"
)

// Injected by goreleaser via -ldflags.
var (
	version = ""
	commit  = ""
	date    = ""
)

func buildVersion() string {
	if version != "" {
		return version
	}
	// Fall back to VCS info embedded by go build.
	if info, ok := debug.ReadBuildInfo(); ok {
		var rev, vcsDate string
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if len(s.Value) > 7 {
					rev = s.Value[:7]
				} else {
					rev = s.Value
				}
			case "vcs.time":
				vcsDate = s.Value
			}
		}
		if rev != "" {
			if vcsDate != "" {
				return fmt.Sprintf("%s (%s)", rev, vcsDate)
			}
			return rev
		}
	}
	return "dev"
}

func printVersion() {
	v := buildVersion()
	if commit != "" {
		v = fmt.Sprintf("%s commit=%s", v, commit)
	}
	if date != "" {
		v = fmt.Sprintf("%s date=%s", v, date)
	}
	fmt.Println(v)
}
