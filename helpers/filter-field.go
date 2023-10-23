package helpers

import "groupie-tracker/types"

func NewField(name, label, typeStr string, minmax, checkbox bool, options bool) types.Field {

	return types.Field{
		Name:       name,
		Label:      label,
		Type:       typeStr,
		Minmax:     minmax,
		Checkbox:   checkbox,
		Options:    options,
		Attributes: make(map[string]string),
		Boxes:      make([]string, 0),
		Values:     make(map[string]bool),
		MMValues:   make(map[string]string),
	}
}
