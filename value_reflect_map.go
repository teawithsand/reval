package reval

import (
	"reflect"
)

type reflectMapValue struct {
	val     reflect.Value
	wrapper Wrapper
}

var _ KeyedValue = &reflectMapValue{}

func (rmv *reflectMapValue) Raw() interface{} {
	return rmv.val.Interface()
}

func (rmv *reflectMapValue) getValue() reflect.Value {
	v := rmv.val
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}
		}

		v = v.Elem()
	}

	return v
}

func (rmv *reflectMapValue) GetField(key interface{}) (res Value, err error) {
	v := rmv.getValue()
	if isReflectZero(v) {
		err = ErrNoField
		return
	}

	raw := v.MapIndex(reflect.ValueOf(key))
	if isReflectZero(raw) {
		err = ErrNoField
		return
	}

	res, err = rmv.wrapper.Wrap(raw.Interface())
	if err != nil {
		return
	}
	return
}

// Returns true if given field exists in value, false otherwise.
func (rmv *reflectMapValue) HasField(name interface{}) bool {
	v := rmv.getValue()
	if isReflectZero(v) {
		return false
	}

	return !isReflectZero(v.MapIndex(reflect.ValueOf(name)))
}

func (rmv *reflectMapValue) ListFields(recv func(name interface{}) (err error)) (err error) {
	v := rmv.getValue()
	if isReflectZero(v) {
		return
	}

	for _, info := range v.MapKeys() {
		err = recv(info.Interface())
		if err != nil {
			return
		}
	}

	return
}

// Returns number of fields.
func (rmv *reflectMapValue) Len() int {
	v := rmv.getValue()
	if isReflectZero(v) {
		return 0
	}

	return v.Len()
}
