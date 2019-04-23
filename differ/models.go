package differ

type DiffResult struct {
	Added []Added
	Changed []Changed
	Removed []Removed
}

type Added struct {
	FieldName string
	FieldType string
}

type Changed struct {
	FieldName string
	OldType string
	NewType string
}

type Removed struct {
	FieldName string
	FieldType string
}
