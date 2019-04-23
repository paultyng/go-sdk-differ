package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/tombuildsstuff/go-sdk-differ/codegen"
	"github.com/tombuildsstuff/go-sdk-differ/locator"
)

// providerImportPaths contains the base path to the Services folder in the SDK for that Provider
var providerImportPaths = map[string]string{
	"azurerm": "github.com/Azure/azure-sdk-for-go/services/",
	"aws": "github.com/aws/aws-sdk-go/service/",
}

// providerCodePaths contains the directory within the provider repo where code exists
var providerCodePaths = map[string]string{
	"azurerm": "azurerm",
	"aws": "aws",
}

func main() {
	log.Printf("Launching SDK Diff Detector..")
	provider := os.Getenv("PROVIDER")
	upgrade := strings.EqualFold(os.Getenv("FORCE_UPGRADE"), "true")

	if err := run(provider, upgrade); err != nil {
		log.Printf("Panic: %s", err)
		os.Exit(1)
	}
}

func run(providerName string, upgrade bool) error {
	goSrcPath := fmt.Sprintf("%s/src", os.Getenv("GOPATH"))
	repositoryPath := fmt.Sprintf("github.com/terraform-providers/terraform-provider-%s", providerName)
	importPath, ok := providerImportPaths[providerName]
	if !ok {
		return fmt.Errorf("Import Path was not found for %q", providerName)
	}

	providerPath, ok := providerCodePaths[providerName]
	if !ok {
		return fmt.Errorf("Provider Path was not found for %q", providerName)
	}

	repositoryFullPath := fmt.Sprintf("%s/%s", goSrcPath, repositoryPath)
	codePath := fmt.Sprintf("%s/%s", repositoryFullPath, providerPath)

	// find all of the packages being used
	importPaths, err := locator.LocateImportPaths(codePath, importPath)
	if err != nil {
		panic(err)
	}

	// find the usages of each package
	for _, oldSdkPath := range *importPaths {
		log.Printf("[DEBUG] Package %q", oldSdkPath)

		// if it's internal we can't compile it
		if strings.Contains(oldSdkPath, "/internal") {
			continue
		}

		types, err := locator.LocateUsagesOfSDK(codePath, oldSdkPath)
		if err != nil {
			return err
		}

		if types == nil {
			log.Printf("[DEBUG] Found no usages of %q", oldSdkPath)
			continue
		}

		sdkTypes := make([]codegen.TerraformTypeInfo, 0)
		for _, typeName := range *types {
			// if it's not public we can't compile it
			isPublic := string(typeName[0]) == strings.ToUpper(string(typeName[0]))
			if !isPublic {
				continue
			}

			sdkTypes = append(sdkTypes, codegen.TerraformTypeInfo{
				Package: oldSdkPath,
				ImportType: typeName,
			})
		}

		log.Printf("[DEBUG] Package %q has %d types", oldSdkPath, len(sdkTypes))
		if len(sdkTypes) == 0 {
			continue
		}

		newSdkPath := oldSdkPath
		if upgrade {
			newPath, err := locator.LocateLatestAzureSDK(goSrcPath, oldSdkPath)
			if err != nil {
				return fmt.Errorf("Error determining if there's a later Azure SDK for %q: %s", oldSdkPath, err)
			}

			newSdkPath = *newPath
		}

		// now that we have this data, let's generate the follow up app
		// we have to do this since the Go compiler doesn't ship the types unless they're being used
		// hence this terrible hack
		if err := generateHack(goSrcPath, repositoryFullPath, oldSdkPath, newSdkPath, sdkTypes); err != nil {
			// we should error here but the AWS SDK references things in an `internal` directory
			// such that we can't guarantee this
			log.Printf("Error running hack: %s", err)
		}
	}

	return nil
}

func generateHack(goPath string, providerRepository string, oldPackage string, newPackage string, types []codegen.TerraformTypeInfo) error {
	folderName := "hack/"
	// delete the folder if it already exists
	os.RemoveAll(folderName)

	// then create it
	if err := os.Mkdir("hack", 0744); err != nil {
		return fmt.Errorf("Error creating directory: %s", err)
	}

	// then create some symlinks
	vendoredPath := fmt.Sprintf("%s/vendor/%s", providerRepository, oldPackage)
	oldFolder := "hack/old"
	cmd := exec.Command("ln", "-s", vendoredPath, oldFolder)
	cmd.Start()

	repoPath := fmt.Sprintf("%s/%s", goPath, newPackage)
	newFolder := "hack/new"
	cmd = exec.Command("ln", "-s", repoPath, newFolder)
	cmd.Start()

	// then generate the hack file
	if err := codegen.GenerateTerraformImports("hack/main.go", types); err != nil {
		return fmt.Errorf("[DEBUG] Error generating hack: %s", err)
	}

	// and finally run it
	if err := codegen.RunTerraformHack("hack/main.go"); err != nil {
		return fmt.Errorf("[DEBUG] Error running hack: %s", err)
	}

	return nil
}