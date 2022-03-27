package sttag

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/teawithsand/reval/stdesc"
)

const unmarshalerTagName = "sttag"

// Unmarshaler unmarshals tag value into type provided using some tag format.
type Unmarshaler interface {
	Unmarshal(tag string, res interface{}) (err error)
}

// defaultUnmarshaler, which loads specified structure from struct tag it's given.
type defaultUnmarshaler struct {
	computer *stdesc.Computer
}

type defaultUnmarshalerSummary struct {
	anonCount int
}

// Unmarshaler, which uses default format.
func NewDefaultUnmarshaler() *defaultUnmarshaler {
	return &defaultUnmarshaler{
		computer: &stdesc.Computer{
			FieldProcessorFactory: stdesc.FieldProcessorFunc(func(pf stdesc.PendingFiled) (options stdesc.FieldOptions, err error) {
				options.Name = pf.Field.Name
				options.Embed = stdesc.IsEmbedField(pf)

				meta := unmarshalTagsFieldMeta{}
				err = meta.Parse(pf.Field.Tag.Get(unmarshalerTagName))
				if err != nil {
					return
				}

				if len(meta.KeyedName) == 0 && meta.AnonymousOffset < 0 {
					meta.KeyedName = pf.Field.Name
				}

				options.Meta = meta
				options.Skip = meta.Skip || !(pf.Field.Type.Kind() == reflect.String || pf.Field.Type.Kind() == reflect.Bool || (pf.Field.Type.Kind() == reflect.Slice && pf.Field.Type.Elem().Kind() == reflect.String))

				return
			}),
			Summarizer: stdesc.SummarizerFunc(func(ctx context.Context, desc stdesc.Descriptor) (meta interface{}, err error) {
				anonCount := 0
				for _, f := range desc.NameToField {
					meta := f.Meta.(unmarshalTagsFieldMeta)
					if meta.Skip {
						continue
					}

					if meta.AnonymousOffset >= 0 {
						anonCount += 1
					}
				}

				meta = defaultUnmarshalerSummary{
					anonCount: anonCount,
				}
				return
			}),
			Cache: &sync.Map{},
		},
	}
}

func (um *defaultUnmarshaler) Unmarshal(tag string, res interface{}) (err error) {
	refRes := reflect.ValueOf(res)

	desc, err := um.computer.ComputeDescriptor(context.Background(), reflect.TypeOf(res))
	if err != nil {
		return
	}

	summ := desc.ComputedSummary.(defaultUnmarshalerSummary)

	opts := SimpleParseOptions{
		AnonymousCount: summ.anonCount,
	}

	simple, err := opts.Parse(tag)
	if err != nil {
		return
	}

	for _, f := range desc.NameToField {
		meta := f.Meta.(unmarshalTagsFieldMeta)

		fmt.Printf("%+#v\n", meta)

		if meta.AnonymousOffset >= 0 {
			v := simple.AnonymousValues.Get(meta.AnonymousOffset)
			if len(v) > 0 {
				if f.Type.Kind() == reflect.String {
					f.MustSet(refRes, reflect.ValueOf(v))
				} else {
					err = ErrInvalidFieldType
					return
				}
			}
		} else if len(meta.KeyedName) > 0 {
			if f.Type.Kind() == reflect.String {
				v := simple.NamedValues.GetFirst(meta.KeyedName)
				f.MustSet(refRes, reflect.ValueOf(v))
			} else if f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.String {
				v := simple.NamedValues.Get(meta.KeyedName)
				f.MustSet(refRes, reflect.ValueOf(v))
			} else if f.Type.Kind() == reflect.Bool {
				_, ok := simple.NamedValues[meta.KeyedName]
				f.MustSet(refRes, reflect.ValueOf(
					ok,
				))
			} else {
				err = ErrInvalidFieldType
				return
			}
		}
	}

	return
}
