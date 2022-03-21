package reval

import (
	"reflect"
)

// Type, which wraps arbitrary type in a Value interface.
type Wrapper interface {
	Wrap(v interface{}) (res Value, err error)
}

func WrapperMustWrap(v Value, err error) (res Value) {
	res = v
	if err != nil {
		panic(err)
	}
	return
}

// TODO(teawithsand): implement support for embedded structures in JSON-like manner

type DefaultWrapper struct {
	// If true, uses JSON names instead of actual field names when structure is wrapped.
	// Behavior is undefined, when json field tags are invalid
	UseJSONNames bool
}

// Util function, which converts go native type to Value.
func (dw *DefaultWrapper) Wrap(data interface{}) (v Value, err error) {
	if data == nil {
		v = nil
		return
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
	case int:
		v = &PrimitiveValue{tdata}
		return
	case reflect.Value: // wrapping reflect is nono
		err = &InvalidValueError{
			Data: data,
		}
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

		if innerRefVal.Kind() == reflect.Map || innerRefVal.Kind() == reflect.Struct {
			v = &mutableReflectKeyedValue{
				reflectKeyedValue: reflectKeyedValue{
					val:     refVal,
					wrapper: dw,
				},
			}
			return
		} else if innerRefVal.Kind() == reflect.Slice || innerRefVal.Kind() == reflect.Array {
			v = &defaultMutableListValue{
				defaultListValue: defaultListValue{
					val:     refVal,
					wrapper: dw,
				},
			}
			return
		} else if refVal.Kind() == reflect.Ptr {
			return dw.Wrap(innerRefVal.Interface())
		} else {
			err = &InvalidValueError{
				Data: data,
			}
			return
		}
	}
}

func (dw *DefaultWrapper) GetAlias(v KeyedValue, key interface{}) (alias interface{}, err error) {
	if !dw.UseJSONNames {
		alias = key
		return
	}

	refRaw := reflect.ValueOf(v.Raw())
	for refRaw.Kind() == reflect.Ptr {
		if refRaw.IsNil() {
			// error here?
			alias = key
			return
		}
		refRaw = refRaw.Elem()
	}

	// ignore non-struct
	if refRaw.Kind() != reflect.Struct {
		alias = key
		return
	}

	// TODO(teawithsand): cache this name map for performance
	nameMap := map[string]string{}

	// TODO(teawithsand): cache this, do something similar like JSON encoder
	length := refRaw.NumField()
	for i := 0; i < length; i++ {
		typeField := refRaw.Type().Field(i)

		jsonName, ok := getJSONFieldName(typeField.Tag.Get("json"))
		if ok {
			nameMap[jsonName] = typeField.Name
		}
	}

	stringKey, ok := key.(string)
	if !ok {
		alias = key
		return
	}

	alias, ok = nameMap[stringKey]
	if !ok {
		alias = key
	}
	return
}
