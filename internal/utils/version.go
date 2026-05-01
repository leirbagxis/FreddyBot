package utils

import (
	"runtime/debug"
)

var (
	// Version can be set at build time using -ldflags:
	// -ldflags "-X 'github.com/leirbagxis/FreddyBot/internal/utils.Version=$(git rev-parse --short HEAD)'"
	Version = "dev"
)

func init() {
	if Version == "dev" || Version == "unknown" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					if len(setting.Value) > 7 {
						Version = setting.Value[:7]
					} else {
						Version = setting.Value
					}

					// Optional: add -dirty suffix if there are uncommitted changes
					for _, s := range info.Settings {
						if s.Key == "vcs.modified" && s.Value == "true" {
							Version += "-dirty"
							break
						}
					}
					return
				}
			}
		}
	}
}
