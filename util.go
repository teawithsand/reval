package reval

import "reflect"

func isReflectZero(v reflect.Value) bool {
	return v == reflect.Value{}
}
