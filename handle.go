package errors

import (
	"errors"
	"fmt"
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

func _handle(err reflect.Value) reflect.Value {
	if IsZero(err) {
		return reflect.Value{}
	}

	if err.Kind() == reflect.Func {
		_ty := err.Type()

		T(_ty.NumIn() == 0 && _ty.IsVariadic(), "func input params num error")
		T(_ty.NumOut() != 1, "func output num error, num: "+strconv.Itoa(_ty.NumOut()))

		_v := valueGet()
		defer valuePut(_v)
		
		err = err.Call(_v)[0]
		return reflect.Value{}
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
		m.tag = ErrTags.UnknownTypeCode
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

func Handle() func() {
	_caller := funcCaller(2)

	return func() {
		err := recover()
		if err == nil || IsZero(reflect.ValueOf(err)) {
			return
		}

		m := _handle(reflect.ValueOf(err))
		if IsZero(m) {
			return
		}

		_m := m.Interface().(*Err)
		panic(&Err{
			sub:    _m,
			tag:    _m.tTag(),
			err:    _m.tErr(),
			caller: _caller,
		})
	}
}
