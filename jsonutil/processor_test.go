package jsonutil_test

import (
	"testing"

	"github.com/teawithsand/reval"
	"github.com/teawithsand/reval/jsonutil"
	"github.com/teawithsand/reval/stdesc"
)

type A struct {
	Field string `json:"field"`
}

func TestProcessor(t *testing.T) {
	dw := reval.DefaultWrapper{
		DescriptorComputer: &stdesc.Computer{
			FieldProcessorFactory: stdesc.FieldProcessorFunc(jsonutil.FieldProcesor),
		},
	}

	assert := func(x bool, msg string) {
		if t.Failed() {
			return
		}

		if !x {
			t.Error("assert filed", msg)
		}
	}

	v := reval.MustWrap(dw.Wrap(A{})).(reval.KeyedValue)

	assert(v.HasField("field"), "has json field")
	assert(!v.HasField("Field"), "does not have go field")
}
