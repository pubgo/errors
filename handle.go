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

func Handle1(fn ...func(m *M)) {
	err := recover()
	if IsNil(err) {
		return
	}

	var _m = newM()
	if len(fn) > 0 && !IsNil(fn[0]) {
		fn[0](_m)
	}

	if len(_m.m) == 0 {
		_m.m = nil
	}

	var caller string
	if _m.caller != "" {
		caller = _m.caller
	} else {
		caller = funcCaller(5)
	}

	if Cfg.Debug {
		log.Println(_m.msg, caller)
	}

	m := _handle(err)
	panic(&Err{
		sub:    m,
		caller: caller,
		err:    m.tErr(),
		msg:    _m.msg,
		tag:    m.tTag(_m.tag),
		m:      _m.m,
	})
}

func Handle(fn func()) {
	assertFn(fn)

	err := recover()
	if IsNil(err) {
		return
	}

	_fn := reflect.ValueOf(fn).Pointer()
	_e := runtime.FuncForPC(_fn)
	file, line := _e.FileLine(_fn)
	ma := strings.Split(_e.Name(), ".")
	caller := strings.TrimPrefix(strings.TrimPrefix(fmt.Sprintf("%s:%d:%s", file, line, ma[len(ma)-1]), srcDir), modDir)

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
