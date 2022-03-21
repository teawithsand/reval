package reval

import (
	"errors"
)

var ErrNotStringable = errors.New("reval: specified value can't be converted to string")
var ErrNoField = errors.New("reval: such field does not exist on given value")
var ErrCantWrap = errors.New("reval: given value type can't be wrapped")
var ErrNilStruct = errors.New("reval: can't obtain struct field from nil pointer struct")
var ErrNilInnerStruct = errors.New("reval:  can't obtain struct field: embedded pointer to structure is nil")
