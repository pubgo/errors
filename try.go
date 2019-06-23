package errors

import (
	"errors"
	"fmt"
	"reflect"
)

func Try(fn interface{}, args ...interface{}) func(...interface{}) (err error) {
	return func(cfn ...interface{}) (err error) {
		defer func() {
			m := kerrGet()
			defer kerrPut(m)

			if r := recover(); !IsZero(r) {
				caller := getCallerFromFn(fn)

				switch d := r.(type) {
				case *Err:
					m = d
					m.caller = caller
				case error:
					m.err = d
					m.msg = d.Error()
					m.caller = caller
				case string:
					m.err = errors.New(d)
					m.msg = d
					m.caller = caller
				default:
					m.msg = fmt.Sprintf("type error %#v", d)
					m.err = errors.New(m.msg)
					m.caller = caller
					m.tag = ErrTag.UnknownErr
				}
			}

			if m.err == nil {
				err = nil
				return
			}
			err = m.copy()
		}()

		_call := FnOf(fn, args...)
		if len(cfn) == 0 {
			_call()
			return
		}

		assertFn(cfn[0])
		reflect.ValueOf(cfn[0]).Call(_call())
		return
	}
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if IsZero(err) {
		return
	}

	if _e, ok := err.(func() (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(*Err); ok {
		if len(fn) > 0 {
			assertFn(fn[0])
			fn[0](_e)
		}
		return
	}

	if _e, ok := err.(error); ok {
		fmt.Println(_e.Error())
		return
	}

	fmt.Printf("%#v\n", err)
}
