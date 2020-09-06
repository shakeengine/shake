package misc

// ErrorCheck : If err is not nil, occur panic
func ErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}
