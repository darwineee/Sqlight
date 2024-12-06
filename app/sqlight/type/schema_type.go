package _type

type SchemaType string

func IsValidSchemaType(t string) bool {
	switch t {
	case "table", "index", "view", "trigger":
		return true
	default:
		return false
	}
}
