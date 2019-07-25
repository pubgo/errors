package internal

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Test struct {
	name string
	desc string
	fn   interface{}
	args []interface{}
}

func (t *Test) In(args ...interface{}) *Test {
	return &Test{fn: t.fn, args: args, desc: t.desc, name: t.name}
}

func (t *Test) _Err(b bool, fn ...interface{}) {
	fmt.Printf("  [Desc func %s] [%s]\n", Green(t.name+" start"), t.desc)
	_err := Try(t.fn)(t.args...)(fn...)
	_isErr := IsNone(_err)

	if (b && !_isErr) || (!b && _isErr) {
		fmt.Printf("  [Desc func %s] --> %s\n", Green(t.name+" ok"), FuncCaller(3))
	} else {
		fmt.Printf("  [Desc func %s] --> %s\n", Red(t.name+" fail"), FuncCaller(3))
	}

	if IsDebug() {
		ErrLog(_err)
	}

	TT((b && _isErr) || (!b && !_isErr), "%s test error", t.name).
		M("input", t.args).
		M("desc", t.desc).
		M("func_name", t.name).
		Done()
}

func (t *Test) IsErr(fn ...interface{}) {
	t._Err(true, fn...)
}

func (t *Test) IsNil(fn ...interface{}) {
	t._Err(false, fn...)
}

func TestRun(fn interface{}, desc func(desc func(string) *Test)) {
	Wrap(AssertFn(reflect.ValueOf(fn)), "func error")

	_name := strings.Split(GetCallerFromFn(reflect.ValueOf(fn)), " ")[1]
	_funcName := strings.Split(GetCallerFromFn(reflect.ValueOf(fn)), " ")[1] + strings.TrimLeft(reflect.TypeOf(fn).String(), "func")
	_path := strings.Split(GetCallerFromFn(reflect.ValueOf(desc)), " ")[0]

	fmt.Printf("[Test func %s] [%s] --> %s\n", Green(_name+" start"), _funcName, _path)
	_err := Try(desc)(func(s string) *Test {
		return &Test{desc: s, fn: fn, name: _name}
	})()
	if _err != nil {
		fmt.Printf("[Test func %s]\n", Red(_name+" fail"))
	} else {
		fmt.Printf("[Test func %s]\n", Green(_name+" success"))
	}
	Panic(_err)
}

func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		fmt.Println(err.Caller(FuncCaller(callDepth)).P())
	})
}

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		fmt.Println(err.P())
	})
}

func Assert() {
	ErrHandle(recover(), func(err *Err) {
		if IsDebug() {
			fmt.Println(err.P())
		}

		os.Exit(1)
	})
}

func Throw(fn interface{}) {
	_fn := reflect.ValueOf(fn)
	T(fn == nil || IsZero(_fn) || _fn.Kind() != reflect.Func, "the input must be func type and not null, input --> %#v",fn)

	ErrHandle(recover(), func(err *Err) {
		panic(err.Caller(GetCallerFromFn(_fn)))
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), func(err *Err) {
		fn(err.Caller(GetCallerFromFn(reflect.ValueOf(fn))))
	})
}

func _handle(err interface{}) *Err {
	if err == nil || IsNone(err) {
		return nil
	}

	if _e, ok := err.(func() (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) (err error)); ok {
		err = _e()
	}

	if _e, ok := err.(func(...interface{}) func(...interface{}) error); ok {
		err = _e()()
	}

	if _e, ok := err.(func(...interface{}) func(...interface{}) func(...interface{}) error); ok {
		err = _e()()()
	}

	if err == nil || IsNone(err) {
		return nil
	}

	m := &Err{}
	switch e := err.(type) {
	case *Err:
		m = e
	case error:
		m.msg = e.Error()
		m.err = e
	case string:
		m.err = errors.New(e)
		m.msg = e
	default:
		m.msg = fmt.Sprintf("unknown type error, input: %#v", e)
		m.err = errors.New("unknown type error")
		m.tag = errTags.UnknownTypeCode
	}
	return m
}

func ErrHandle(err interface{}, fn ...func(err *Err)) {
	if err == nil || IsNone(err) {
		return
	}

	_m := _handle(err)
	if _m == nil || IsNone(_m) || _m.err == nil {
		return
	}

	if len(fn) == 0 {
		return
	}

	Wrap(AssertFn(reflect.ValueOf(fn[0])), "func error")
	fn[0](_m)
}
