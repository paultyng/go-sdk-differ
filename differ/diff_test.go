package differ

import "testing"

type FirstStruct struct {
	Id *string `json:"id"`
	Name *string `json:"name"`
}

type SecondStruct struct {
	Id *int `json:"id"`
	Tags *[]string `json:"tags"`
	Children *[]FirstChild `json:"children"`
}

type FirstChild struct {
	Id *string `json:"id"`
}

type ThirdStruct struct {
	Id *int `json:"id"`
	Tags *[]string `json:"tags"`
	Children *[]SecondChild `json:"children"`
}

type SecondChild struct {
	Id *string `json:"id"`
	Name *string `json:"name"`
}

func TestDiff(t *testing.T) {
	data := []struct{
		Input             interface{}
		Output            interface{}
		ShouldHaveChanges bool
		Expected          DiffResult
	}{
		{
			Input:             struct{}{},
			Output:            struct{}{},
			ShouldHaveChanges: false,
		},
		{
			Input: FirstStruct{},
			Output: FirstStruct{},
			ShouldHaveChanges: false,
		},
		{
			Input: FirstStruct{},
			Output: SecondStruct{
				Children: &[]FirstChild{},
			},
			ShouldHaveChanges: true,
		},
		{
			Input: SecondStruct{
				Children: &[]FirstChild{},
			},
			Output: SecondStruct{
				Children: &[]FirstChild{},
			},
			ShouldHaveChanges: false,
		},
		{
			Input: SecondStruct{
				Children: &[]FirstChild{},
			},
			Output: ThirdStruct{
				Children: &[]SecondChild{},
			},
			ShouldHaveChanges: true,
		},
	}

	for _, v := range data {
		diff, err := Diff(v.Input, v.Output)
		if err != nil {
			t.Fatal(err)
		}

		if v.ShouldHaveChanges && !diff.HasChanges() {
			t.Fatal("Expected changes but got nothing")
		}

		if !v.ShouldHaveChanges && diff.HasChanges() {
			t.Fatalf("Expected no changes but got %s", diff.Print())
		}
	}
}
