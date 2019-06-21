package errors

import (
	"errors"
	"fmt"
	"reflect"
)

func Try(fn interface{}, cfn ...interface{}) (err error) {
	defer func() {
		m := kerrGet()
		defer kerrPut(m)

		if r := recover(); r != nil {
			switch d := r.(type) {
			case *Err:
				m = d
			case error:
				m.err = d
				m.msg = d.Error()
				m.caller = funcCaller(callDepth)
			case string:
				m.err = errors.New(d)
				m.msg = d
				m.caller = funcCaller(callDepth)
			default:
				m.msg = fmt.Sprintf("type error %#v", d)
				m.err = errors.New(m.msg)
				m.caller = funcCaller(callDepth)
				m.tag = ErrTag.UnknownErr
			}
		}

		if m.err == nil {
			err = nil
			return
		}
		err = m.copy()
	}()

	assertFn(fn)
	v := FnOf(fn)()
	if len(cfn) > 0 {
		assertFn(cfn[0])
		reflect.ValueOf(cfn[0]).Call(v)
	}
	return nil
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if IsNil(err) {
		return
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

	fmt.Printf("%#v", err)
}
