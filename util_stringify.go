package reval

import (
	"fmt"
	"reflect"
)

// stringifyValue stringifies value for some things like getting it's length
// It should not be used in strict mode.
func stringifyValue(val Value) (res string, err error) {
	raw := val.Raw()

	switch trv := raw.(type) {
	case string:
		res = trv
	case fmt.Stringer:
		res = trv.String()
	default:
		// TODO(teawithsand): fix this, this is just unsound
		//  do separate casts for different types
		//  and fail in default branch

		ty := reflect.TypeOf(val)
		for ty.Kind() == reflect.Ptr {
			ty = ty.Elem()
		}

		k := ty.Kind()
		if k == reflect.Struct || k == reflect.Interface || k == reflect.Func {
			err = ErrNotStringable
			return
		}

		res = fmt.Sprintf("%d", trv)
	}

	return
}
