# errors
error handle for go

```go
package errors_test

import (
	es "errors"
	"fmt"
	"github.com/pubgo/errors"
	"github.com/rs/zerolog/log"
	"testing"
	"time"
)

func TestErrLog(t *testing.T) {
	defer errors.Log()

	errors.T(true, "test t")
}

func TestRetry(t *testing.T) {
	defer errors.Log()

	errors.Retry(5, func() {
		errors.T(true, "test t")
	})
}

func TestIf(t *testing.T) {
	defer errors.Log()

	log.Info().Msg(errors.If(true, "test true", "test false").(string))
}

func TestT(t *testing.T) {
	defer errors.Log()

	errors.T(true, "test t")
}

func TestTT(t *testing.T) {
	defer errors.Log()

	errors.TT(true, "test tt").
		M("k", "v").
		M("tag", "tag").
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

func testFunc_1() {
	defer errors.Handle(func() {})

	errors.WrapM(es.New("sbhbhbh"), "test shhh").
		M("ss", 1).
		M("input", 2).
		Done()
}

func testFunc() {
	defer errors.Handle(func() {})

	errors.ErrLog(errors.Try(testFunc_1))
}

func TestPanic(t *testing.T) {
	defer errors.Debug()

	errors.ErrLog(errors.Try(testFunc)())
}

func init11() {
	defer errors.Handle(func() {})

	//T(true, "sss")
	errors.T(true, "test tt")
}

func TestT2(t *testing.T) {
	defer errors.Log()

	init11()
}

func TestTry(t *testing.T) {
	defer errors.Log()

	errors.Panic(errors.Try(errors.T, true, "sss"))
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
		defer errors.Handle(func() {})

		errors.Wrap(es.New("hello error"), "sss")
	}()

}

func TestErrHandle(t *testing.T) {
	defer errors.Log()

	errors.ErrHandle(errors.Try(errors.T, true, "test 12345"), func(err *errors.Err) {
		err.P()
	})

	errors.ErrHandle(errors.Try(errors.T, true, "test 12345"))

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
	errors.T(errors.IsZero(1), "")
	errors.T(errors.IsZero(1.2), "")
	errors.T(!errors.IsZero(nil), "")
	errors.T(errors.IsZero("ss"), "")
	errors.T(errors.IsZero(map[string]interface{}{}), "")
	errors.T(errors.IsZero(ss()), "")
	errors.T(!errors.IsZero(ss1()), "")
	errors.T(errors.IsZero(&s), "")
	errors.T(!errors.IsZero(ss2), "")
}

func TestResp(t *testing.T) {
	defer errors.Resp(func(err *errors.Err) {
		err.Log()
	})

	errors.T(true, "data handle")
}

func TestTicker(t *testing.T) {
	defer errors.Handle(func() {})

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
```