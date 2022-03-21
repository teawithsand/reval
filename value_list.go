package reval

import "reflect"

type defaultListValue struct {
	val     reflect.Value
	wrapper Wrapper
}

var _ ListValue = &defaultListValue{}

// returns element of pointer if any.
func (rkv *defaultListValue) getInnerValue() (res reflect.Value) {
	res = rkv.val
	for res.Kind() == reflect.Ptr {
		res = res.Elem()
	}
	return
}

func (lvw *defaultListValue) Raw() interface{} {
	return lvw.val.Interface()
}

// Panics if index is < 0 or out of bounds.

func (lvw *defaultListValue) GetIndex(i int) (res Value, err error) {
	iv := lvw.getInnerValue()
	res, err = lvw.wrapper.Wrap(iv.Index(i).Interface())
	return
}

type defaultMutableListValue struct {
	defaultListValue
}

var _ MutableListValue = &defaultMutableListValue{}

func (lvw *defaultListValue) IsAssignable(val Value) (ok bool) {
	return isAssignable(lvw.getInnerValue().Type().Elem(), val)
}

func (lvw *defaultListValue) SetIndex(i int, val Value) (err error) {
	iv := lvw.getInnerValue()
	err = assignList(iv, i, val)
	return
}

// Returns number of elements.
func (lvw *defaultListValue) Len() int {
	iv := lvw.getInnerValue()
	return iv.Len()
}
