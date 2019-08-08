package internal

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

var _funcCaller = func(callDepth int) []string {
	return []string{FuncCaller(callDepth), FuncCaller(callDepth + 1)}
}

func T(b bool, msg string, args ...interface{}) {
	if !b {
		return
	}

	panic(&Err{
		err:    fmt.Errorf(msg, args...),
		caller: _funcCaller(callDepth + 1),
	})
}

func TT(b bool, fn func(err *Err)) {
	if !b {
		return
	}

	_err := &Err{caller: _funcCaller(callDepth + 2)}
	fn(_err)

	if _err.msg == "" {
		log.Fatalf("msg is null")
	}
	_err.err = errors.New(_err.msg)

	panic(_err)
}

func Panic(err interface{}) {
	if err == nil || IsNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || IsNone(_m) {
		return
	}

	panic(&Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: _funcCaller(callDepth + 1),
	})
}

func Wrap(err interface{}, msg string, args ...interface{}) {
	if err == nil || IsNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || IsNone(_m) {
		return
	}

	panic(&Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: _funcCaller(callDepth + 1),
	})
}

func WrapM(err interface{}, fn func(err *Err)) {
	if err == nil || IsNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || IsNone(_m) {
		return
	}

	_err := &Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: _funcCaller(callDepth + 2),
	}
	fn(_err)
	panic(_err)
}

func AssertFn(fn reflect.Value) error {
	if IsZero(fn) || fn.Kind() != reflect.Func {
		return fmt.Errorf("the func is nil[%#v] or not func type[%s]", fn, fn.Kind().String())
	}
	return nil
}
