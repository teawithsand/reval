package reval

import (
	"reflect"
)

type ProcessFieldOptions struct {
	Name  string
	Embed bool
}

type structFieldMeta struct {
	index    int
	accessor FieldAccessor
}

type defaultFieldAccessor struct {
	ty       reflect.Type
	registry map[string]structFieldMeta
}

func (fa *defaultFieldAccessor) Get(v reflect.Value, key interface{}) (res interface{}, err error) {
	v = dereferenceValue(v)
	if isReflectZero(v) {
		err = ErrNilPointer
		return
	}

	stringKey, ok := key.(string)
	if !ok {
		err = ErrFieldNotFound
		return
	}

	if v.Kind() != reflect.Struct {
		err = ErrInvalidType
		return
	}

	fieldData, ok := fa.registry[stringKey]
	if !ok {
		err = ErrFieldNotFound
		return
	}

	f := v.Field(fieldData.index)

	if fieldData.accessor == nil {
		res = f.Interface()
		return
	} else {
		res, err = fieldData.accessor.Get(f, stringKey)
		return
	}
}

func (fa *defaultFieldAccessor) ListFields(recv func(key interface{}) (err error)) (err error) {
	for key, reg := range fa.registry {
		if reg.accessor != nil {
			err = reg.accessor.ListFields(recv)
			if err != nil {
				return
			}
		} else {
			err = recv(key)
			if err != nil {
				return
			}
		}
	}

	return
}

type StructFieldAccessorFactory struct{}

// Creates struct accessor for given struct using processor given.
func (fac *StructFieldAccessorFactory) MakeAccessor(processor StructFieldProcessor, ty reflect.Type) (res FieldAccessor, err error) {
	ty = dereferenceType(ty)

	v, ok := cache.Load(ty)
	if ok {
		res = v.(FieldAccessor)
		return
	}

	if ty.Kind() != reflect.Struct {
		err = ErrInvalidType
		return
	}

	length := ty.NumField()

	fieldRegistry := map[string]structFieldMeta{}

	for i := 0; i < length; i++ {
		f := ty.Field(i)

		var options ProcessFieldOptions
		options, err = processor.ProcessField(f)
		if err != nil {
			return
		}

		meta := structFieldMeta{
			index: i,
		}

		if options.Embed {
			meta.accessor, err = fac.MakeAccessor(processor, ty)
			if err != nil {
				return
			}
		}

		fieldRegistry[options.Name] = meta
	}

	accessor := &defaultFieldAccessor{}
	cache.Store(ty, accessor)
	// res = accessor
	panic("NIY")
	return
}
