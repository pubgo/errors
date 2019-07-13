package errors

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		err.caller = append(err.caller[:len(err.caller)-1], funcCaller(callDepth))
		fmt.Println(err.P())
	})
}

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		fmt.Println(err.P())
	})
}

func Assert() {
	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		fmt.Println(err.P())
		os.Exit(1)
	})
}

func Throw(fn func()) {
	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		err.caller = append(err.caller, getCallerFromFn(reflect.ValueOf(fn)))
		panic(err)
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		err.caller = append(err.caller, getCallerFromFn(reflect.ValueOf(fn)))
		fn(err)
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
