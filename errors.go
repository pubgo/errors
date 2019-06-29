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
		return &Err{}
	}

	_err := fmt.Errorf(msg, args...)
	return &Err{
		err:    _err,
		msg:    _err.Error(),
		caller: funcCaller(callDepth),
	}
}

func WrapM(err interface{}, msg string, args ...interface{}) *Err {
	m := _handle(err)
	if IsZero(m) {
		return &Err{}
	}

	return &Err{
		sub:    m,
		tag:    m.tTag(),
		err:    m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: funcCaller(callDepth),
	}
}

func Wrap(err interface{}, msg string, args ...interface{}) {
	m := _handle(err)
	if IsZero(m) {
		return
	}

	panic(&Err{
		sub:    m,
		tag:    m.tTag(),
		err:    m.tErr(),
		msg:    fmt.Sprintf(msg, args...),
		caller: funcCaller(callDepth),
	})
}

func Panic(err interface{}) {
	m := _handle(err)
	if IsZero(m) {
		return
	}

	panic(&Err{
		sub:    m,
		tag:    m.tTag(),
		err:    m.tErr(),
		caller: funcCaller(callDepth),
	})
}

func P(d ...interface{}) {
	defer Handle(func() {})

	for _, i := range d {
		if IsZero(i) {
			return
		}

		if dt, err := json.MarshalIndent(i, "", "\t"); err != nil {
			Panic(err)
		} else {
			log.Info().Msg(string(dt))
		}
	}
}

func assertFn(fn interface{}) {
	T(IsZero(fn), "the func is nil")

	_v := reflect.TypeOf(fn)
	T(_v.Kind() != reflect.Func, "func type error(%s)", _v)
}
