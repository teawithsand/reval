package sttag

import "strconv"

type unmarshalTagsFieldMeta struct {
	Skip            bool
	AnonymousOffset int
	KeyedName       string
}

func (meta *unmarshalTagsFieldMeta) Parse(tag string) (err error) {
	opts := SimpleParseOptions{
		AnonymousCount: 2,
	}

	res, err := opts.Parse(tag)
	if err != nil {
		return
	}

	keyedName := res.AnonymousValues.Get(0)
	meta.KeyedName = keyedName

	offset := res.AnonymousValues.Get(1)
	if len(offset) > 0 {
		var parsedOffset uint64
		parsedOffset, err = strconv.ParseUint(offset, 10, 32)
		if err != nil {
			return
		}

		meta.AnonymousOffset = int(parsedOffset)
	} else {
		meta.AnonymousOffset = -1
	}

	meta.Skip = len(offset) == 0 && meta.KeyedName == "-"
	return
}
