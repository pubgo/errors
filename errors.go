package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func T(b bool, msg string, args ...interface{}) {
	if !b {
		return
	}

	_err := fmt.Errorf(msg, args...)
	panic(&Err{
		err:    _err,
		msg:    _err.Error(),
		caller: []string{funcCaller(callDepth), funcCaller(callDepth + 1)},
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
		caller: []string{funcCaller(callDepth), funcCaller(callDepth + 1)},
	}
}

func WrapM(err interface{}, msg string, args ...interface{}) *Err {
	if err == nil {
		return nil
	}

	_m := _handle(err)
	if _m == nil || IsZero(reflect.ValueOf(_m)) {
		return nil
	}

	return &Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: []string{funcCaller(callDepth), funcCaller(callDepth + 1)},
	}
}

func Wrap(err interface{}, msg string, args ...interface{}) {
	if err == nil {
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
		msg:    fmt.Sprintf(msg, args...),
		caller: []string{funcCaller(callDepth), funcCaller(callDepth + 1)},
	})
}

func Panic(err interface{}) {
	if err == nil {
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
		caller: []string{funcCaller(callDepth), funcCaller(callDepth + 1)},
	})
}

func P(d ...interface{}) {
	for _, i := range d {
		if IsZero(reflect.ValueOf(i)) {
			continue
		}

		dt, err := json.MarshalIndent(i, "", "\t")
		Wrap(err, "P json MarshalIndent error")
		fmt.Println(string(dt))
	}
}

func assertFn(fn reflect.Value) {
	T(IsZero(fn), "the func is nil")

	T(fn.Kind() != reflect.Func, "func type error: "+fn.Kind().String())
}
