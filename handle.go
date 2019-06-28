package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		err.P()
	})
}

func Log() {
	ErrHandle(recover(), func(err *Err) {
		err.Log()
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), fn)
}

func _handle(err interface{}) *Err {
	if e, ok := err.(func() (err error)); ok {
		err = e()
	}

	if e, ok := err.(func(...interface{}) (err error)); ok {
		err = e()
	}

	if IsZero(err) {
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
		m.msg = fmt.Sprintf("handle type error %#v", e)
		m.err = errors.New(m.msg)
		m.tag = ErrTag.UnknownErr
		_t := reflect.TypeOf(err)
		m.m["type"] = _t.String()
		m.m["kind"] = _t.Kind()
		m.m["name"] = _t.Name()

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

	m := _handle(err)
	if IsZero(m) {
		return
	}

	panic(&Err{
		sub:    m,
		caller: getCallerFromFn(fn),
		err:    m.tErr(),
	})
}
