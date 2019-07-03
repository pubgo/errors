package errors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
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
		caller: funcCaller(callDepth),
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
		caller: funcCaller(callDepth),
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
		caller: funcCaller(callDepth),
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
		caller: funcCaller(callDepth),
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
		caller: funcCaller(callDepth),
	})
}

func P(d ...interface{}) {
	defer Handle()()

	for _, i := range d {
		if IsZero(reflect.ValueOf(i)) {
			continue
		}

		dt, err := json.MarshalIndent(i, "", "\t")
		Wrap(err, "P json MarshalIndent error")
		log.Info().Msg(string(dt))
	}
}

func assertFn(fn reflect.Value) {
	T(IsZero(fn), "the func is nil")

	T(fn.Kind() != reflect.Func, "func type error: "+fn.Kind().String())
}
