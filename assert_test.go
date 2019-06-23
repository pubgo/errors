package errors_test

import (
	es "errors"
	"github.com/pkg/profile"
	"github.com/pubgo/errors"
	"testing"
)

func TestT(t *testing.T) {
	defer errors.Debug()

	errors.T(true, "test t")
}

func TestTT(t *testing.T) {
	defer errors.Debug()

	errors.TT(true, func(m *errors.M) {
		m.Msg("test tt").
			Tag("tag").
			M("k", "v")

	})
}

func TestWrap(t *testing.T) {
	defer errors.Debug()

	errors.Wrap(es.New("test"), "test")
}

func TestWrapM(t *testing.T) {
	defer errors.Debug()

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
	errors.Cfg.Debug=false
	defer errors.Debug()

	// 开始性能分析, 返回一个停止接口
	stopper := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	// 在main()结束时停止性能分析
	defer stopper.Stop()

	for i:=0;i<10000;i++{
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
	defer errors.Debug()

	init11()
}

func TestTry(t *testing.T) {
	defer errors.Debug()

	errors.Panic(errors.Try(errors.T, true, "sss"))
}

func TestTask(t *testing.T) {
	defer errors.Debug()

	errors.Wrap(errors.Try(func() {
		errors.Wrap(es.New("dd"), "err ")
	}), "test wrap")
}

func TestHandle(t *testing.T) {
	defer errors.Debug()

	func() {
		defer errors.Handle(func() {})

		errors.Wrap(es.New("hello error"), "sss")
	}()

}

func TestErrHandle(t *testing.T) {
	defer errors.Debug()

	errors.ErrHandle(errors.Try(errors.T, true, "test 12345"), func(err *errors.Err) {
		err.P()
	})

	errors.ErrHandle(errors.Try(errors.T, true, "test 12345"))

	errors.ErrHandle("ttt")
	errors.ErrHandle(es.New("eee"))
	errors.ErrHandle([]string{"dd"})
}

func TestIsZero(t *testing.T) {
	defer errors.Debug()

	var ss = func() map[string]interface{} {
		return make(map[string]interface{})
	}

	var ss1 = func() map[string]interface{} {
		return nil
	}

	var s = 1
	var ss2 map[string]interface{}
	t.Log(errors.IsZero(1))
	t.Log(errors.IsZero(1.2))
	t.Log(errors.IsZero(nil))
	t.Log(errors.IsZero("ss"))
	t.Log(errors.IsZero(map[string]interface{}{}))
	t.Log(errors.IsZero(ss()))
	t.Log(errors.IsZero(ss1()))
	t.Log(errors.IsZero(&s))
	t.Log(errors.IsZero(ss2))
}

func TestResp(t *testing.T) {
	defer errors.Resp(func(err *errors.Err) {
		err.StackTrace()
	})

	errors.T(true, "data handle")
}
