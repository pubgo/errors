package internal

import (
	"fmt"
	"reflect"
)

var _funcCaller = func(callDepth int) []string {
	return []string{FuncCaller(callDepth), FuncCaller(callDepth + 1)}
}

func T(b bool, msg string, args ...interface{}) {
	if !b {
		return
	}

	_err := fmt.Errorf(msg, args...)
	panic(&Err{
		err:    _err,
		msg:    _err.Error(),
		caller: _funcCaller(callDepth + 1),
	})
}

func TT(b bool, msg string, args ...interface{}) *Err {
	if !b {
		return nil
	}

	_err := fmt.Errorf(msg, args...)
	return &Err{
		err:    _err,
		msg:    _err.Error(),
		caller: _funcCaller(callDepth + 1),
	}
}

func Panic(err interface{}) {
	if err == nil {
		return
	}

	_m := _handle(err)
	if  _m == nil || IsNone(_m) {
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
	if err == nil {
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

func WrapM(err interface{}, msg string, args ...interface{}) *Err {
	if err == nil {
		return nil
	}

	_m := _handle(err)
	if _m == nil || IsNone(_m) {
		return nil
	}

	return &Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: _funcCaller(callDepth + 1),
	}
}

func AssertFn(fn reflect.Value) error {
	if IsZero(fn) || fn.Kind() != reflect.Func {
		return fmt.Errorf("the func is nil(%#v) or type error(%s)", fn, fn.Kind().String())
	}
	return nil
}
