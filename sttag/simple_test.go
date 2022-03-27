package sttag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teawithsand/reval/sttag"
)

func TestSimpleParseTag(t *testing.T) {
	opts := sttag.SimpleParseOptions{
		AnonymousCount: 2,
	}

	assertParseOk := func(tag string, exp sttag.SimpleParsedTag) {
		if t.Failed() {
			return
		}

		res, err := opts.Parse(tag)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(res, exp) {
			t.Error(fmt.Errorf("expected different value while parsing `%s`, got\n%+#v expected\n%+#v", tag, res, exp))
			return
		}
	}

	assertParseOk(``, sttag.SimpleParsedTag{})
	assertParseOk(`asdf`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf"},
	})
	assertParseOk(`"asdf"`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf"},
	})
	assertParseOk(`"asdf",asdf`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
	})
	assertParseOk(`"asdf",asdf,keyedBlank`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
		NamedValues: sttag.NamedValues{
			"keyedBlank": nil,
		},
	})
	assertParseOk(`"asdf",asdf,keyedBlank:`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
		NamedValues: sttag.NamedValues{
			"keyedBlank": []string{""},
		},
	})
	assertParseOk(`"asdf",asdf,keyedBlank:""`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
		NamedValues: sttag.NamedValues{
			"keyedBlank": []string{""},
		},
	})
	assertParseOk(`"asdf",asdf,keyedBlank:,keyed:fdsa`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
		NamedValues: sttag.NamedValues{
			"keyedBlank": []string{""},
			"keyed":      []string{"fdsa"},
		},
	})
	assertParseOk(`"asdf",asdf,keyedBlank:,keyed:"fd\"sa"`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
		NamedValues: sttag.NamedValues{
			"keyedBlank": []string{""},
			"keyed":      []string{"fd\"sa"},
		},
	})
	assertParseOk(`"asdf",asdf,keyedBlank:,keyed:"fd\"sa"`, sttag.SimpleParsedTag{
		AnonymousValues: sttag.AnonymousValues{"asdf", "asdf"},
		NamedValues: sttag.NamedValues{
			"keyedBlank": []string{""},
			"keyed":      []string{"fd\"sa"},
		},
	})
}
