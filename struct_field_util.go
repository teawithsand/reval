package reval

import "reflect"

func copyIndices(data []int) []int {
	res := make([]int, len(data))
	copy(res, data)
	return res
}

func dereferenceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

func dereferenceType(ty reflect.Type) reflect.Type {
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty
}

func isReflectZero(v reflect.Value) bool {
	return v == reflect.Value{}
}
