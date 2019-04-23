package locator

import (
	"fmt"
	"os/exec"
	"strings"
)

func LocateImportPaths(path string, sdkImportPrefix string) (*[]string, error) {
	// grep -R "sdkImportPrefix" .
	cmd := exec.Command("/usr/bin/grep", "-R", sdkImportPrefix, ".")
	cmd.Dir = path

	str, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error reading output of grep %q: %s", path, err)
	}

	importPaths := make(map[string]struct{}, 0)

	output := string(str)
	for _, line := range strings.Split(output, "\n") {
		split := strings.Split(line, "\t")
		if len(split) < 2 {
			continue
		}

		importPathWithoutQuotes := strings.Split(split[1], "\"")
		if len(importPathWithoutQuotes) != 3 {
			continue
		}

		importPaths[importPathWithoutQuotes[1]] = struct{}{}
	}

	paths := make([]string, 0)
	for k := range importPaths {
		paths = append(paths, k)
	}

	return &paths, nil
}
