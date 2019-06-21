package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
)

func T(b bool, msg string, args ...interface{}) {
	if !b {
		return
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
	fn(_m)

	if len(_m.m) == 0 {
		_m.m = nil
	}

	if Cfg.Debug {
		log.Println(_m.msg, funcCaller(callDepth))
	}

	panic(&Err{
		caller: If(_m.caller == "", funcCaller(callDepth), _m.caller).(string),
		msg:    _m.msg,
		err:    errors.New(_m.msg),
		tag:    _m.tag,
		m:      _m.m,
	})
}

func WrapM(err error, fn func(m *M)) {
	if IsNil(err) {
		return
	}

	m := _handle(err)
	_m := newM()
	fn(_m)

	if len(_m.m) == 0 {
		_m.m = nil
	}

	if Cfg.Debug {
		log.Println(_m.msg, funcCaller(callDepth))
	}

	panic(&Err{
		sub:    m,
		caller: If(_m.caller == "", funcCaller(callDepth), _m.caller).(string),
		msg:    _m.msg,
		err:    m.tErr(),
		tag:    m.tTag(_m.tag),
		m:      _m.m,
	})
}



func Wrap(err error, msg string, args ...interface{}) {
	if IsNil(err) {
		return
	}

	WrapM(err, func(m *M) {
		m.Msg(msg, args...)
		m.Caller(5)
	})
}

func Panic(err error) {
	if IsNil(err) {
		return
	}

	WrapM(err, func(m *M) {
		m.Caller(5)
	})
}

func P(d ...interface{}) {
	defer Handle(func(m *M) {
	})

	for _, i := range d {
		if IsNil(i) {
			return
		}

		if dt, err := json.MarshalIndent(i, "", "\t"); err != nil {
			Panic(err)
		} else {
			fmt.Println(reflect.ValueOf(i).String(), "->", string(dt))
		}
	}
}
