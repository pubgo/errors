package errors

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

func Try(fn interface{}, args ...interface{}) func(...interface{}) (err error) {
	return func(cfn ...interface{}) (err error) {
		defer func() {

			var m *Err
			if r := recover(); !IsZero(r) {
				m = new(Err)
				caller := getCallerFromFn(fn)
				switch d := r.(type) {
				case *Err:
					m = d
					m.caller = caller
				case error:
					m.err = d
					m.msg = d.Error()
					m.caller = caller
				case string:
					m.err = errors.New(d)
					m.msg = d
					m.caller = caller
				default:
					m.msg = fmt.Sprintf("try type error %#v", d)
					m.err = errors.New(m.msg)
					m.caller = caller
					m.tag = ErrTag.UnknownErr
					_t := reflect.TypeOf(err)
					m.m["type"] = _t.String()
					m.m["kind"] = _t.Kind()
					m.m["name"] = _t.Name()
				}
			}

			if IsZero(m) || m.err == nil {
				err = nil
				return
			}
			err = m
		}()

		_call := fnOf(fn, args...)
		if len(cfn) == 0 {
			_call()
			return
		}

		assertFn(cfn[0])
		reflect.ValueOf(cfn[0]).Call(_call())
		return
	}
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if IsZero(err) {
		return
	}

	if _e, ok := err.(func() (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) (err error)); ok {
		err = _e()
	}

	if IsZero(err) {
		return
	}

	if _e, ok := err.(*Err); ok {
		if len(fn) > 0 {
			assertFn(fn[0])
			fn[0](_e)
		}
		return
	}

	if _e, ok := err.(error); ok {
		log.Error().Err(_e).Msg("ErrHandle")
		return
	}

	fmt.Printf("other type: %#v\n", err)
	fmt.Printf("is zero: %#v\n", IsZero(err) || err == nil)
	fmt.Printf("Kind: %#v\n", reflect.TypeOf(err).Kind())
	fmt.Printf("String: %#v\n", reflect.TypeOf(err).String())

}
