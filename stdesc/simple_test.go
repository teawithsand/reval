package stdesc_test

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/teawithsand/reval/stdesc"
)

type C struct {
	F int
}

type B struct {
	*C
	N int
}

type A struct {
	B
	P int
	Q string
}

var fp = stdesc.FieldProcessorFunc(func(pf stdesc.PendingFiled) (options stdesc.FieldOptions, err error) {
	options.Name = pf.Field.Name
	options.Embed = pf.Field.Anonymous
	return
})

func TestSimpleStruct(t *testing.T) {
	c := stdesc.Comptuer{
		FieldProcessorFactory: fp,
	}

	desc, err := c.ComputeDescriptor(context.Background(), reflect.TypeOf(A{}))
	if err != nil {
		t.Error(err)
		return
	}

	var fields sort.StringSlice
	for k := range desc.NameToField {
		fields = append(fields, k)
	}

	sort.Sort(fields)

	if !reflect.DeepEqual([]string{"F", "N", "P", "Q"}, []string(fields)) {
		t.Error("fields not equal, got: ", fields)
		return
	}
}

func Test_FieldGet(t *testing.T) {
	d := A{
		P: 42,
		B: B{
			C: &C{
				F: 21,
			},
		},
	}

	c := stdesc.Comptuer{
		FieldProcessorFactory: fp,
	}

	desc, err := c.ComputeDescriptor(context.Background(), reflect.TypeOf(A{}))
	if err != nil {
		t.Error(err)
		return
	}

	fField := desc.NameToField["F"]

	res := fField.MustGet(reflect.ValueOf(d)).Interface()
	if res != d.F {
		t.Error("expected different value")
		return
	}

	res = fField.MustGet(reflect.ValueOf(&d)).Interface()
	if res != d.F {
		t.Error("expected different value")
		return
	}
}

func Test_FieldSet(t *testing.T) {
	d := &A{
		P: 42,
		B: B{
			C: &C{
				F: 21,
			},
		},
	}

	c := stdesc.Comptuer{
		FieldProcessorFactory: fp,
	}

	desc, err := c.ComputeDescriptor(context.Background(), reflect.TypeOf(A{}))
	if err != nil {
		t.Error(err)
		return
	}

	fField := desc.NameToField["F"]
	fField.MustSet(reflect.ValueOf(d), reflect.ValueOf(11))
	if 11 != d.F {
		t.Error("expected different value")
		return
	}
}
