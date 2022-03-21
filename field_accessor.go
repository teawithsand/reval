package reval

import (
	"reflect"
	"sync"
)

type StructFieldProcessor interface {
	// Processes field, and decides what to do with it:
	// 1. what name it should be accessible with
	// 2. Should this struct be omitted / embedded
	ProcessField(f reflect.StructField) (res ProcessFieldOptions, err error)
}

type FieldAccessor interface {
	Get(v reflect.Value, key interface{}) (res interface{}, err error)
	ListFields(recv func(key interface{}) (err error)) (err error)
	Len() int
	HasField(key interface{}) bool
}

var cache sync.Map // map[reflect.Type]FieldAccessor // TODO(teawithsand): move to generics
