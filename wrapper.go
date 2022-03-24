package reval

import (
	"context"
	"reflect"

	"github.com/teawithsand/reval/stdesc"
)

// Type, which wraps arbitrary type in a Value interface.
type Wrapper interface {
	Wrap(v interface{}) (res Value, err error)
}

func MustWrap(res Value, err error) Value {
	if err != nil {
		panic(err)
	}

	return res
}

// TODO(teawithsand): implement support for embedded structures in JSON-like manner

type DefaultWrapper struct {
	// Computer, which will be used in structures in order to compute descriptor for field access.
	// Fallbacks to default computer if not set.
	DescriptorComputer *stdesc.Comptuer
}

// Util function, which converts go native type to Value.
func (dw *DefaultWrapper) Wrap(data interface{}) (v Value, err error) {
	if data == nil {
		v = nil
		return
	}

	var c *stdesc.Comptuer
	c = dw.DescriptorComputer
	if c == nil {
		c = &stdesc.Comptuer{}
	}

	switch tdata := data.(type) {
	case Value:
		v = tdata
		return
	case string:
		v = &PrimitiveValue{tdata}
		return
	case float64:
		v = &PrimitiveValue{tdata}
		return
	case float32:
		v = &PrimitiveValue{tdata}
		return
	case int:
		v = &PrimitiveValue{tdata}
		return
	case uint:
		v = &PrimitiveValue{tdata}
		return
	case int8:
		v = &PrimitiveValue{tdata}
		return
	case int16:
		v = &PrimitiveValue{tdata}
		return
	case int32:
		v = &PrimitiveValue{tdata}
		return
	case int64:
		v = &PrimitiveValue{tdata}
		return
	case uint8:
		v = &PrimitiveValue{tdata}
		return
	case uint16:
		v = &PrimitiveValue{tdata}
		return
	case uint32:
		v = &PrimitiveValue{tdata}
		return
	case uint64:
		v = &PrimitiveValue{tdata}
		return
	case reflect.Value: // wrapping reflect is nono
		err = ErrCantWrap
		return
	// TODO(teawithsand): add more primitive types here
	default:
		refVal := reflect.ValueOf(data)
		innerRefVal := refVal
		for innerRefVal.Kind() == reflect.Ptr {
			if innerRefVal.IsNil() {
				v = nil
				return
			}
			innerRefVal = innerRefVal.Elem()
		}

		// not nil, so we can operate on it
		// nil maps are not allowed, since we can't assign to them
		// we could construct one
		// but it's not really what we want from Wrap function
		if innerRefVal.Kind() == reflect.Map {
			v = &reflectMapValue{
				val:     innerRefVal,
				wrapper: dw,
			}
			return
		} else if innerRefVal.Kind() == reflect.Struct {
			var descriptor stdesc.Descriptor
			descriptor, err = c.ComputeDescriptor(context.Background(), refVal.Type())
			if err != nil {
				return
			}

			v = &reflectStructValue{
				val:        refVal,
				wrapper:    dw,
				descriptor: descriptor,
			}
			return
		} else if innerRefVal.Kind() == reflect.Slice || innerRefVal.Kind() == reflect.Array {
			v = &defaultListValue{
				val:     refVal,
				wrapper: dw,
			}
			return
		} else if refVal.Kind() == reflect.Ptr {
			return dw.Wrap(innerRefVal.Interface())
		} else {
			err = ErrCantWrap
			return
		}
	}
}
