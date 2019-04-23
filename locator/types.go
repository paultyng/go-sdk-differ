package locator

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func LocateUsagesOfSDK(path string, importPath string) (*[]string, error) {
	// the type name is the last bit, for example:
	importTypeNames := strings.Split(importPath, "/")
	importType := importTypeNames[len(importTypeNames)-1]

	// first locate all of the files with this import
	fileNamesContainingImports := make(map[string]struct{}, 0)
	cmd := exec.Command("/usr/bin/grep", "-R", importPath, ".")
	cmd.Dir = path

	str, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error reading output of first grep %q: %s", path, err)
	}

	output := string(str)
	for _, line := range strings.Split(output, "\n") {
		// example:
		// ./azurerm/resource_arm_eventhub_namespace_authorization_rule.go:	parameters := eventhub.AuthorizationRule{
		split := strings.Split(line, "\t")
		if len(split) < 2 {
			continue
		}

		fileName := strings.TrimSpace(split[0])
		fileNamesContainingImports[fileName] = struct{}{}
	}

	// the format's important here since we only care about assignments
	// grep -R ":= importPath." .
	cmd = exec.Command("/usr/bin/grep", "-R", fmt.Sprintf(":= %s.", importType), ".")
	cmd.Dir = path

	str, err = cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// no files match this prefix
				// for example, if an SDK's been imported but not used
				if status.ExitStatus() == 1 {
					return nil, nil
				}
			}

		}
		return nil, fmt.Errorf("Error reading output of second grep %q: %s", path, err)
	}

	types := make(map[string]struct{}, 0)

	output = string(str)
	for _, line := range strings.Split(output, "\n") {
		// example:
		// ./azurerm/resource_arm_eventhub_namespace_authorization_rule.go:	parameters := eventhub.AuthorizationRule{
		split := strings.Split(line, "\t")
		if len(split) < 2 {
			continue
		}

		// this match is only required if it's also got the same import
		fileName := strings.TrimSpace(split[0])
		if _, exists := fileNamesContainingImports[fileName]; !exists {
			continue
		}

		code := strings.TrimSpace(split[1])
		if code == "" {
			continue
		}

		// if this ends in a { or a {} then it's a type and we can do some more gnarly string hacking
		if strings.HasSuffix(code, "{") || strings.HasSuffix(code, "{}") {
			importStr := fmt.Sprintf("%s.", importType)
			indexOfType := strings.Index(code, importStr)
			startOfType := indexOfType + len(importStr)
			indexOfBracket := strings.Index(code, "{")
			name := code[startOfType:indexOfBracket]

			// handles
			// if props := resp.Props; props != nil {
			if strings.Contains(name, " ") {
				continue
			}

			types[name] = struct{}{}
			continue
		}

		// TODO: implement me
		/*
		// otherwise we're making a type e.g.
		// outputs := make([]apimanagement.ParameterContract, 0)
		if strings.Contains(code, "make([]") {
			importStr := fmt.Sprintf("make([]%s.", importType)
			indexOfType := strings.Index(code, importStr)
			startOfType := indexOfType + len(importStr)
			indexOfBracket := strings.Index(code, ")")
			name := code[startOfType:indexOfBracket]
			types[name] = struct{}{}
			continue
		}
		 */

		continue

	}

	paths := make([]string, 0)
	for k := range types {
		//typeName := fmt.Sprintf("%s.%s", importType, k)
		paths = append(paths, k)
	}

	return &paths, nil
}
