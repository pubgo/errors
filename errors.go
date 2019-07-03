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

	m := _handle(reflect.ValueOf(err))
	if IsZero(m) {
		return nil
	}

	_m := m.Interface().(*Err)
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

	m := _handle(reflect.ValueOf(err))
	if IsZero(m) {
		return
	}

	_m := m.Interface().(*Err)
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

	m := _handle(reflect.ValueOf(err))
	if IsZero(m) {
		return
	}

	_m := m.Interface().(*Err)
	panic(&Err{
		sub:    _m,
		tag:    _m.tTag(),
		err:    _m.tErr(),
		caller: funcCaller(callDepth),
	})
}

func P(d ...interface{}) {
	defer Handle(func() {})

	for _, i := range d {
		if IsZero(reflect.ValueOf(i)) {
			return
		}

		if dt, err := json.MarshalIndent(i, "", "\t"); err != nil {
			Panic(err)
		} else {
			log.Info().Msg(string(dt))
		}
	}
}

func assertFn(fn reflect.Value) {
	T(IsZero(fn), "the func is nil")

	T(fn.Kind() != reflect.Func, "func type error(%s)", fn.Kind())
}
