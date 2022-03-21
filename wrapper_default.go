package reval

func MustWrap(data interface{}) (v Value) {
	v, err := Wrap(data)
	if err != nil {
		panic(err)
	}
	return
}

func Wrap(data interface{}) (v Value, err error) {
	dw := DefaultWrapper{}
	return dw.Wrap(data)
}
