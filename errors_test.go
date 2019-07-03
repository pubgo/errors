package errors_test

import (
	es "errors"
	"fmt"
	"github.com/pubgo/errors"
	"reflect"
	"testing"
	"time"
)

func TestT(t *testing.T) {
	defer errors.Debug()

	errors.T(true, "test t")
}

func TestRetry(t *testing.T) {
	defer errors.Debug()

	errors.Wrap(errors.Retry(3, func() {
		errors.T(true, "test t")
	}), "test Retry error")
}

func TestIf(t *testing.T) {
	defer errors.Log()

	errors.T(errors.If(true, "test true", "test false").(string) != "test true", "")
}

func TestTT(t *testing.T) {
	defer errors.Debug()

	errors.TT(true, "test tt").
		M("k", "v").
		SetTag("ss", 12).
		Done()
}

func TestWrap(t *testing.T) {
	defer errors.Debug()

	errors.Wrap(es.New("test"), "test")
}

func TestWrapM(t *testing.T) {
	defer errors.Debug()

	errors.WrapM(es.New("dd"), "test").
		Done()
}

func testFunc_2() {
	defer errors.Handle()()

	errors.WrapM(es.New("testFunc_1"), "test shhh").
		M("ss", 1).
		M("input", 2).
		Done()
}

func testFunc_1() {
	defer errors.Handle()()

	testFunc_2()
}

func testFunc() {
	defer errors.Handle()()

	errors.Wrap(errors.Try(testFunc_1)()(), "errors.Wrap")
}

func TestErrLog(t *testing.T) {
	defer errors.Debug()

	errors.Wrap(errors.Try(testFunc)()(), "errors.Wrap11111")
}

func init11() {
	defer errors.Handle()()

	//T(true, "sss")
	errors.T(true, "test tt")
}

func TestT2(t *testing.T) {
	defer errors.Debug()

	init11()
}

func TestTry(t *testing.T) {
	defer errors.Log()

	errors.Panic(errors.Try(errors.T)(true, "sss"))
}

func TestTask(t *testing.T) {
	defer errors.Log()

	errors.Wrap(errors.Try(func() {
		errors.Wrap(es.New("dd"), "err ")
	}), "test wrap")
}

func TestHandle(t *testing.T) {
	defer errors.Log()

	func() {
		defer errors.Handle()()

		errors.Wrap(es.New("hello error"), "sss")
	}()

}

func TestErrHandle(t *testing.T) {
	defer errors.Log()

	errors.ErrHandle(errors.Try(errors.T)(true, "test 12345"), func(err *errors.Err) {
		err.P()
	})

	errors.ErrHandle(errors.Try(errors.T)(true, "test 12345"))

	errors.ErrHandle("ttt")
	errors.ErrHandle(es.New("eee"))
	errors.ErrHandle([]string{"dd"})
}

func TestIsZero(t *testing.T) {
	defer errors.Log()

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
	defer errors.Resp(func(err *errors.Err) {
		err.Log()
	})

	errors.T(true, "data handle")
}

func TestTicker(t *testing.T) {
	defer errors.Handle()()

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

func _2nit(a string, m ...int) {

}
func TestNam1e(t *testing.T) {
	_fn := reflect.ValueOf(_2nit)
	fmt.Println(_fn.Type().IsVariadic())
	fmt.Println(_fn.Type().NumIn())

	fmt.Println(_fn.Type().In(_fn.Type().NumIn() - 1).Elem())

	//var variadicType reflect.Value
	//var isVariadic = _fn.Type().IsVariadic()
	//if isVariadic {
	//	variadicType = reflect.New(_fn.Type().In(0)).Elem()
	//}
}
func TestErr(t *testing.T) {
	errors.ErrHandle(errors.Try(func() {
		errors.ErrHandle(errors.Try(func() {
			errors.T(true, "90999 error")
		}), func(err *errors.Err) {
			errors.Wrap(err, "wrap")
		})
	}), func(err *errors.Err) {
		err.P()
	})
}
