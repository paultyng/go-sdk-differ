package locator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func LocateLatestAzureSDK(goPath string, sdkInUse string) (*string, error) {
	// github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-03-01/compute
	paths := strings.Split(sdkInUse, "/")
	currentVersion := paths[len(paths)-2]
	sdkName := paths[len(paths)-1]

	// then strip off the version number (and the forward slash at the end)
	versionsPath := sdkInUse[0:strings.Index(sdkInUse, currentVersion)-1]
	workingDirectory := fmt.Sprintf("%s/%s", goPath, versionsPath)

	cmd := exec.Command("/bin/bash", "-c", "ls | sort | tail -1")
	cmd.Dir = workingDirectory

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	latestVersion := strings.TrimSpace(string(output))
	latestSdkPath := fmt.Sprintf("%s/%s/%s", versionsPath, latestVersion, sdkName)

	// ensure this path actually exists since RP's can have different packages within different API versions (resources)
	if _, err := os.Stat(fmt.Sprintf("%s/%s", goPath, latestSdkPath)); os.IsNotExist(err) {
		log.Printf("      SDK Upgrade available but not for this Resource Provider..")
		return &sdkInUse, nil
	}

	if latestVersion != currentVersion {
		log.Printf("      SDK Upgrade available: %s", latestSdkPath)
		return &sdkInUse, nil
	}

	return &latestSdkPath, nil
}
