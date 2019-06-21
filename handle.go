package errors

import (
	"errors"
	"fmt"
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

func Handle(fn ...func(m *M)) {
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

	m := _handle(err)
	panic(&Err{
		sub:    m,
		caller: funcCaller(4),
		err:    m.tErr(),
		msg:    _m.msg,
		tag:    m.tTag(_m.tag),
		m:      _m.m,
	})
}
