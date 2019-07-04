package errors

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"runtime"
	"strconv"
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
	if err == nil || IsZero(reflect.ValueOf(err)) {
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

	if err == nil || IsZero(reflect.ValueOf(err)) {
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
		m.tag = ErrTags.UnknownTypeCode
	}
	return m
}

func getCallerFromFn(fn reflect.Value) string {
	_fn := fn.Pointer()
	_e := runtime.FuncForPC(_fn)
	file, line := _e.FileLine(_fn)
	ma := strings.Split(_e.Name(), ".")

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(".")
	buf.WriteString(ma[len(ma)-1])
	return strings.TrimPrefix(strings.TrimPrefix(buf.String(), srcDir), modDir)
}

func Handle() func() {
	_caller := funcCaller(2)

	if _l := log.Debug(); _l.Enabled() {
		_l.Msg(_caller)
	}

	return func() {
		err := recover()
		if err == nil || IsZero(reflect.ValueOf(err)) {
			return
		}

		_m := _handle(err)
		if _m == nil || IsZero(reflect.ValueOf(_m)) {
			return
		}

		panic(&Err{
			sub:    _m,
			tag:    _m.tTag(),
			err:    _m.tErr(),
			caller: _caller,
		})
	}
}
