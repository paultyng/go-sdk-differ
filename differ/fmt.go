package differ

import "fmt"

func (r DiffResult) Print() string {
	if !r.HasChanges() {
		return fmt.Sprintf("No changes!")
	}

	output := ""
	if len(r.Added) > 0 {
		output += "Added:\n"
		for _, v := range r.Added {
			output += fmt.Sprintf("  %s (%s)\n", v.FieldName, v.FieldType)
		}
	}

	if len(r.Changed) > 0 {
		output += "Changed:\n"
		for _, v := range r.Changed {
			output += fmt.Sprintf("  %s (from %s to %s)\n", v.FieldName, v.OldType, v.NewType)
		}
	}

	if len(r.Removed) > 0 {
		output += "Removed:\n"
		for _, v := range r.Removed {
			output += fmt.Sprintf("  %s (%s)\n", v.FieldName, v.FieldType)
		}
	}

	return output
}
