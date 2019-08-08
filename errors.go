package errors

import (
	"github.com/pubgo/errors/internal"
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
var Assert = internal.Assert
var Resp = func(fn func(err *Err)) {
	internal.Resp(fn)
}
var RespErr = internal.RespErr
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
