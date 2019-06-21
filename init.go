package errors

import (
	"fmt"
	"reflect"
	"runtime"
)

const callDepth = 2

func assertFn(fn interface{}) {
	T(IsNil(fn), "the func is nil")

	_v := reflect.TypeOf(fn)
	T(_v.Kind() != reflect.Func, "func type error(%s)", _v.String())
}

func funcCaller(callDepth int) string {
	fn, _, line, ok := runtime.Caller(callDepth)
	if !ok {
		return "no func caller"
	}
	return fmt.Sprintf("%s:%d", runtime.FuncForPC(fn).Name(), line)
}

func IsNil(p interface{}) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = false
		}
	}()

	if p == nil {
		return true
	}

	if !reflect.ValueOf(p).IsValid() {
		return true
	}

	return reflect.ValueOf(p).IsNil()
}
