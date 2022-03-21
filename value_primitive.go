package reval

import "reflect"

type PrimitiveValue struct {
	Val interface{}
}

func (pv *PrimitiveValue) Raw() interface{} {
	if pv == nil {
		return nil
	}
	return pv.Val
}

// Returns value after stripping pointer layer.
// Returns nil if any pointer is nil.
func (pv *PrimitiveValue) RawDereferenced() interface{} {
	v := reflect.ValueOf(pv.Val)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Interface()
}
