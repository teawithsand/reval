package reval

import (
	"reflect"

	"github.com/teawithsand/reval/stdesc"
)

// reflectStructValue wraps any map or struct into value.
type reflectStructValue struct {
	descriptor stdesc.Descriptor
	val        reflect.Value
	wrapper    Wrapper
}

var _ KeyedValue = &reflectStructValue{}

// returns element of pointer if any.
func (rkv *reflectStructValue) getInnerValue() (res reflect.Value) {
	res = rkv.val
	for res.Kind() == reflect.Ptr {
		if res.IsNil() {
			return reflect.Value{}
		}
		res = res.Elem()
	}
	return
}

func (rkv *reflectStructValue) Raw() interface{} {
	return rkv.val.Interface()
}

// Panics when no such field.
// Must not return nil in that case.
//
// Returns nil value if field was not found.
func (rkv *reflectStructValue) GetField(key interface{}) (res Value, err error) {
	v := rkv.getInnerValue()

	if isReflectZero(v) {
		err = ErrNilStruct
		return
	} else {
		stringKey, ok := key.(string)
		if !ok {
			err = ErrNoField
			return
		}

		field, ok := rkv.descriptor.NameToField[stringKey]
		if !ok {
			err = ErrNoField
			return
		}

		rawResult := field.MustGet(v)

		if isReflectZero(rawResult) {
			err = ErrNilInnerStruct
			return
		}

		res, err = rkv.wrapper.Wrap(rawResult.Interface())
		return
	}
}

// Returns true if given field exists in value, false otherwise.
func (rkv *reflectStructValue) HasField(key interface{}) bool {
	v := rkv.getInnerValue()

	if isReflectZero(v) {
		return false
	} else {
		stringKey, ok := key.(string)
		if !ok {
			return false
		}

		field, ok := rkv.descriptor.NameToField[stringKey]
		if !ok {
			return false
		}

		rawResult := field.MustGet(v)

		if isReflectZero(rawResult) {
			return false
		}

		return true
	}
}

// Iteration must stop when non-nil error is returned.
// This error must be returned from top-level function.
//
// Note: field name yielded here is not value but primitive go type, like string or int.
func (rkv *reflectStructValue) ListFields(recv func(name interface{}) (err error)) (err error) {
	for _, f := range rkv.descriptor.NameToField {
		if rkv.HasField(f.Name) {
			err = recv(f.Name)
			if err != nil {
				return
			}
		}
	}

	return
}

// Returns number of fields.
func (rkv *reflectStructValue) Len() int {
	return len(rkv.descriptor.NameToField)
}
