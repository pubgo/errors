package tests_test

import (
	es "errors"
	"fmt"
	"github.com/pubgo/errors"
	"github.com/pubgo/errors/internal"
	"reflect"
	"testing"
	"time"
)

func init() {
	internal.InitDebug()
}

func TestCfg(t *testing.T) {
	errors.P("errors.Cfg", errors.Cfg)
}

func TestT(t *testing.T) {
	errors.TestRun(errors.T, func(desc func(string) *errors.Test) {
		desc("params is true").In(true, "test t").IsErr()
		desc("params is false").In(false, "test t").IsNil()
	})
}

func TestErrLog2(t *testing.T) {
	errors.TestRun(errors.ErrLog, func(desc func(string) *errors.Test) {
		desc("err log params").In(es.New("sss")).IsNil()
		desc("nil params").In(es.New("sss")).IsNil()
	})
}

func TestRetry(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.Retry, func(desc func(string) *errors.Test) {
		desc("retry(3)").In(3, func() {
			errors.T(true, "test t")
		}).IsErr(func(err error) {
			errors.Wrap(err, "test Retry error")
		})
	})
}

func TestIf(t *testing.T) {
	defer errors.Assert()

	errors.T(errors.If(true, "test true", "test false").(string) != "test true", "")
}

func TestTT(t *testing.T) {
	defer errors.Assert()

	_fn := func(b bool) {
		errors.TT(b, "test tt").M("k", "v").SetTag("12").Done()
	}

	errors.TestRun(_fn, func(desc func(string) *errors.Test) {
		desc("true params 1").In(true).IsErr()
		desc("true params 2").In(true).IsErr()
		desc("true params 3").In(true).IsErr()
		desc("false params").In(false).IsNil()
	})
}

func TestWrap(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(es.New("test"), "test")
}

func TestWrapM(t *testing.T) {
	defer errors.Assert()

	errors.WrapM(es.New("dd"), "test").
		Done()
}

func testFunc_2() {
	errors.WrapM(es.New("testFunc_1"), "test shhh").
		M("ss", 1).
		M("input", 2).
		Done()
}

func testFunc_1() {
	testFunc_2()
}

func testFunc() {
	errors.Wrap(errors.Try(testFunc_1), "errors.Wrap")
}

func TestErrLog(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(testFunc, func(desc func(string) *errors.Test) {
		desc("test func").In().IsErr()
	})
}

func init11() {
	errors.T(true, "test tt")
}

func TestT2(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(init11, func(desc func(string) *errors.Test) {
		desc("simple test").In().IsErr()
	})
}

func TestTry(t *testing.T) {
	defer errors.Assert()

	errors.Panic(errors.Try(errors.T)(true, "sss"))
}

func TestTask(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(errors.Try(func() {
		errors.Wrap(es.New("dd"), "err ")
	}), "test wrap")
}

func TestHandle(t *testing.T) {
	defer errors.Assert()

	func() {
		errors.Wrap(es.New("hello error"), "sss")
	}()

}

func TestErrHandle(t *testing.T) {
	defer errors.Assert()

	errors.ErrHandle(errors.Try(func() {
		errors.T(true, "test T")
	}), func(err *errors.Err) {
		err.P()
	})

	errors.ErrHandle("ttt", func(err *errors.Err) {
		err.P()
	})
	errors.ErrHandle(es.New("eee"), func(err *errors.Err) {
		err.P()
	})
	errors.ErrHandle([]string{"dd"}, func(err *errors.Err) {
		err.P()
	})
}

func TestIsZero(t *testing.T) {
	//defer errors.Log()

	var ss = func() map[string]interface{} {
		return make(map[string]interface{})
	}

	var ss1 = func() map[string]interface{} {
		return nil
	}

	var s = 1
	var ss2 map[string]interface{}
	errors.T(errors.IsZero(reflect.ValueOf(1)), "")
	errors.T(errors.IsZero(reflect.ValueOf(1.2)), "")
	errors.T(!errors.IsZero(reflect.ValueOf(nil)), "")
	errors.T(errors.IsZero(reflect.ValueOf("ss")), "")
	errors.T(errors.IsZero(reflect.ValueOf(map[string]interface{}{})), "")
	errors.T(errors.IsZero(reflect.ValueOf(ss())), "")
	errors.T(!errors.IsZero(reflect.ValueOf(ss1())), "")
	errors.T(errors.IsZero(reflect.ValueOf(&s)), "")
	errors.T(!errors.IsZero(reflect.ValueOf(ss2)), "")
}

