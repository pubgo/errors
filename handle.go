package errors

import (
	"errors"
	"fmt"
	"os"
)

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		fmt.Println(err.P())
	})
}

func Assert() {
	ErrHandle(recover(), func(err *Err) {
		fmt.Println(err.P())
		os.Exit(1)
	})
}

func Throw() {
	ErrHandle(recover(), func(err *Err) {
		fmt.Println(err.P())
		panic(err)
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), fn)
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
		m.err = errors.New(m.msg)
		m.tag = ErrTags.UnknownTypeCode
	}
	return m
}
