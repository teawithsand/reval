package sttag

import (
	"errors"
	"io"
	"strings"
)

type SimpleParseOptions struct {
	AnonymousCount int
}

type AnonymousValues []string

func (av AnonymousValues) Has(value string) bool {
	for _, v := range av {
		if v == value {
			return true
		}
	}

	return false
}

func (av AnonymousValues) HasAfter(value string, i int) bool {
	if i >= len(av) || i < 0 {
		return false
	}

	return av[i:].Has(value)
}

func (av AnonymousValues) Get(i int) string {
	if i >= len(av) || i < 0 {
		return ""
	}
	return av[i]
}

type NamedValues map[string][]string

func (av NamedValues) IsSet(key string) bool {
	if av == nil {
		return false
	}
	return len(av[key]) > 0
}

func (av NamedValues) GetFirst(key string) string {
	if av == nil {
		return ""
	}

	imm := av[key]

	if len(imm) > 0 {
		return imm[0]
	}
	return ""
}

func (av NamedValues) Get(key string) []string {
	if av == nil {
		return nil
	}

	imm := av[key]
	return imm
}

type SimpleParsedTag struct {
	AnonymousValues AnonymousValues
	NamedValues     NamedValues
}

var ErrTagInvalid = errors.New("reval/sttag: invalid tag structure")

const sepNone = '*'
const quoteChar = '"'

func (spo *SimpleParseOptions) readString(rd *strings.Reader, sep map[rune]struct{}) (res string, sepUsed rune, err error) {
	state := 0
	for {
		var r rune
		r, _, err = rd.ReadRune()
		if errors.Is(err, io.EOF) {
			if state == 0 || state == -1 {
				if len(res) > 0 {
					err = nil
					sepUsed = sepNone
				}

				return
			} else {
				err = io.ErrUnexpectedEOF
				return
			}
		} else if err != nil {
			return
		}

		if state == -1 {
			if _, ok := sep[r]; !ok {
				err = ErrTagInvalid
				return
			} else {
				sepUsed = r
				return
			}
		} else if state == 0 {
			if r == quoteChar {
				state = 1
			} else if _, ok := sep[r]; ok {
				sepUsed = r
				return
			} else {
				res += string(r)
			}
		} else if state == 1 {
			if r == '\\' {
				state = 2
			} else if r == quoteChar {
				state = -1
			} else {
				res += string(r)
			}
		} else if state == 2 {
			res += string(r)
			state = 1
		}
	}
}

func (spo *SimpleParseOptions) Parse(tag string) (res SimpleParsedTag, err error) {
	rd := strings.NewReader(tag)

	for i := 0; i < spo.AnonymousCount; i++ {
		var text string
		text, _, err = spo.readString(rd, map[rune]struct{}{
			',': {},
		})
		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
		res.AnonymousValues = append(res.AnonymousValues, text)
	}

	exit := false
	for {
		if exit {
			break
		}

		var key, value string
		var topLevelSep rune
		key, topLevelSep, err = spo.readString(rd, map[rune]struct{}{
			',': {},
			':': {},
		})
		if errors.Is(err, io.EOF) {
			err = nil
			return
		} else if err != nil {
			return
		} else if topLevelSep == ',' {
			continue
		}

		if res.NamedValues == nil {
			res.NamedValues = NamedValues{}
		}

		if topLevelSep == ':' {
			value, _, err = spo.readString(rd, map[rune]struct{}{
				',': {},
			})
			if errors.Is(err, io.EOF) {
				err = nil
				res.NamedValues[key] = append(res.NamedValues[key], "")
				return
			}

			res.NamedValues[key] = append(res.NamedValues[key], value)
		} else if topLevelSep == sepNone {
			_, ok := res.NamedValues[key]
			if !ok {
				res.NamedValues[key] = nil
			}
		}

	}

	return
}
