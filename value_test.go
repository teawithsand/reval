package reval_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teawithsand/reval"
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

func TestCanWrapValue(t *testing.T) {
	assert := func(v interface{}) {
		dw := reval.DefaultWrapper{}

		wrapped, err := dw.Wrap(v)
		if err != nil {
			t.Error(err, fmt.Sprintf("%T", v))
			return
		}
		rv := reflect.ValueOf(v)
		isNil := v == nil || rv.Kind() == reflect.Ptr && rv.IsNil()
		if isNil {
			if wrapped != nil {
				t.Error("nil mismatch")
			}
		} else {
			if reflect.TypeOf(v).Kind() != reflect.Ptr {
				if !reflect.DeepEqual(wrapped.Raw(), v) {
					t.Error("wrapped and raw not equal wrapped=", wrapped.Raw(), "v=", v)
					return
				}
			}
		}
	}

	assert(int(0))
	assert(string("asdf"))
	assert(float64(0))
	assert(float32(0))
	assert(rune('a'))
	assert(A{})
	assert(&A{})
	assert([2]int{1, 2})
	assert(make([]int, 3))
	assert(make(map[int]int))
	assert(nil)
	assert((*int)(nil))
	iv := 42
	assert(&iv)
	iva := &iv
	assert(&iva)
}

func Test_Struct_HasGetField(t *testing.T) {
	dw := reval.DefaultWrapper{}

	value := reval.MustWrap(dw.Wrap(A{
		B: B{
			C: &C{
				F: 42,
			},
		},
	})).(reval.KeyedValue)

	v, err := value.GetField("F")
	if err != nil {
		t.Error(err)
		return
	}
	if v.Raw() != 42 {
		t.Error("expected different value")
		return
	}
	if v.(*reval.PrimitiveValue).RawDereferenced() != 42 {
		t.Error("expected different value")
		return
	}

	if value.HasField("C") {
		t.Error("expected false")
		return
	}

	if !value.HasField("Q") {
		t.Error("expected true")
		return
	}

	if value.Len() != 4 {
		t.Error("invalid length")
		return
	}
}

func Test_Map_GetField(t *testing.T) {
	dw := reval.DefaultWrapper{}
	value := reval.MustWrap(dw.Wrap(map[int]int{
		2: 3,
		4: 5,
	})).(reval.KeyedValue)

	assert := func(x bool, msg string) {
		if t.Failed() {
			return
		}

		if !x {
			t.Error("assert filed", msg)
		}
	}

	assert(value.Len() == 2, "invalid len")
	assert(value.HasField(2), "invalid has field")
	assert(!value.HasField(42), "invalid has field 2")
	res, err := value.GetField(2)
	if err != nil {
		t.Error(err)
		return
	}
	assert(res.Raw() == 3, "invalid value")
}
