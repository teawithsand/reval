package stdesc

import "reflect"

// Returns true, if field should be embedded, in order to get encoding/json like embedding.
// It's used by default FieldProcessor.
func IsEmbedField(pf PendingFiled) bool {
	return pf.Field.Anonymous && (pf.Field.Type.Kind() == reflect.Struct || (pf.Field.Type.Kind() == reflect.Ptr && pf.Field.Type.Elem().Kind() == reflect.Struct))
}
