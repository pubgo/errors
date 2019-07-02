package errors

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
)

func _Try(fn reflect.Value, args ...reflect.Value) func(value reflect.Value) (err error) {
	_call := FnOf(fn, args...)
	return func(cfn reflect.Value) (err error) {
		defer func() {
			var m *Err
			r := reflect.ValueOf(recover())
			if !IsZero(r) {
				m = new(Err)
				switch d := r.Interface().(type) {
				case *Err:
					m = d
				case error:
					m.err = d
					m.msg = d.Error()
				case string:
					m.err = errors.New(d)
					m.msg = d
				default:
					m.msg = fmt.Sprintf("try type error %#v", d)
					m.err = errors.New(m.msg)
					m.tag = ErrTag.UnknownErr
					_t := r.Type()
					m.m["type"] = _t.String()
					m.m["kind"] = _t.Kind()
					m.m["name"] = _t.Name()
				}
				m.caller = getCallerFromFn(fn)
			}

			if m == nil || IsZero(reflect.ValueOf(m)) || m.err == nil {
				err = nil
				return
			}
			err = m
		}()

		_c := _call()
		if !IsZero(cfn) {
			assertFn(cfn)
			cfn.Call(_c)
		}
		return
	}
}

func Try(fn interface{}, args ...interface{}) func(...interface{}) (err error) {
	_fn := reflect.ValueOf(fn)
	assertFn(_fn)

	var variadicType reflect.Value
	var isVariadic = _fn.Type().IsVariadic()
	if isVariadic {
		variadicType = reflect.New(_fn.Type().In(0)).Elem()
	}

	var _args = make([]reflect.Value, len(args))
	for i, k := range args {
		_v := reflect.ValueOf(k)
		if k == nil || IsZero(_v) {
			if isVariadic {
				args[i] = variadicType
			} else {
				args[i] = reflect.New(_fn.Type().In(i)).Elem()
			}
		}
		_args[i] = _v
	}

	return func(cfn ...interface{}) (err error) {

		var _cfn reflect.Value
		if len(cfn) > 0 && !IsZero(reflect.ValueOf(cfn[0])) {
			_cfn = reflect.ValueOf(cfn[0])
		}

		return _Try(reflect.ValueOf(fn), _args...)(_cfn)
	}
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if err == nil || IsZero(reflect.ValueOf(err)) {
		return
	}

	if _e, ok := err.(func() (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) (err error)); ok {
		err = _e()
	}

	if err == nil || IsZero(reflect.ValueOf(err)) {
		return
	}

	if _e, ok := err.(*Err); ok {
		if len(fn) > 0 {
			assertFn(reflect.ValueOf(fn[0]))
			fn[0](_e)
		}
		return
	}

	if l := log.Debug(); l.Enabled() {
		if _e, ok := err.(error); ok {
			l.Err(_e).Msg("err msg")
			return
		}

		l.Interface("other type", err).
			Bool("is zero", IsZero(reflect.ValueOf(err)) || err == nil).
			Str("Kind", reflect.TypeOf(err).String()).
			Msgf("%#v", err)
	}
}

func Retry(num int, fn func()) {
	defer Handle(func() {})

	T(num < 1, "the num param must be more than 0")

	var err error
	var all = 0
	for i := 0; i < num; i++ {
		if err = Try(fn)(); err == nil {
			return
		}

		all += i
		log.Warn().Caller().Str("method", "retry").Int("cur_sleep_time", i).Int("all_sleep_time", all).Msg(err.Error())
		time.Sleep(time.Second * time.Duration(i))
	}

	Wrap(err, "retry error,retry_num(%d)", num)
}

func RetryAt(t time.Duration, fn func(at time.Duration)) {
	defer Handle(func() {})

	var err error
	var all = time.Duration(0)
	for {
		if err = Try(fn, all)(); err == nil {
			return
		}

		all += t
		log.Warn().Caller().Str("method", "retry_at").Float64("cur_sleep_time", t.Seconds()).Float64("all_sleep_time", all.Seconds()).Msg(err.Error())
		time.Sleep(t)
	}
}

func Ticker(fn func(dur time.Time) time.Duration) {
	defer Handle(func() {})

	var _dur = time.Duration(0)
	for i := 0; ; i++ {
		ErrHandle(Try(func() {
			_dur = fn(time.Now())
		}), func(err *Err) {
			err.Log()
		})

		if _dur < 0 {
			return
		}

		if _dur == 0 {
			_dur = time.Second
		}

		time.Sleep(_dur)
	}
}

func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		err.P()
	})
}
