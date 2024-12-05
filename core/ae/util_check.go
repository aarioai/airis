package ae

func Check(es ...*Error) *Error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

func CheckErrors(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func PanicIf(cond bool, tip ...any) {
	if cond {
		if len(tip) > 0 {
			panic(tip)
		}
		panic("PanicIf")
	}
}

// PanicOn 如果存在服务器错误则触发 panic
func PanicOn(es ...*Error) {
	if len(es) == 0 {
		return
	}

	for _, e := range es {
		if e != nil {
			panic("app.PanicOn: " + e.Text())
		}
	}
}

// PanicOnErrors 断言检查标准错误，如果存在错误则触发 panic
func PanicOnErrors(errs ...error) {
	if len(errs) == 0 {
		return
	}

	for _, err := range errs {
		if err != nil {
			panic("app.PanicOnErrors: " + err.Error())
		}
	}
}
