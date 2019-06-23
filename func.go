package errors

import (
	"reflect"
	"time"
)

func If(b bool, t, f interface{}) interface{} {

	if b {
		return t
	}

	return f
}

func FnCost(f func()) time.Duration {
	t1 := time.Now()
	ErrHandle(Try(f)())
	return time.Now().Sub(t1)
}

type FnT func(cfn ...interface{}) (err error)

func FnOf(fn interface{}, args ...interface{}) func() []reflect.Value {
	assertFn(fn)

	t := reflect.ValueOf(fn)
	return func() []reflect.Value {
		var vs []reflect.Value
		for i, p := range args {
			var _v reflect.Value
			if IsZero(p) {
				if t.Type().IsVariadic() {
					i = 0
				}
				_v = reflect.New(t.Type().In(i)).Elem()
			} else {
				_v = reflect.ValueOf(p)
			}

			vs = append(vs, _v)
		}
		return t.Call(vs)
	}
}
