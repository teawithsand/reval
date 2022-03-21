package reval

import "errors"

var ErrFieldNotFound = errors.New("reval: value with specified name was not found")
var ErrNilPointer = errors.New("reval: can query null pointer for field")
var ErrInvalidType = errors.New("reval: this accessor requires different type than provided")
