package stdesc

import (
	"context"
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
	Type reflect.Type // Field type, the most inner one, which can be acessed via get/set

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

	field := target.FieldByIndex(f.Path)

	if v.Type().AssignableTo(field.Type()) {
		field.Set(v)
	} else if reflect.PointerTo(v.Type()).AssignableTo(field.Type()) {
		if v.CanAddr() {
			field.Set(v.Addr())
		} else {
			panic(fmt.Errorf("value of type %s is assignable to type %s but can't take address of value", v.Type(), field.Type()))
		}
	} else {
		panic(fmt.Errorf("value of type %s (nor pointer to it) is not assignable to field type %s", v.Type(), field.Type()))
	}
}

type Descriptor struct {
	Ty reflect.Type

	// also alows listing fields
	NameToField map[string]Field

	// Comptued summary of descriptor, performed by user code.
	ComputedSummary interface{}
}

type FieldOptions struct {
	Skip bool // If true, field is skipped. Ignores all other settings.

	Embed bool // If true, and type is struct, then all fields from child structs are embedded as if they were in parent.
	Name  string
	Meta  interface{}

	Override bool // if true, and field with same name is set then overrides previous definiton in case it was set.
}

type PendingFiled struct {
	Field reflect.StructField
	Path  []int
}

type Summarizer interface {
	SummaryzeDescriptor(ctx context.Context, desc Descriptor) (meta interface{}, err error)
}
type SummarizerFunc func(ctx context.Context, desc Descriptor) (meta interface{}, err error)

func (f SummarizerFunc) SummaryzeDescriptor(ctx context.Context, desc Descriptor) (meta interface{}, err error) {
	return f(ctx, desc)
}

// FieldProcessor decides how field should be processed.
type FieldProcessor interface {
	ProcessField(pf PendingFiled) (options FieldOptions, err error)
}

type FieldProcessorFunc func(pf PendingFiled) (options FieldOptions, err error)

func (f FieldProcessorFunc) ProcessField(pf PendingFiled) (options FieldOptions, err error) {
	return f(pf)
}

func (f FieldProcessorFunc) MakeFieldProcessor(ctx context.Context, ty reflect.Type) (fp FieldProcessor, err error) {
	fp = f
	return
}

type FieldProcessorFactory interface {
	MakeFieldProcessor(ctx context.Context, ty reflect.Type) (fp FieldProcessor, err error)
}

type FieldProcessorFactoryFunc func(ctx context.Context, ty reflect.Type) (fp FieldProcessor, err error)

func (f FieldProcessorFactoryFunc) MakeFieldProcessor(ctx context.Context, ty reflect.Type) (fp FieldProcessor, err error) {
	return f(ctx, ty)
}

type Computer struct {
	// Fallbacks to processor, which embeds all anonmous fields and sets name to field name.
	// Note: path parameter may not be modified by this function.
	FieldProcessorFactory FieldProcessorFactory

	// Performs descriptor summarization.
	Summarizer Summarizer

	// Descriptor cache for each type.
	// Used if not nil.
	Cache *sync.Map
}

func (c *Computer) innerComputeDescriptor(
	ctx context.Context,
	fp FieldProcessor,
	path []int,
	ty reflect.Type,
	desc *Descriptor,
) (err error) {
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
		field.Type = structField.Type

		var options FieldOptions
		options, err = fp.ProcessField(PendingFiled{
			Field: structField,
			Path:  path,
		})
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

		err = c.innerComputeDescriptor(ctx, fp, path, ef.sf.Type, desc)
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
func (c *Computer) ComputeDescriptor(ctx context.Context, ty reflect.Type) (desc Descriptor, err error) {
	desc.NameToField = map[string]Field{}

	var fp FieldProcessor
	if c.FieldProcessorFactory != nil {
		fp, err = c.FieldProcessorFactory.MakeFieldProcessor(ctx, ty)
		if err != nil {
			return
		}
	}

	if fp == nil {
		fp = FieldProcessorFunc(func(pf PendingFiled) (options FieldOptions, err error) {
			options.Name = pf.Field.Name
			options.Embed = IsEmbedField(pf)
			return
		})
	}

	err = c.innerComputeDescriptor(ctx, fp, []int{}, ty, &desc)
	if err != nil {
		return
	}

	if c.Summarizer != nil {
		var summary interface{}
		summary, err = c.Summarizer.SummaryzeDescriptor(ctx, desc)
		if err != nil {
			return
		}

		desc.ComputedSummary = summary
	}
	return
}
