package errors_test

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
	errors.P("errors.Cfg", errors.Cfg)
}

func TestT(t *testing.T) {
	defer errors.Assert()
	errors.T(true, "test t")
}

func TestErrLog2(t *testing.T) {
	defer errors.Resp(func(err *errors.Err) {
		errors.ErrLog(err)
	})

	errors.T(true, "test t")
}

func TestDebug(t *testing.T) {
	defer errors.Assert()

	errors.T(true, "test t")
}

func TestRetry(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(internal.Retry(3, func() {
		errors.T(true, "test t")
	}), "test Retry error")
}

func TestIf(t *testing.T) {
	defer errors.Assert()

	errors.T(errors.If(true, "test true", "test false").(string) != "test true", "")
}

func TestTT(t *testing.T) {
	defer errors.Assert()

	errors.TT(true, "test tt").
		M("k", "v").
		SetTag("12").
		Done()
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
	errors.Wrap(internal.Try(testFunc_1), "errors.Wrap")
}

func TestErrLog(t *testing.T) {
	defer errors.Assert()

	testFunc()
}

func init11() {
	errors.T(true, "test tt")
}

func TestT2(t *testing.T) {
	defer errors.Assert()

	init11()
}

func TestTry(t *testing.T) {
	defer errors.Assert()

	errors.Panic(internal.Try(errors.T)(true, "sss"))
}

func TestTask(t *testing.T) {
	defer errors.Assert()

	errors.Wrap(internal.Try(func() {
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

	errors.ErrHandle(internal.Try(func() {
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

	errors.T(true, "data handle")
}

func TestTicker(t *testing.T) {
	defer errors.Assert()

	internal.Ticker(func(dur time.Time) time.Duration {
		fmt.Println(dur)
		return time.Second
	})
}

func TestRetryAt(t *testing.T) {
	internal.RetryAt(time.Second*2, func(dur time.Duration) {
		fmt.Println(dur.String())

		errors.T(true, "test RetryAt")
	})
}

func TestErr(t *testing.T) {
	errors.ErrHandle(internal.Try(func() {
		errors.ErrHandle(internal.Try(func() {
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
	errors.TestRun("AssertFn", errors.AssertFn, func(t *errors.Test) {
		t.In(reflect.ValueOf(func() {})).IsNil()
		t.In(reflect.ValueOf(nil)).IsErr()
	})

	fn()
}

func TestGetCallerFromFn(t *testing.T) {
	defer errors.Assert()

	errors.TestRun("GetCallerFromFn", _GetCallerFromFn1, func(t *errors.Test) {
		t.In(_GetCallerFromFn2).IsNil()
		t.In(nil).IsErr()
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

	errors.TestRun("AssertFn", errors.AssertFn, func(t *errors.Test) {
		t.In(reflect.ValueOf(func() {})).IsNil()
		t.In(reflect.ValueOf(nil)).IsErr()
	})
}
