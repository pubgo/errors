package assert

import (
	"errors"
	"testing"
)

func TestT(t *testing.T) {
	defer Debug()

	T(true, "test t")
}

func TestTT(t *testing.T) {
	defer Debug()

	TT(true, func(m *M) {
		m.Msg("test tt").
			Tag("tag").
			M("k", "v")

	})
}

func TestWrap(t *testing.T) {
	defer Debug()

	Wrap(errors.New("test"), "test")
}

func TestWrapM(t *testing.T) {
	defer Debug()

	WrapM(errors.New("dd"), func(m *M) {
		m.Msg("test")
	})
}

func a1() {
	defer Handle(func(m *M) {
		m.Msg("test SWrap")
	})

	WrapM(errors.New("sbhbhbh"), func(m *M) {
		m.Msg("test shhh").
			M("ss", 1).
			M("input", 2)
	})
}

func TestName(t *testing.T) {
	defer Debug()

	ErrHandle(Try(a1), func(err *Err) {
		err.P()
	})
}

func TestTry(t *testing.T) {
	defer Debug()

	Cfg.Debug = true

	T(true, "sss")
}

func TestTask(t *testing.T) {
	defer Debug()

	Wrap(Try(func() {
		Wrap(errors.New("dd"), "err ")
	}), "test wrap")
}

func test123() {
	defer Handle(func(m *M) {
		m.Msg("test panic %d", 33)
	})

	Wrap(errors.New("hello error"), "sss")
}

func TestExpect11(t *testing.T) {
	defer Debug()

	Cfg.Debug = true

	test123()
}

func TestIsNil(t *testing.T) {
	defer Debug()

	var ss = func() map[string]interface{} {
		return make(map[string]interface{})
	}

	var ss1 = func() map[string]interface{} {
		return nil
	}

	var s = 1
	var ss2 map[string]interface{}
	t.Log(IsNil(1))
	t.Log(IsNil(1.2))
	t.Log(IsNil(nil))
	t.Log(IsNil("ss"))
	t.Log(IsNil(map[string]interface{}{}))
	t.Log(IsNil(ss()))
	t.Log(IsNil(ss1()))
	t.Log(IsNil(&s))
	t.Log(IsNil(ss2))
}

func TestResp(t *testing.T) {
	defer Resp(func(err *Err) {
		err.Tag()
		err.StackTrace()
	})

	T(true, "data handle")
}
