package jsonutil

import (
	"github.com/teawithsand/reval/stdesc"
)

// FieldProcessor, which mimics encoding/json package.
// It makes fields appear as-if they were JSON fields.
func FieldProcesor(pf stdesc.PendingFiled) (options stdesc.FieldOptions, err error) {
	fieldName, ok := GetJSONFieldName(pf.Field.Tag.Get("json"))
	if len(fieldName) == 0 && ok {
		options.Skip = true
		return
	}

	options.Name = fieldName
	if len(fieldName) == 0 {
		options.Name = pf.Field.Name
	}

	if pf.Field.IsExported() && pf.Field.Anonymous && fieldName == "" {
		options.Embed = true
	}

	return
}

var _ stdesc.FieldProcessor = stdesc.FieldProcessorFunc(FieldProcesor)
