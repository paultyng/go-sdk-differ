package codegen

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type TerraformTypeInfo struct {
	Package string
	ImportType string
}

type TerraformTypes struct {
	Info TerraformTypeInfo
	OldType reflect.Type
	NewType reflect.Type
}

func GenerateTerraformImports(fileName string, types []TerraformTypeInfo) error {
	codegen := ""
	for _, v := range types {
		codegen += fmt.Sprintf(`{
			Info: codegen.TerraformTypeInfo{
				Package: %q,
				ImportType: %q,
			},
			OldType: reflect.TypeOf(old.%s{}),
			NewType: reflect.TypeOf(new.%s{}),
		},
`, v.Package, v.ImportType, v.ImportType, v.ImportType)
	}

	template := `package main

import (
    "log"
	"reflect"

    "github.com/tombuildsstuff/go-sdk-differ/codegen"
	"github.com/tombuildsstuff/go-sdk-differ/differ"
    new "github.com/tombuildsstuff/go-sdk-differ/hack/new"
    old "github.com/tombuildsstuff/go-sdk-differ/hack/old"
)

func main() {
	for _, v := range GetTerraformTypes() {
		old := reflect.New(v.OldType)
		new := reflect.New(v.NewType)
		result, err := differ.Diff(old, new)
		if err != nil {
			panic(err)
		}

        if !result.HasChanges() {
          log.Printf("        Type %q - no changes", v.Info.ImportType)
        } else {
		  log.Printf("        Type %q", v.Info.ImportType)
		  log.Printf("          %s", result.Print())
		}
	}
}

func GetTerraformTypes() []codegen.TerraformTypes {
	return []codegen.TerraformTypes{
        [[types]]
	}
}
`

	templated := strings.Replace(template, "[[types]]", codegen, 1)

	if err := ioutil.WriteFile(fileName, []byte(templated), 0744); err != nil {
		return fmt.Errorf("Error writing template: %s", err)
	}

	return nil
}
