package errors

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var errType = reflect.TypeOf(&Err{})

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

func _handle(err reflect.Value) reflect.Value {
	if IsZero(err) {
		return reflect.Value{}
	}

	for {
		if err.Kind() != reflect.Func {
			break
		}

		_ty := err.Type()

		T(_ty.NumIn() == 0 || _ty.IsVariadic(), "func input params num error")
		T(_ty.NumOut() != 1 || _ty.Out(0) != errType, "func output num and type error")

		err = err.Call([]reflect.Value{})[0]
		break
	}

	if IsZero(err) {
		return reflect.Value{}
	}

	m := &Err{}
	switch e := err.Interface().(type) {
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
		m.tag = _ErrTags.UnknownTypeCode
		_t := err.Type()
		m.m["type"] = _t.String()
		m.m["kind"] = _t.Kind()
		m.m["name"] = _t.Name()
	}
	return reflect.ValueOf(m)
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

func Handle(fn func()) {
	if fn == nil {
		log.Error().Msg("fn is nil")
		os.Exit(-1)
	}

	err := recover()
	if err == nil || IsZero(reflect.ValueOf(err)) {
		return
	}

	m := _handle(reflect.ValueOf(err))
	if IsZero(m) {
		return
	}

	_m:=m.Interface().(*Err)
	panic(&Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: getCallerFromFn(reflect.ValueOf(fn)),
	})
}
