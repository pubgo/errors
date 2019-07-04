package errors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
	"sync"
)

type _Err struct {
	Tag    string                 `json:"tag,omitempty"`
	M      map[string]interface{} `json:"m,omitempty"`
	Err    error                  `json:"err,omitempty"`
	Msg    string                 `json:"msg,omitempty"`
	Caller string                 `json:"caller,omitempty"`
	Sub    *_Err                  `json:"sub,omitempty"`
}

func (t *_Err) String() string {
	defer Handle()()

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
	caller string
	sub    *Err
}

func (t *Err) reset() {
	t.tag = ""
	t.m = nil
	t.err = nil
	t.msg = ""
	t.caller = ""
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

	_t := t
	err := t._err()
	_err := err
	for _t.sub != nil {
		_err.Sub = _t.sub._err()
		_t = _t.sub
	}
	return err
}

func (t *Err) tErr() (err error) {
	err, t.err = t.err, nil
	return
}

func (t *Err) tTag() (tag string) {
	tag, t.tag = t.tag, ""
	return
}

func (t *Err) P() {
	if t.isNil() {
		return
	}

	_t := t
	for _t != nil {
		fmt.Print(_t.caller)
		P(_t._err())
		_t = _t.sub
	}
}

func (t *Err) isNil() bool {
	return t == nil || t.err == nil || IsZero(reflect.ValueOf(t))
}

func (t *Err) Log() {
	if t.isNil() {
		return
	}

	_t := t
	for _t != nil {
		_l := log.Error()

		if _t.err != nil {
			_l = _l.Err(_t.err)
		}

		if _t.tag != "" {
			_l = _l.Str("err_tag", _t.tag)
		}

		if _t.caller != "" {
			_l = _l.Str("err_caller", _t.caller)
		}

		if _t.m != nil {
			_l = _l.Interface("err_m", _t.m)
		}

		_l.Msg(_t.msg)

		_t = _t.sub
	}
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
		t.caller = funcCaller(depth)
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
