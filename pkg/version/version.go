package version

import (
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"
)

var changelogVersionRegex = regexp.MustCompile(`(?m)^##\s+\[([0-9]+\.[0-9]+\.[0-9]+[^\]]*)\]`)

// Get automatically resolves the CrashBench version dynamically without static hardcoding.
// Order of resolution:
// 1. Exact git release tag via 'git describe --tags --exact-match'
// 2. Annotated git tag via 'git describe --tags --always'
// 3. Latest semantic release header in CHANGELOG.md (e.g. ## [1.1.0])
// 4. Go runtime VCS revision from debug.ReadBuildInfo()
// 5. Default fallback to v1.1.0
func Get() string {
	// 1. Exact git tag
	if out, err := exec.Command("git", "describe", "--tags", "--exact-match").Output(); err == nil {
		v := strings.TrimSpace(string(out))
		if v != "" {
			if !strings.HasPrefix(v, "v") {
				v = "v" + v
			}
			return v
		}
	}

	// 2. Git describe with tags
	if out, err := exec.Command("git", "describe", "--tags", "--always").Output(); err == nil {
		v := strings.TrimSpace(string(out))
		if strings.HasPrefix(v, "v") {
			return v
		}
	}

	// 3. Read latest version from CHANGELOG.md
	if v := readChangelogVersion(); v != "" {
		return v
	}

	// 4. Go runtime build info VCS commit hash
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" && !strings.HasPrefix(info.Main.Version, "v0.0.0-") {
			return info.Main.Version
		}
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return "v1.1.0-dev (" + s.Value[:7] + ")"
			}
		}
	}

	return "v1.1.0"
}

func readChangelogVersion() string {
	paths := []string{"CHANGELOG.md", "../CHANGELOG.md", "crashbench/CHANGELOG.md"}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			matches := changelogVersionRegex.FindSubmatch(data)
			if len(matches) > 1 {
				v := string(matches[1])
				if !strings.HasPrefix(v, "v") {
					v = "v" + v
				}
				return v
			}
		}
	}
	return ""
}
