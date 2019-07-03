package errors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"sync"
)

type _Err struct {
	Tag    uint16                 `json:"tag,omitempty"`
	M      map[string]interface{} `json:"m,omitempty"`
	Err    error                  `json:"err,omitempty"`
	Msg    string                 `json:"msg,omitempty"`
	Caller string                 `json:"caller,omitempty"`
	Sub    *_Err                  `json:"sub,omitempty"`
}

func (t *_Err) String() string {
	defer Handle()()

	_dt, err := json.Marshal(t)
	Wrap(err, "json marshal error")

	return string(_dt)
}

var _errPool = sync.Pool{
	New: func() interface{} {
		return new(Err)
	},
}

func errGet() *Err {
	return _errPool.Get().(*Err)
}

type Err struct {
	tag    uint16
	m      map[string]interface{}
	err    error
	msg    string
	caller string
	sub    *Err
}

func (t *Err) put() {
	t.reset()
	_errPool.Put(t)
}

func (t *Err) reset() {
	t.tag = 0
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
	err := t._err()
	c := err
	for t.sub != nil {
		c.Sub = t.sub._err()
		t.sub = t.sub.sub
		c = c.Sub
	}
	return err
}

func (t *Err) tErr() (err error) {
	err, t.err = t.err, nil
	return
}

func (t *Err) tTag() (tag uint16) {
	tag, t.tag = t.tag, 0
	return
}

func (t *Err) P() {
	P(t.StackTrace())
}

func (t *Err) isNil() bool {
	return t == nil || t.err == nil
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

		if _t.tag != 0 {
			_l = _l.Uint16("err_tag", _t.tag)
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

func (t *Err) SetTag(k string, v uint16) *Err {
	if t.isNil() {
		return t
	}

	t.tag = v
	return t
}

func (t *Err) Tag() uint16 {
	return t.tag
}
