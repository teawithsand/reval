package jsonutil

import (
	"reflect"

	"github.com/teawithsand/reval/stdesc"
)

// FieldProcessor, which mimics encoding/json package.
// It makes fields appear as-if they were JSON fields.
func FieldProcesor(field reflect.StructField, path []int) (options stdesc.FieldOptions, err error) {
	fieldName, ok := GetJSONFieldName(field.Tag.Get("json"))
	if len(fieldName) == 0 && ok {
		options.Skip = true
		return
	}

	options.Name = fieldName
	if len(fieldName) == 0 {
		options.Name = field.Name
	}

	if field.IsExported() && field.Anonymous && fieldName == "" {
		options.Embed = true
	}

	return
}

var _ stdesc.FieldProcessor = FieldProcesor
