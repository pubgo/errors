package errors

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"strconv"
	"time"
)

func TryRaw(fn reflect.Value) func(...reflect.Value) func(...reflect.Value) (err error) {
	assertFn(fn)

	var variadicType reflect.Value
	var isVariadic = fn.Type().IsVariadic()
	if isVariadic {
		variadicType = reflect.New(fn.Type().In(fn.Type().NumIn() - 1).Elem()).Elem()
	}

	_NumIn := fn.Type().NumIn()
	return func(args ...reflect.Value) func(...reflect.Value) (err error) {
		T(isVariadic && len(args) < _NumIn-1, "func input params is error,func(%d) now(%d)", _NumIn, len(args))
		T(!isVariadic && _NumIn != len(args), "func input params is not match,func(%d) now(%d)", _NumIn, len(args))

		for i, k := range args {
			if IsZero(k) {
				args[i] = reflect.New(fn.Type().In(i)).Elem()
				continue
			}

			if isVariadic {
				args[i] = variadicType
			}

			args[i] = k
		}

		return func(cfn ...reflect.Value) (err error) {
			defer func() {
				var m *Err
				if r := recover(); r != nil || !IsZero(reflect.ValueOf(r)) {
					m = new(Err)
					switch d := r.(type) {
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
						m.tag = ErrTags.UnknownTypeCode
					}
				}

				if m == nil || IsZero(reflect.ValueOf(m)) || m.err == nil {
					err = nil
					return
				}
				err = m
			}()

			_c := fn.Call(args)
			if len(cfn) > 0 && !IsZero(cfn[0]) {
				assertFn(cfn[0])
				cfn[0].Call(_c)
			}
			return
		}
	}
}

func Try(fn interface{}) func(...interface{}) func(...interface{}) (err error) {

	_tr := TryRaw(reflect.ValueOf(fn))

	return func(args ...interface{}) func(...interface{}) (err error) {

		var _args = valueGet()
		defer valuePut(_args)

		for _, k := range args {
			_args = append(_args, reflect.ValueOf(k))
		}
		_tr1 := _tr(_args...)

		return func(cfn ...interface{}) (err error) {

			var _cfn = valueGet()
			defer valuePut(_cfn)

			for _, k := range cfn {
				_cfn = append(_cfn, reflect.ValueOf(k))
			}
			return _tr1(_cfn...)
		}
	}
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if err == nil || IsZero(reflect.ValueOf(err)) {
		return
	}

	_m := _handle(err)
	if _m == nil || IsZero(reflect.ValueOf(_m)) {
		return
	}

	if len(fn) == 0 {
		return
	}

	assertFn(reflect.ValueOf(fn[0]))
	_m.caller = append(_m.caller, funcCaller(callDepth))
	fn[0](_m)
}

func Retry(num int, fn func()) (err error) {
	defer Resp(func(_err *Err) {
		err = _err
	})

	T(num < 1, "the num param must be more than 0")

	var all = 0
	var _fn = TryRaw(reflect.ValueOf(fn))
	var _cfn = reflect.Value{}
	for i := 0; i < num; i++ {
		if err = _fn()(_cfn); err == nil {
			return
		}

		all += i
		log.Debug().
			Err(err).
			Str("method", "retry").
			Int("cur_sleep_time", i).
			Int("all_sleep_time", all).
			Msg("")
		time.Sleep(time.Second * time.Duration(i))
	}

	Wrap(err, "retry error,retry_num: "+strconv.Itoa(num))
	return
}

func RetryAt(t time.Duration, fn func(at time.Duration)) {

	var err error
	var all = time.Duration(0)
	var _cfn = reflect.Value{}
	var _fn = TryRaw(reflect.ValueOf(fn))
	for {
		if err = _fn(reflect.ValueOf(all))(_cfn); err == nil {
			return
		}

		all += t
		T(all > Cfg.MaxRetryDur, "more than the max retry duration")
		if _l := log.Debug(); _l.Enabled() {
			_l.Caller().
				Err(err).
				Str("method", "retry_at").
				Float64("cur_retry_time", t.Seconds()).
				Float64("all_retry_time", all.Seconds()).
				Msg("")
		}
		time.Sleep(t)
	}
}

func Ticker(fn func(dur time.Time) time.Duration) {

	var _err error
	var _dur = time.Duration(0)
	var _all = time.Duration(0)
	var _fn = TryRaw(reflect.ValueOf(fn))
	var _cfn = reflect.ValueOf(func(t time.Duration) {
		_dur = t
	})

	for i := 0; ; i++ {
		_err = _fn(reflect.ValueOf(time.Now()))(_cfn)
		if _dur < 0 {
			return
		}

		if _dur == 0 {
			_dur = time.Second
		}

		_all += _dur
		T(_all > Cfg.MaxRetryDur, "more than the max ticker time")
		if _l := log.Debug(); _l.Enabled() && _err != nil {
			_l.Caller().
				Err(_err).
				Str("method", "ticker").
				Int("retry_count", i).
				Float64("retry_all_time", _all.Seconds()).
				Msg("")
		}

		time.Sleep(_dur)
	}
}

func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		err.P()
	})
}
