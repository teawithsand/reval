package reval

import "reflect"

func getMapOrStructField(structOrMapVal reflect.Value, key interface{}) (v reflect.Value) {
	if structOrMapVal.Kind() != reflect.Struct && structOrMapVal.Kind() != reflect.Map {
		panic("reval: required struct/map reflect value")
	}

	if structOrMapVal.Kind() == reflect.Map {
		v = structOrMapVal.MapIndex(reflect.ValueOf(key))
		return
	} else {
		switch typedKey := key.(type) {
		case string:
			v = structOrMapVal.FieldByName(typedKey)
		case int:
			if typedKey >= structOrMapVal.NumField() {
				// report error?
			} else {
				v = structOrMapVal.Field(typedKey)
			}
		default:
			// report error?
		}
		return
	}
}

func getStructFieldType(structVal reflect.Type, key interface{}) (field reflect.StructField, ok bool) {
	if structVal.Kind() != reflect.Struct {
		panic("reval: required reflect value")
	}

	switch typedKey := key.(type) {
	case string:
		field, ok = structVal.FieldByName(typedKey)
	case int:
		if typedKey >= structVal.NumField() {
			// report error?
		} else {
			field = structVal.Field(typedKey)
			ok = true
		}
	default:
		// report error?
	}
	return
}
