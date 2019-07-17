package internal

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"strings"
)

type Test struct {
	desc string
	fn   interface{}
	args []interface{}
}

func (t *Test) In(args ...interface{}) *Test {
	return &Test{fn: t.fn, args: args, desc: t.desc}
}

func (t *Test) IsErr(fn ...interface{}) {
	fmt.Printf("[%s] --> %sstart%s\n", t.desc, Red, Reset)
	_err := Try(t.fn)(t.args...)(fn...)
	TT(_err == nil, "[%s] -->%sfail%s", t.desc, Red, Reset).
		M("input", t.args).
		Done()

	if _l := log.Debug(); _l.Enabled() {
		ErrLog(_err)
	}
	fmt.Printf("[%s] --> %sok%s\n\n", t.desc, Red, Reset)
}

func (t *Test) IsNil(fn ...interface{}) {
	fmt.Printf("[%s] --> %sstart%s\n", t.desc, Red, Reset)
	WrapM(Try(t.fn)(t.args...)(fn...), "[%s] -->%sfail%s", t.desc, Red, Reset).
		M("input", t.args).
		Done()
	fmt.Printf("[%s] --> %sok%s\n\n", t.desc, Red, Reset)
}

func TestRun(fn interface{}, desc func(desc func(string) *Test)) {
	defer Assert()

	_funcName := strings.Split(GetCallerFromFn(reflect.ValueOf(fn)), " ")[1]
	_path := strings.Split(GetCallerFromFn(reflect.ValueOf(desc)), " ")[0]
	fmt.Printf("test func [%s] start: %s\n", _funcName, _path)
	Wrap(Try(desc)(func(s string) *Test {
		return &Test{desc: s, fn: fn}
	}), "test error")
	fmt.Printf("test func [%s] %sover%s: %s\n\n", _funcName, Red, Reset, _path)
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
		fmt.Println(err.Caller(FuncCaller(callDepth)).P())
		os.Exit(1)
	})
}

func Throw(fn interface{}) {
	_fn := reflect.ValueOf(fn)
	T(fn == nil || IsZero(_fn) || _fn.Kind() != reflect.Func, "the input must be func type and not null")

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
