package errors

import "github.com/pubgo/errors/internal"

// Err
type Err = internal.Err

// error assert
var Panic = internal.Panic
var Wrap = internal.Wrap
var WrapM = internal.WrapM
var TT = internal.TT
var T = internal.T

// error handle
var Throw = internal.Throw
var Assert = internal.Assert
var Resp = internal.Resp
var ErrLog = internal.ErrLog
var ErrHandle = internal.ErrHandle
var Debug = internal.Debug

// test
type Test = internal.Test

var TestRun = internal.TestRun

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


// try
var Try = internal.Try
var Retry = internal.Retry
var RetryAt = internal.RetryAt
var Ticker = internal.Ticker
