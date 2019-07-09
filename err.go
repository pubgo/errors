package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

type _Err struct {
	Tag    string                 `json:"tag,omitempty"`
	M      map[string]interface{} `json:"m,omitempty"`
	Err    error                  `json:"err,omitempty"`
	Msg    string                 `json:"msg,omitempty"`
	Caller []string               `json:"caller,omitempty"`
	Sub    *_Err                  `json:"sub,omitempty"`
}

func (t *_Err) String() string {

	buf := &strings.Builder{}
	defer buf.Reset()

	_err := json.NewEncoder(buf).Encode(t)
	Wrap(_err, "json marshal error")
	return buf.String()
}

var _valuePool = sync.Pool{
	New: func() interface{} {
		return []reflect.Value{}
	},
}

func valueGet() []reflect.Value {
	return _valuePool.Get().([]reflect.Value)
}

func valuePut(v []reflect.Value) {
	v = v[:0]
	_valuePool.Put(v)
}

type Err struct {
	tag    string
	m      map[string]interface{}
	err    error
	msg    string
	caller []string
	sub    *Err
}

func (t *Err) reset() {
	t.tag = ""
	t.m = nil
	t.err = nil
	t.msg = ""
	t.caller = t.caller[:0]
	t.sub = nil
}

func (t *Err) _err() *_Err {
	return &_Err{
		Tag:    t.tag,
		M:      t.m,
		Err:    t.err,
		Msg:    t.msg,
		Caller: t.caller,
	}
}

func (t *Err) Err() error {
	return t.err
}

func (t *Err) Error() string {
	return t.err.Error()
}

func (t *Err) StackTrace() *_Err {
	if t.isNil() {
		return nil
	}

	kerr := t._err()
	c := kerr
	for t.sub != nil {
		c.Sub = t.sub._err()
		t.sub = t.sub.sub
		c = c.Sub
	}
	return kerr
}

func (t *Err) tErr() (err error) {
	err, t.err = t.err, nil
	return
}

func (t *Err) tTag() (tag string) {
	tag, t.tag = t.tag, ""
	return
}

func (t *Err) P() string {
	if t.isNil() {
		return ""
	}

	var buf = &strings.Builder{}
	defer buf.Reset()

	var _errs []*_Err
	_t := t
	for _t != nil {
		_errs = append(_errs, _t._err())
		_t = _t.sub
	}

	for i := len(_errs) - 1; i > -1; i-- {
		if len(_errs[i].Caller) > 0 {
			var _err string
			if _errs[i].Err != nil {
				_err = _errs[i].Err.Error()
			}

			var _m string
			if !IsZero(reflect.ValueOf(_errs[i].M)) {
				dt, err := json.MarshalIndent(_errs[i].M, "", "\t")
				Wrap(err, "P json MarshalIndent error")
				_m = string(dt)
			}

			buf.WriteString(fmt.Sprintf("[Debug] %s %s\n  msg: %s\n  %serr%s: %s\n  tag: %s\n  m: %s \n",
				time.Now().Format("2006/01/02 - 15:04:05"), _errs[i].Caller[0], _errs[i].Msg, red, reset, _err, _errs[i].Tag, _m))

			for _, k := range _errs[i].Caller[1:] {
				if strings.Contains(k, "handle.go") {
					continue
				}

				if strings.Contains(k, "testing/testing.go") {
					continue
				}

				buf.WriteString(time.Now().Format("[Debug] 2006/01/02 - 15:04:05 "))
				buf.WriteString(fmt.Sprintln(k))
				buf.WriteString("========================================================================================================================\n\n")
			}
		}
	}

	return buf.String()
}

func (t *Err) isNil() bool {
	return t == nil || t.err == nil || IsZero(reflect.ValueOf(t))
}

func (t *Err) Done() {
	if t.isNil() {
		return
	}

	panic(t)
}

func (t *Err) _msg(msg string, args ...interface{}) *Err {
	if t.isNil() {
		return t
	}

	t.msg = fmt.Sprintf(msg, args...)
	return t
}

func (t *Err) Caller(depth int) *Err {
	if t.isNil() {
		return t
	}

	if t.err != nil && !IsZero(reflect.ValueOf(t.err)) {
		t.caller = append(t.caller, funcCaller(depth))
	}
	return t
}

func (t *Err) M(k string, v interface{}) *Err {
	if t.isNil() {
		return t
	}

	if t.m == nil {
		t.m = make(map[string]interface{}, Cfg.MaxObj)
	}

	t.m[k] = v
	return t
}

func (t *Err) SetTag(tag string) *Err {
	if t.isNil() {
		return t
	}

	t.tag = tag
	return t
}

func (t *Err) Tag() string {
	return t.tag
}
