package errors

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"reflect"
)

func T(b bool, msg string, args ...interface{}) {
	if !b {
		return
	}

	if _l := log.Debug(); _l.Enabled() {
		_l.Msgf(msg, args...)
	}

	TT(b, func(m *M) {
		m.Msg(msg, args...)
		m.Caller(5)
	})
}

func TT(b bool, fn func(m *M)) {
	if !b {
		return
	}

	_m := newM()
	fn(&_m)

	if len(_m.m) == 0 {
		_m.m = nil
	}

	if _l := log.Debug(); _l.Enabled() {
		_l.Msg(_m.msg)
	}

	panic(&Err{
		caller: If(_m.caller != "", _m.caller, funcCaller(callDepth)).(string),
		msg:    _m.msg,
		err:    errors.New(_m.msg),
		tag:    _m.tag,
		m:      _m.m,
	})
}

func WrapM(err interface{}, fn func(m *M)) {
	if IsZero(err) {
		return
	}

	m := _handle(err)
	if IsZero(m) {
		return
	}

	_m := newM()
	fn(&_m)

	if len(_m.m) == 0 {
		_m.m = nil
	}

	if _l := log.Debug(); _l.Enabled() {
		_l.Msg(_m.msg)
	}

	panic(&Err{
		sub:    m,
		caller: If(_m.caller != "", _m.caller, funcCaller(callDepth)).(string),
		msg:    _m.msg,
		err:    m.tErr(),
		tag:    m.tTag(_m.tag),
		m:      _m.m,
	})
}

func Wrap(err interface{}, msg string, args ...interface{}) {
	if IsZero(err) {
		return
	}

	if _l := log.Debug(); _l.Enabled() {
		_l.Msgf(msg,args...)
	}

	WrapM(err, func(m *M) {
		m.Msg(msg, args...)
		m.Caller(5)
	})
}

func Panic(err interface{}) {
	if IsZero(err) {
		return
	}

	if _l := log.Debug(); _l.Enabled() {
		_l.Msg(err.(error).Error())
	}

	WrapM(err, func(m *M) {
		m.Caller(5)
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
	T(_v.Kind() != reflect.Func, "func type error(%s)", _v.String())
}