func TestResp(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.Resp, func(desc func(string) *errors.Test) {
		desc("resp ok").In(func(err *errors.Err) {
			err = err.Caller(errors.FuncCaller(2))
		}).IsNil()
	})

}

func TestTicker(t *testing.T) {
	defer errors.Assert()

	errors.Ticker(func(dur time.Time) time.Duration {
		fmt.Println(dur)
		return time.Second
	})
}

func TestRetryAt(t *testing.T) {
	errors.RetryAt(time.Second*2, func(dur time.Duration) {
		fmt.Println(dur.String())

		errors.T(true, "test RetryAt")
	})
}

func TestErr(t *testing.T) {
	errors.ErrHandle(errors.Try(func() {
		errors.ErrHandle(errors.Try(func() {
			errors.T(true, "90999 error")
		}), func(err *errors.Err) {
			errors.Wrap(err, "wrap")
		})
	}), func(err *errors.Err) {
		fmt.Println(err.P())
	})
}

func _GetCallerFromFn2() {
	errors.WrapM(es.New("test 123"), "test GetCallerFromFn").
		M("ss", "dd").
		Done()
}

func _GetCallerFromFn1(fn func()) {
	fn()
}

func TestGetCallerFromFn(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(_GetCallerFromFn1, func(desc func(string) *errors.Test) {
		desc("GetCallerFromFn ok").In(_GetCallerFromFn2).IsErr()
		desc("GetCallerFromFn nil").In(nil).IsErr()
	})
}

func TestErrTagRegistry(t *testing.T) {
	defer errors.Assert()

	errors.ErrTagRegistry("errors_1")
	errors.ErrTagRegistry("errors_2")
	fmt.Printf("%#v\n", errors.ErrTags())

	errors.T(errors.ErrTagsMatch("errors") == true, "errors match error")
	errors.T(errors.ErrTagsMatch("errors_1") == false, "errors_1 not match")
}

func TestTest(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.AssertFn, func(desc func(string) *errors.Test) {
		desc("params is func 1").
			In(reflect.ValueOf(func() {})).
			IsNil(func(err error) {
				errors.Wrap(err, "check error")
			})

		desc("params is func 2").
			In(reflect.ValueOf(func() {})).
			IsNil(func(err error) {
				errors.Wrap(err, "check error")
			})

		desc("params is func 3").
			In(reflect.ValueOf(func() {})).
			IsNil(func(err error) {
				errors.Wrap(err, "check error")
			})

		desc("params is nil").
			In(reflect.ValueOf(nil)).
			IsErr(func(err error) {
				errors.Wrap(err, "check error ok")
			})
	})
}

func TestThrow(t *testing.T) {
	defer errors.Assert()

	errors.TestRun(errors.Throw, func(desc func(string) *errors.Test) {
		desc("not func type params").In(es.New("ss")).IsErr()
		desc("func type params").In(func() {}).IsNil()
		desc("nil type params").In(nil).IsErr()
	})
}

type a1 struct {
	name string
}

func F1(name string) func(func(err *errors.Err)) *a1 {
	return func(i func(err *errors.Err)) *a1 {
		defer errors.Resp(func(err *internal.Err) {
			i(err)
		})

		errors.T(name == "", "name is null")
		return &a1{name: name}
	}
}

func F2(name string) func(func(err *errors.Err)) *a1 {
	return func(i func(err *errors.Err)) *a1 {
		defer errors.Resp(i)

		return F1(name)(func(err *errors.Err) {
			panic(err)
		})
	}
}

func TestWrapCall(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(errors.Try(F2)("")(func(a *a1) {

	}), "")

	f := F2("")(func(err *errors.Err) {
		fmt.Println(err)
	})
	fmt.Println(f.name)
}
