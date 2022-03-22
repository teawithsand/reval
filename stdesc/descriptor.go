package stdesc

import (
	"fmt"
	"reflect"
	"sync"
)

func copyIndices(data []int) []int {
	res := make([]int, len(data))
	copy(res, data)
	return res
}

type Field struct {
	Name string
	Path []int // used to access fields via indexes, always at least one element in here
	Meta interface{}

	// Path to field, which contains struct, which contains this field
	ParentPath []int
}

// Returns value of given field or panics
func (f *Field) MustGet(v reflect.Value) (res reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			panic("reval/sdesc: pointer to structure is nil")
		}
		v = v.Elem()
	}

	return v.FieldByIndex(f.Path)
}

// Sets value of field.
// Note: panics if value is not addressable.
func (f *Field) MustSet(target, v reflect.Value) {
	if target.Kind() == reflect.Ptr {
		if target.IsNil() {
			panic("reval/sdesc: pointer to structure is nil")
		}
		target = target.Elem()
	}

	target.FieldByIndex(f.Path).Set(v)
}

type Descriptor struct {
	Ty reflect.Type

	// also alows listing fields
	NameToField map[string]Field
}

type FieldOptions struct {
	Skip bool // If true, field is skipped. Ignores all other settings.

	Embed bool // If true, and type is struct, then all fields from child structs are embedded as if they were in parent.
	Name  string
	Meta  interface{}

	Override bool // if true, and field with same name is set then overrides previous definiton in case it was set.
}

// FieldProcessor decides how field should be processed.
type FieldProcessor func(field reflect.StructField, path []int) (options FieldOptions, err error)

type Comptuer struct {
	// Fallbacks to processor, which embeds all anonmous fields and sets name to field name.
	// Note: path parameter may not be modified by this function.
	FieldProcessor FieldProcessor

	// Descriptor cache for each type.
	// Used if not nil.
	Cache *sync.Map
}

func (c *Comptuer) innerComputeDescriptor(path []int, ty reflect.Type, desc *Descriptor) (err error) {
	if ty.Kind() == reflect.Pointer {
		ty = ty.Elem()
	}

	if c.Cache != nil {
		cachedDescriptor, ok := c.Cache.Load(ty)
		if ok {
			*desc = cachedDescriptor.(Descriptor)
			return
		}
	}

	fp := c.FieldProcessor

	if fp == nil {
		fp = func(field reflect.StructField, path []int) (options FieldOptions, err error) {
			options.Name = field.Name
			options.Embed = field.Anonymous && (field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct))
			return
		}
	}

	// use this trick to embed fields of structure before embeded structure fields
	var embedFields []struct {
		index int
		sf    reflect.StructField
	}

	length := ty.NumField()
	for i := 0; i < length; i++ {
		path = append(path, i)

		var field Field

		structField := ty.Field(i)

		var options FieldOptions
		options, err = fp(structField, path)
		if err != nil {
			return
		}

		if !options.Skip || !structField.IsExported() {
			if options.Embed {
				fieldType := structField.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				if fieldType.Kind() == reflect.Struct {
					embedFields = append(embedFields, struct {
						index int
						sf    reflect.StructField
					}{
						index: i,
						sf:    structField,
					})
				} else {
					err = fmt.Errorf("reval/sdesc: can't embed non-struct or not pointer-to-struct field; path: %+#v", path)
					return
				}
			} else {
				field.Name = options.Name
				field.Meta = options.Meta
				field.Path = copyIndices(path)

				_, exists := desc.NameToField[options.Name]
				if !exists || options.Override {
					desc.NameToField[options.Name] = field
				}
			}
		}

		path = path[:len(path)-1]
	}

	for _, ef := range embedFields {
		path = append(path, ef.index)

		err = c.innerComputeDescriptor(path, ef.sf.Type, desc)
		if err != nil {
			return
		}

		path = path[:len(path)-1]
	}

	if c.Cache != nil {
		c.Cache.Store(ty, *desc)
	}

	return
}

// Note: returned descriptor should be deep copied before modifying
// since it may be stored in cache.
func (c *Comptuer) ComputeDescriptor(ty reflect.Type) (desc Descriptor, err error) {
	desc.NameToField = map[string]Field{}

	err = c.innerComputeDescriptor([]int{}, ty, &desc)
	if err != nil {
		return
	}
	return
}
