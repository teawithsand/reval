package sttag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teawithsand/reval/sttag"
)

type TagOne struct {
	One       string   `sttag:",0"`
	Two       string   `sttag:",1"`
	Three     string   `sttag:"three"`
	Four      string   `sttag:"four"`
	Slice     []string `sttag:"slice"`
	Tag       bool     `sttag:"tag"`
	Anonymous string
}

func TestUnmarshal(t *testing.T) {
	du := sttag.NewDefaultUnmarshaler()

	assert := func(tag string, expected TagOne) {
		if t.Failed() {
			return
		}

		var tagOne TagOne
		err := du.Unmarshal(tag, &tagOne)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(tagOne, expected) {
			t.Error(fmt.Errorf("expected %+#v\ngot %+#v", expected, tagOne))
			return
		}
	}

	assert(`1,2,three:3,four:4,slice:1,slice:2,tag`, TagOne{
		One:   "1",
		Two:   "2",
		Three: "3",
		Four:  "4",
		Slice: []string{"1", "2"},
		Tag:   true,
	})
	assert(`,,tag`, TagOne{
		Tag: true,
	})
	assert(`,,`, TagOne{})
	assert(``, TagOne{})
	assert(`,,Anonymous:123`, TagOne{
		Anonymous: "123",
	})
}
