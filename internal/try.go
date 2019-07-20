package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

func TryRaw(fn reflect.Value) func(...reflect.Value) func(...reflect.Value) (err error) {
	Wrap(AssertFn(fn), "func error")

	var variadicType reflect.Value
	var isVariadic = fn.Type().IsVariadic()
	if isVariadic {
		variadicType = reflect.New(fn.Type().In(fn.Type().NumIn() - 1).Elem()).Elem()
	}

	_NumIn := fn.Type().NumIn()
	return func(args ...reflect.Value) func(...reflect.Value) (err error) {
		T(isVariadic && len(args) < _NumIn-1, "func %s input params is error,func(%d,%d)", fn.Type(), _NumIn, len(args))
		T(!isVariadic && _NumIn != len(args), "func %s input params is not match,func(%d,%d)", fn.Type(), _NumIn, len(args))

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

		_call := FuncCaller(3)
		return func(cfn ...reflect.Value) (err error) {
			defer func() {
				ErrHandle(recover(), func(_err *Err) {
					err = _err.Caller(_call)
				})
			}()

			_c := fn.Call(args)
			if len(cfn) > 0 && !IsZero(cfn[0]) {
				Wrap(AssertFn(cfn[0]), "func error")
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

func Retry(num int, fn func()) (err error) {
	defer Resp(func(_err *Err) {
		err = _err.Caller(GetCallerFromFn(reflect.ValueOf(fn)))
	})

	T(num < 1, "the num is less than 0")

	var all = 0
	var _fn = TryRaw(reflect.ValueOf(fn))
	for i := 0; i < num; i++ {
		if err = _fn()(); err == nil {
			return
		}

		all += i
		if IsDebug() {
			fmt.Printf("cur_sleep_time: %d, all_sleep_time: %d", i, all)
			ErrLog(err)
		}
		time.Sleep(time.Second * time.Duration(i))
	}

	Wrap(err, "retry error,retry_num: "+strconv.Itoa(num))
	return
}

func RetryAt(t time.Duration, fn func(at time.Duration)) {
	defer Throw(fn)

	var err error
	var all = time.Duration(0)
	var _fn = TryRaw(reflect.ValueOf(fn))
	for {
		if err = _fn(reflect.ValueOf(all))(); err == nil {
			return
		}

		all += t
		if all > Cfg.MaxRetryDur {
			T(true, "more than the max(%s) retry duration", Cfg.MaxRetryDur.String())
		}

		if IsDebug() {
			fmt.Printf("cur_retry_time: %d, all_retry_time: %f", t.Seconds(), all.Seconds())
			ErrLog(err)
		}
		time.Sleep(t)
	}
}

func Ticker(fn func(dur time.Time) time.Duration) {
	defer Throw(fn)

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
		if IsDebug() {
			fmt.Printf("retry_count: %d, retry_all_time: %f", i, _all.Seconds())
			ErrLog(_err)
		}
		time.Sleep(_dur)
	}
}

var _valuePool = sync.Pool{
	New: func() interface{} {
		return []reflect.Value{}
	},
}

func valueGet() []reflect.Value {
	return _valuePool.Get().([]reflect.Value)
}

func valuePut(v []reflect.Value) {
	v = v[:0]
	_valuePool.Put(v)
}
