package internal

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
)

func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		err.Caller(FuncCaller(callDepth))
		fmt.Println(err.P())
	})
}

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		fmt.Println(err.P())
		debug.PrintStack()
	})
}

func Throw(fn interface{}) {
	_fn := reflect.ValueOf(fn)
	T(fn == nil || IsZero(_fn) || _fn.Kind() != reflect.Func, "the input must be func type and not null, input --> %#v", fn)

	ErrHandle(recover(), func(err *Err) {
		err.Caller(GetCallerFromFn(_fn))
		panic(err)
	})
}

func _handle(err interface{}) *Err {
	if err == nil || IsNone(err) {
		return nil
	}

	if _e, ok := err.(func() (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) func(...interface{}) error); ok {
		err = _e()()
	}

	if _e, ok := err.(func(...interface{}) func(...interface{}) func(...interface{}) error); ok {
		err = _e()()()
	}

	if err == nil || IsNone(err) {
		return nil
	}

	m := &Err{}
	switch e := err.(type) {
	case *Err:
		m = e
	case error:
		m.msg = e.Error()
		m.err = e
	case string:
		m.err = errors.New(e)
		m.msg = e
	default:
		m.msg = fmt.Sprintf("unknown type error, input: %#v", e)
		m.err = errors.New("unknown type error")
		m.tag = errTags.UnknownTypeCode
	}
	return m
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if err == nil || IsNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || IsNone(_m) || _m.err == nil {
		return
	}

	if len(fn) == 0 {
		return
	}

	Wrap(AssertFn(reflect.ValueOf(fn[0])), "func error")
	fn[0](_m)
}
