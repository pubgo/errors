package errors

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		err.P()
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), fn)
}

func _handle(err interface{}) *Err {
	if e, ok := err.(func() (err error)); ok {
		err = e()
	}

	if e, ok := err.(FnT); ok {
		err = e()
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
		m.msg = fmt.Sprintf("type error %#v", e)
		m.err = errors.New(m.msg)
		m.tag = ErrTag.UnknownErr
	}
	return m
}

func getCallerFromFn(fn interface{}) string {
	_fn := reflect.ValueOf(fn).Pointer()
	_e := runtime.FuncForPC(_fn)
	file, line := _e.FileLine(_fn)
	ma := strings.Split(_e.Name(), ".")
	return strings.TrimPrefix(strings.TrimPrefix(fmt.Sprintf("%s:%d:%s", file, line, ma[len(ma)-1]), srcDir), modDir)
}

func Handle(fn func()) {
	assertFn(fn)

	err := recover()
	if IsZero(err) {
		return
	}

	caller := getCallerFromFn(fn)

	if Cfg.Debug {
		log.Println("handle", caller)
	}

	m := _handle(err)
	panic(&Err{
		sub:    m,
		caller: caller,
		err:    m.tErr(),
	})
}
