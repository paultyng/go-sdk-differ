package differ

import (
	"reflect"
	"strings"
)

func Diff(oldObj interface{}, newObj interface{}) (*DiffResult, error) {
	oldElem := reflect.ValueOf(oldObj)
	newElem := reflect.ValueOf(newObj)

	fields := make(map[string]string, 0)

	result := DiffResult{
		Added: []Added{},
		Changed: []Changed{},
		Removed: []Removed{},
	}

	// iterate over the new object to get all it's fields and types
	for i := 0; i < newElem.NumField(); i++ {
		field := newElem.Type().Field(i)
		fields[field.Name] = field.Type.String()
	}

	// then iterate over the old object
	for i := 0; i < oldElem.NumField(); i++ {
		field := oldElem.Type().Field(i)

		newFieldType, ok := fields[field.Name]
		if !ok {
			// this fields been removed
			result.Removed = append(result.Removed, Removed{
				FieldName: field.Name,
				FieldType: field.Type.String(),
			})
			delete(fields, field.Name)
			continue
		}

		// has the type changed?
		if !strings.EqualFold(newFieldType, field.Type.String()) {
			result.Changed = append(result.Changed, Changed{
				FieldName: field.Name,
				OldType: field.Type.String(),
				NewType: newFieldType,
			})
			delete(fields, field.Name)
			continue
		}

		delete(fields, field.Name)
	}

	// anything remaining must be new
	if len(fields) > 0 {
		for k, v := range fields {
			result.Added = append(result.Added, Added{
				FieldName: k,
				FieldType: v,
			})
		}
	}

	return &result, nil
}

func (r DiffResult) HasChanges() bool {
	return len(r.Added) != 0 || len(r.Changed) != 0 || len(r.Removed) != 0
}