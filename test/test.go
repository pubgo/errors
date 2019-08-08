package test

import (
	"fmt"
	"github.com/pubgo/errors/internal"
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
	fmt.Printf("  [Desc func %s] [%s]\n", internal.Green(t.name+" start"), t.desc)
	_err := internal.Try(t.fn)(t.args...)(fn...)
	_isErr := internal.IsNone(_err)

	if (b && !_isErr) || (!b && _isErr) {
		fmt.Printf("  [Desc func %s] --> %s\n", internal.Green(t.name+" ok"), internal.FuncCaller(3))
	} else {
		fmt.Printf("  [Desc func %s] --> %s\n", internal.Red(t.name+" fail"), internal.FuncCaller(3))
	}

	if internal.IsDebug() {
		internal.ErrLog(_err)
	}

	internal.TT((b && _isErr) || (!b && !_isErr), func(err *internal.Err) {
		err.Msg("%s test error", t.name)
		err.M("input", t.args)
		err.M("desc", t.desc)
		err.M("func_name", t.name)
	})

}

func (t *Test) IsErr(fn ...interface{}) {
	t._Err(true, fn...)
}

func (t *Test) IsNil(fn ...interface{}) {
	t._Err(false, fn...)
}

func Run(fn interface{}, desc func(desc func(string) *Test)) {
	internal.Wrap(internal.AssertFn(reflect.ValueOf(fn)), "func error")

	_name := strings.Split(internal.GetCallerFromFn(reflect.ValueOf(fn)), " ")[1]
	_funcName := strings.Split(internal.GetCallerFromFn(reflect.ValueOf(fn)), " ")[1] + strings.TrimLeft(reflect.TypeOf(fn).String(), "func")
	_path := strings.Split(internal.GetCallerFromFn(reflect.ValueOf(desc)), " ")[0]

	fmt.Printf("[Test func %s] [%s] --> %s\n", internal.Green(_name+" start"), _funcName, _path)
	_err := internal.Try(desc)(func(s string) *Test {
		return &Test{desc: s, fn: fn, name: _name}
	})()
	if _err != nil {
		fmt.Printf("[Test func %s]\n", internal.Red(_name+" fail"))
	} else {
		fmt.Printf("[Test func %s]\n", internal.Green(_name+" success"))
	}
	internal.Panic(_err)
}
