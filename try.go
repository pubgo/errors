package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
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

		_call := FnOf(fn, args...)
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

	if l := log.Debug(); l.Enabled() {
		if _e, ok := err.(error); ok {
			l.Err(_e).Msg("err msg")
			return
		}

		l.Interface("other type", err).
			Bool("is zero", IsZero(err) || err == nil).
			Str("Kind", reflect.TypeOf(err).String()).
			Msgf("%#v",err)
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
			if dt, err := json.MarshalIndent(i, "", "\t"); err != nil {
				log.Error().Caller().Err(err).Msg("json MarshalIndent error")
			} else {
				log.Warn().Caller().Msg(string(dt))
			}
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
