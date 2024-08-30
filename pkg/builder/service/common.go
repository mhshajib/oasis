package builder

type Field struct {
	Name    string
	Type    string
	JsonTag string
}
type CriteriaField struct {
	Name    string
	Type    string
	JsonTag string
}

func adjustFieldTypeForCriteria(fieldType string) string {
	// Check if it's a slice type
	if len(fieldType) > 2 && fieldType[:2] == "[]" {
		return "[]*" + fieldType[2:]
	}
	// Otherwise, make it a pointer type
	return "*" + fieldType
}
