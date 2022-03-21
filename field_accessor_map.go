package reval

import "reflect"

// MapFieldAccessor wraps map of any type to any into accessor
type MapFieldAccessor struct {
	reflect.Value
}

func (accessor *MapFieldAccessor) Get(v reflect.Value, key interface{}) (res interface{}, err error) {
	v = dereferenceValue(v)
	if isReflectZero(v) {
		err = ErrNilPointer
		return
	}
	if v.Kind() != reflect.Map {
		err = ErrInvalidType
		return
	}

	val := v.MapIndex(reflect.ValueOf(key))
	if isReflectZero(val) {
		err = ErrFieldNotFound
		return
	}

	res = val.Interface()

	return
}

func (accessor *MapFieldAccessor) ListFields(recv func(key interface{}) (err error)) (res interface{}, err error) {
	for _, k := range accessor.Value.MapKeys() {
		err = recv(k.Interface())
		if err != nil {
			return
		}
	}
	return
}
