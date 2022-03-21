package reval

import (
	"errors"
	"reflect"
)

type accessorMapValue struct {
	val      reflect.Value
	isNil    bool
	wrapper  Wrapper
	accessor FieldAccessor
}

var _ KeyedValue = &accessorMapValue{}

func (rmv *accessorMapValue) Raw() interface{} {
	return rmv.val.Interface()
}

func (rmv *accessorMapValue) GetField(key interface{}) (res Value, err error) {
	if rmv.isNil {
		return
	}

	rawRes, err := rmv.accessor.Get(rmv.val, key)
	if errors.Is(err, ErrFieldNotFound) || errors.Is(err, ErrNilPointer) {
		err = nil
		return
	} else if err != nil {
		return
	}

	res, err = rmv.wrapper.Wrap(rawRes)
	return
}

// Returns true if given field exists in value, false otherwise.
func (rmv *accessorMapValue) HasField(name interface{}) bool {
	if rmv.isNil {
		return false
	}

	return rmv.accessor.HasField(name)
}

func (rmv *accessorMapValue) ListFields(recv func(name interface{}) (err error)) (err error) {
	if rmv.isNil {
		return
	}
	return rmv.accessor.ListFields(recv)
}

// Returns number of fields.
func (rmv *accessorMapValue) Len() int {
	return rmv.accessor.Len()
}
