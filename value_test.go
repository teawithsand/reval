package reval_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teawithsand/reval"
)

type someStruct struct {
	A int
	B *int
	C string
}

func TestCanWrapValue(t *testing.T) {
	assert := func(v interface{}) {
		wrapped, err := reval.Wrap(v)
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
	assert(someStruct{})
	assert(&someStruct{})
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

func TestCanOperateOnStructValue_WithStruct(t *testing.T) {
	b := 3
	v, err := reval.Wrap(&someStruct{
		A: 2,
		B: &b,
		C: "asdf",
	})
	if err != nil {
		t.Error(err)
	}
	mkv := v.(reval.MutableKeyedValue)

	checkBValOk := func(b int) {
		bField, err := mkv.GetField("B")
		if err != nil {
			t.Error(err)
			return
		}
		if bField.Raw() != b {
			t.Error("invalid value stored in b=", bField.Raw())
			return
		}
	}
	checkBValOk(b)

	nv := 4
	err = mkv.SetField("B", reval.MustWrap(&nv))
	if err != nil {
		t.Error(err)
		return
	}
	checkBValOk(nv)

	nv = 5
	err = mkv.SetField("B", reval.MustWrap(nv))
	if err != nil {
		t.Error(err)
		return
	}

	checkBValOk(nv)

	err = mkv.SetField("C", reval.MustWrap("asdf"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCanOperateOnStructValue_WithMap(t *testing.T) {
	b := 3
	v, err := reval.Wrap(&map[string]interface{}{
		"A": 2,
		"B": &b,
		"C": "asdf",
	})
	if err != nil {
		t.Error(err)
	}
	mkv := v.(reval.MutableKeyedValue)

	checkBValOk := func(b int) {
		bField, err := mkv.GetField("B")
		if err != nil {
			t.Error(err)
			return
		}
		if bField.Raw() != b {
			t.Error("invalid reval stored in b=", bField.Raw())
			return
		}
	}
	checkBValOk(b)

	nv := 4
	err = mkv.SetField("B", reval.MustWrap(&nv))
	if err != nil {
		t.Error(err)
		return
	}
	checkBValOk(nv)

	nv = 5
	err = mkv.SetField("B", reval.MustWrap(nv))
	if err != nil {
		t.Error(err)
		return
	}

	checkBValOk(nv)

	err = mkv.SetField("C", reval.MustWrap("asdf"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCanOperateOnListValue_Slice(t *testing.T) {
	v, err := reval.Wrap([]int{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	mlv := v.(reval.MutableListValue)

	err = mlv.SetIndex(0, reval.MustWrap(4))
	if err != nil {
		t.Error(err)
		return
	}

	{
		v, err := mlv.GetIndex(0)
		if err != nil {
			t.Error(err)
			return
		}

		if v.Raw() != int(4) {
			t.Error("expected reval to be equal to newly set")
			return
		}
	}
}
