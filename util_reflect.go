package reval

import "reflect"

func isAssignable(ty reflect.Type, val Value) (ok bool) {
	return reflect.TypeOf(val.Raw()).AssignableTo(ty) || reflect.PtrTo(reflect.TypeOf(val.Raw())).AssignableTo(ty)
}

// assignValue
func assignValue(fieldRef reflect.Value, value Value) (err error) {
	if !fieldRef.CanSet() {
		return &NotSettableValueError{
			Data: value,
		}
	}

	refValue := reflect.ValueOf(value.Raw())

	if refValue.Type().AssignableTo(fieldRef.Type()) {
		fieldRef.Set(refValue)
		return
	} else if reflect.PtrTo(refValue.Type()).AssignableTo(fieldRef.Type()) {
		var setValue reflect.Value
		if !refValue.CanAddr() {
			setValue = reflect.New(refValue.Type())
			// copying is allowed when we receive non-ptr type
			// since primitive types are always non-ptr
			setValue.Elem().Set(refValue)
		} else {
			setValue = refValue.Addr()
		}

		fieldRef.Set(setValue)
		return
	}

	err = &NotAssignableValueError{
		To:    fieldRef.Type(),
		Value: value,
	}

	return
}

func assignMap(mapRefVal reflect.Value, key interface{}, value Value) (err error) {
	if mapRefVal.Kind() != reflect.Map {
		panic("reval: required map reflect value")
	}

	refValue := reflect.ValueOf(value.Raw())

	if refValue.Type().AssignableTo(mapRefVal.Type().Elem()) {
		mapRefVal.SetMapIndex(reflect.ValueOf(key), refValue)
		return
	} else if reflect.PtrTo(refValue.Type()).AssignableTo(mapRefVal.Type().Elem()) {
		var setValue reflect.Value
		if !refValue.CanAddr() {
			setValue = reflect.New(refValue.Type())
			// copying is allowed when we receive non-ptr type
			// since primitive types are always non-ptr
			setValue.Elem().Set(refValue)
		} else {
			setValue = refValue.Addr()
		}

		mapRefVal.SetMapIndex(reflect.ValueOf(key), setValue)
		return
	} else {
		err = &NotAssignableValueError{
			To:    mapRefVal.Type().Elem(),
			Value: value,
		}
		return
	}
}

func assignList(listVal reflect.Value, i int, value Value) (err error) {
	if listVal.Kind() != reflect.Array && listVal.Kind() != reflect.Slice {
		panic("reval: required array/slice reflect value")
	}

	if i > listVal.Len() || i < 0 {
		// TODO(teawithsand): OOB index handling here
	}

	err = assignValue(listVal.Index(i), value)
	return
}

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
