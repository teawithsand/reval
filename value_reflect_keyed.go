package reval

import (
	"reflect"
)

// reflectKeyedValue wraps any map or struct into value.
type reflectKeyedValue struct {
	val     reflect.Value
	wrapper Wrapper
}

var _ KeyedValue = &reflectKeyedValue{}

func (rkv *reflectKeyedValue) translareKey(key interface{}) (res interface{}, err error) {
	tw, ok := rkv.wrapper.(FieldAliasWrapper)
	if !ok {
		res = key
		return
	}
	res, err = tw.GetAlias(rkv, key)
	return
}

// returns element of pointer if any.
func (rkv *reflectKeyedValue) getInnerValue() (res reflect.Value) {
	res = rkv.val
	for res.Kind() == reflect.Ptr {
		res = res.Elem()
	}
	return
}

func (rkv *reflectKeyedValue) Raw() interface{} {
	return rkv.val.Interface()
}
func (rkv *reflectKeyedValue) innerGetReflectField(key interface{}) (v reflect.Value, err error) {
	iv := rkv.getInnerValue()
	key, err = rkv.translareKey(key)
	if err != nil {
		return
	}
	v = getMapOrStructField(iv, key)
	return
}

// Panics when no such field.
// Must not return nil in that case.
//
// Returns nil value if field was not found.
func (rkv *reflectKeyedValue) GetField(key interface{}) (res Value, err error) {
	v, err := rkv.innerGetReflectField(key)
	if err != nil {
		return
	}
	if isReflectZero(v) {
		return
	} else {
		res, err = rkv.wrapper.Wrap(v.Interface())
		return
	}

}

// Returns true if given field exists in value, false otherwise.
func (rkv *reflectKeyedValue) HasField(key interface{}) bool {
	f, err := rkv.innerGetReflectField(key)
	if err != nil {
		return false
	}
	return !isReflectZero(f)
}

// Iteration must stop when non-nil error is returned.
// This error must be returned from top-level function.
//
// Note: field name yielded here is not value but primitive go type, like string or int.
func (rkv *reflectKeyedValue) ListFields(recv func(name interface{}) (err error)) (err error) {
	iv := rkv.getInnerValue()

	if iv.Kind() == reflect.Map {
		for _, k := range iv.MapKeys() {
			err = recv(k.Interface())
			if err != nil {
				return
			}
		}
	} else {
		sz := iv.NumField()
		for i := 0; i < sz; i++ {
			err = recv(iv.Field(i).Interface())
			if err != nil {
				return
			}
		}
	}

	return
}

// Returns number of fields.
func (rkv *reflectKeyedValue) Len() int {
	iv := rkv.getInnerValue()

	if iv.Kind() == reflect.Map {
		return iv.Len()
	} else {
		return iv.NumField()
	}
}

type mutableReflectKeyedValue struct {
	reflectKeyedValue
}

func (mrkv *mutableReflectKeyedValue) IsAssignable(key interface{}, value Value) bool {
	if !mrkv.HasField(key) {
		return false
	}

	iv := mrkv.getInnerValue()
	if iv.Kind() == reflect.Map {
		return isAssignable(iv.Type().Elem(), value)
	} else {
		ift, err := mrkv.reflectKeyedValue.innerGetReflectField(key)
		if err != nil {
			return false
		}
		fieldType := ift.Type()
		return isAssignable(fieldType, value)
	}
}

func (mrkv *mutableReflectKeyedValue) SetField(key interface{}, value Value) (err error) {
	if !mrkv.HasField(key) {
		err = &NoFieldError{
			Value: mrkv,
			Field: key,
		}
		return
	}

	iv := mrkv.getInnerValue()
	if iv.Kind() == reflect.Map {
		err = assignMap(iv, key, value)
		if err != nil {
			return
		}
		return
	} else {
		var fieldRef reflect.Value
		fieldRef, err = mrkv.reflectKeyedValue.innerGetReflectField(key)
		if err != nil {
			return
		}

		err = assignValue(fieldRef, value)
		if err != nil {
			return
		}
	}

	return
}
