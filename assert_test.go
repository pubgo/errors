package errors_test

import (
	es "errors"
	"fmt"
	"github.com/pubgo/errors"
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

	t.Log(errors.If(true, "", "").(string))
}

func TestT(t *testing.T) {
	defer errors.Log()

	errors.T(true, "test t")
}

func TestTT(t *testing.T) {
	defer errors.Log()

	errors.TT(true, func(m *errors.M) {
		m.Msg("test tt").
			Tag("tag").
			M("k", "v")

	})
}

func TestWrap(t *testing.T) {
	defer errors.Log()

	errors.Wrap(es.New("test"), "test")
}

func TestWrapM(t *testing.T) {
	defer errors.Log()

	errors.WrapM(es.New("dd"), func(m *errors.M) {
		m.Msg("test")
	})
}

func testFunc() {
	defer errors.Handle(func() {})

	errors.WrapM(es.New("sbhbhbh"), func(m *errors.M) {
		m.Msg("test shhh").
			M("ss", 1).
			M("input", 2)
	})
}

func TestPanic(t *testing.T) {
	defer errors.Log()

	for i := 0; i < 10000; i++ {
		errors.Try(testFunc)()
		t.Log("ok")
	}

}

func init11() {
	defer errors.Handle(func() {})

	//T(true, "sss")
	errors.TT(true, func(m *errors.M) {
		m.Msg("test tt")
	})
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
