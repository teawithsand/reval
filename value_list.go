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
		if res.IsNil() {
			return reflect.Value{}
		}
		res = res.Elem()
	}
	return
}

func (lvw *defaultListValue) Raw() interface{} {
	return lvw.val.Interface()
}

func (lvw *defaultListValue) GetIndex(i int) (res Value, err error) {
	iv := lvw.getInnerValue()
	if isReflectZero(iv) || i < 0 || i >= lvw.Len() {
		err = ErrNoField
		return
	}
	res, err = lvw.wrapper.Wrap(iv.Index(i).Interface())
	return
}

// Returns number of elements.
func (lvw *defaultListValue) Len() int {
	iv := lvw.getInnerValue()
	if isReflectZero(iv) {
		return 0
	}
	return iv.Len()
}
