package reval

type Value interface {
	Raw() interface{}
}

// Value which consists of keys and value.
type KeyedValue interface {
	Value

	// Panics when no such field.
	// Must not return nil in that case.
	//
	// Returns nil value if field was not found.
	GetField(key interface{}) (res Value, err error)

	// Returns true if given field exists in value, false otherwise.
	HasField(name interface{}) bool

	// Iteration must stop when non-nil error is returned.
	// This error must be returned from top-level function.
	//
	// Note: field name yielded here is not value but primitive go type, like string or int.
	ListFields(recv func(name interface{}) (err error)) (err error)

	// Returns number of fields.
	Len() int
}

type MutableKeyedValue interface {
	KeyedValue

	// Returns true, if SetField with these parameters can be called without returning error.
	IsAssignable(name interface{}, value Value) bool

	// Panics when no such field.
	// Retruns error if not assignable.
	SetField(name interface{}, value Value) (err error)
}

type ListValue interface {
	Value

	// Panics if index is < 0 or out of bounds.
	GetIndex(i int) (res Value, err error)

	// Returns number of elements.
	Len() int
}

type MutableListValue interface {
	ListValue

	// Returns true, if SetField with these parameters can be called without returning error.
	IsAssignable(value Value) bool

	// Panics if index is < 0 or out of bounds.
	SetIndex(i int, value Value) (err error)
}
