package errors

import (
	"fmt"
	"github.com/pubgo/errors/internal"
	"os"
	"reflect"
	"runtime/debug"
)

// Err
type Err = internal.Err

// error assert
var Panic = internal.Panic
var Wrap = internal.Wrap
var WrapM = func(err interface{}, fn func(err *Err)) {
	internal.WrapM(err, fn)
}
var TT = func(b bool, fn func(err *Err)) {
	internal.TT(b, fn)
}
var T = internal.T

// error handle
var Throw = internal.Throw

func Assert() {
	ErrHandle(recover(), func(err *Err) {
		if internal.IsDebug() {
			fmt.Println(err.P())
			debug.PrintStack()
		}
		os.Exit(1)
	})
}

func Resp(fn func(err *Err)) {
	ErrHandle(recover(), func(err *Err) {
		err.Caller(GetCallerFromFn(reflect.ValueOf(fn)))
		fn(err)
	})
}

func RespErr(err *error) {
	ErrHandle(recover(), func(_err *Err) {
		*err = _err
	})
}

var ErrLog = internal.ErrLog
var ErrHandle = internal.ErrHandle
var Debug = internal.Debug

// config
var Cfg = &internal.Cfg

// err tag
var ErrTagRegistry = internal.ErrTagRegistry
var ErrTags = internal.ErrTags
var ErrTagsMatch = internal.ErrTagsMatch

// utils
var AssertFn = internal.AssertFn
var If = internal.If
var IsZero = internal.IsZero
var IsNone = internal.IsNone
var P = internal.P
var FuncCaller = internal.FuncCaller
var GetCallerFromFn = internal.GetCallerFromFn
var LoadEnvFile = internal.LoadEnvFile
var InitDebug = internal.InitDebug

// try
var Try = internal.Try
var Retry = internal.Retry
var RetryAt = internal.RetryAt
var Ticker = internal.Ticker
