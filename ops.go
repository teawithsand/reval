package reval

import "reflect"

type OPs interface {
	// Returns true if values are equal.
	// Note: this equal is like JS equal, so in some cases is not symmetric, transistent and so on.
	Eq(v1, v2 Value) (ok bool)

	// Returns "length" of value.
	Len(v Value) (len int, err error)

	// Casts value to string or error if possible.
	Stringify(v Value) (res string, err error)

	// Empty thing is now handled by value itself.
	// It's here for future use. It should always return value's IsEmpty().
	IsEmpty(v Value) bool
}

type DefaultOPs struct{}

func (dvo *DefaultOPs) Eq(v1, v2 Value) (ok bool) {
	// if same then equal
	if v1 == v2 {
		return true
	}

	/*
		// If raw equal then equal
		if v1.Raw() == v2.Raw() {
			return true
		}
	*/

	// If deep equal then equal
	if reflect.DeepEqual(v1.Raw(), v2.Raw()) {
		return true
	}

	// nasty mode - if strings equal then equal
	// Note: this is not compliant since it makes NaN == NaN
	// which is not valid...
	sv1, err := stringifyValue(v1)
	if err != nil {
		return false
	}
	sv2, err := stringifyValue(v2)
	if err != nil {
		return false
	}
	if sv1 == sv2 {
		return true
	}

	return false
}

func (dvo *DefaultOPs) Stringify(v Value) (res string, err error) {
	res, err = stringifyValue(v)
	return
}

func (dvo *DefaultOPs) Len(val Value) (sz int, err error) {
	res, err := stringifyValue(val)
	if err != nil {
		return
	}
	sz = len(res)
	return
}

func (dvo *DefaultOPs) IsEmpty(v Value) bool {
	if v == nil {
		return true
	}
	r := v.Raw()
	refRaw := reflect.ValueOf(r)
	for refRaw.Kind() == reflect.Ptr {
		if refRaw.IsNil() {
			return true
		}
		refRaw = refRaw.Elem()
	}

	unwrappedRaw := refRaw.Interface()
	return unwrappedRaw == nil || unwrappedRaw == ""
}
