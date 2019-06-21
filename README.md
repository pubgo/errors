# errors
error handle for go

```go
func TestT(t *testing.T) {
	defer Debug()

	T(true, "test t")
}
```

```go
func TestTT(t *testing.T) {
	defer Debug()

	TT(true, func(m *M) {
		m.Msg("test tt").
			Tag("tag").
			M("k", "v")

	})
}
```

```go
func TestWrap(t *testing.T) {
	defer Debug()

	Wrap(errors.New("test"), "test")
}
```

```go
func TestWrapM(t *testing.T) {
	defer Debug()

	WrapM(errors.New("dd"), func(m *M) {
		m.Msg("test")
	})
}
```

```go
func testFunc() {
	defer Handle(func() {})

	WrapM(errors.New("sbhbhbh"), func(m *M) {
		m.Msg("test shhh").
			M("ss", 1).
			M("input", 2)
	})
}
```

```go
func TestPanic(t *testing.T) {
	Cfg.Debug = true
	defer Debug()

	testFunc()
	t.Log("ok")
}
```

```go
func init11() {
	defer Handle(func() {})

	//T(true, "sss")
	TT(true, func(m *M) {
		m.Msg("test tt")
	})
}
```

```go
func TestT2(t *testing.T) {
	Cfg.Debug = true
	defer Debug()

	init11()
}
```

```go
func TestTry(t *testing.T) {
	defer Debug()

	Cfg.Debug = true

	Panic(Try(FnOf(T, true, "sss")))
}
```

```go
func TestTask(t *testing.T) {
	defer Debug()

	Wrap(Try(func() {
		Wrap(errors.New("dd"), "err ")
	}), "test wrap")
}
```

```go
func TestHandle(t *testing.T) {
	defer Debug()

	defer Handle(func() {})

	Wrap(errors.New("hello error"), "sss")
}
```

```go
func TestErrHandle(t *testing.T) {
	defer Debug()

	ErrHandle(Try(func() {
		T(true, "test 12345")
	}), func(err *Err) {
		err.P()
	})

	ErrHandle(Try(func() {
		T(true, "test 12345")
	}))

	ErrHandle("ttt")
	ErrHandle(errors.New("eee"))
	ErrHandle([]string{"dd"})
}
```

```go
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
```

```go
func TestResp(t *testing.T) {
	defer Resp(func(err *Err) {
		err.Tag()
		err.StackTrace()
	})

	T(true, "data handle")
}
```
