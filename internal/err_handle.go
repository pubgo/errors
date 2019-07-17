package internal

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

type Test struct {
	desc string
	fn   interface{}
	args []interface{}
}

func (t *Test) IsErr(fn ...interface{}) {
	_err := Try(t.fn)(t.args...)(fn...)
	TT(_err == nil, "test func %s fail", FuncCaller(2)).
		M("input", t.args).
		Done()
	fmt.Printf("test func %s ok\n", FuncCaller(2))
}

func (t *Test) In(args ...interface{}) *Test {
	return &Test{fn: t.fn, args: args, desc: t.desc}
}

func (t *Test) IsNil(fn ...interface{}) {
	WrapM(Try(t.fn)(t.args...)(fn...), "test func %s fail", FuncCaller(2)).
		M("input", t.args).
		Done()
	fmt.Printf("test func %s ok\n", FuncCaller(2))
}

func TestRun(desc string, fn interface{}, t func(t *Test)) {
	fmt.Println(desc,"start")
	t(&Test{fn: fn, desc: desc})
	fmt.Printf("%s over\n\n", desc)
}

func ErrLog(err interface{}) {
	ErrHandle(err, func(err *Err) {
		err.caller = append(err.caller[:len(err.caller)-1], FuncCaller(callDepth))
		fmt.Println(err.P())
	})
}

func Debug() {
	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		fmt.Println(err.P())
	})
}

func Assert() {
	ErrHandle(recover(), func(err *Err) {
		//err.caller = err.caller[:len(err.caller)-1]
		fmt.Println(err.P())
		os.Exit(1)
	})
}

func Throw(fn interface{}) {
	_fn := reflect.ValueOf(fn)
	T(fn == nil || IsZero(_fn) || _fn.Kind() != reflect.Func, "the input must be func type and not null")

	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		err.caller = append(err.caller, GetCallerFromFn(_fn))
		panic(err)
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), func(err *Err) {
		err.caller = err.caller[:len(err.caller)-1]
		err.caller = append(err.caller, GetCallerFromFn(reflect.ValueOf(fn)))
		fn(err)
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

	AssertFn(reflect.ValueOf(fn[0]))
	_m.caller = append(_m.caller, FuncCaller(callDepth))
	fn[0](_m)
}
