package errors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
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
	defer Handle(func() {})

	_dt, err := json.Marshal(t)
	Wrap(err, "json marshal error")

	return string(_dt)
}

type Err struct {
	tag    string
	m      map[string]interface{}
	err    error
	msg    string
	caller string
	sub    *Err
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

func (t *Err) tTag() (tag string) {
	tag, t.tag = t.tag, ""
	return
}

func (t *Err) P() {
	P(t.StackTrace())
}

func (t *Err) Log() {
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
	if IsZero(t.err) {
		return
	}
	panic(t)
}

func (t *Err) _msg(msg string, args ...interface{}) *Err {
	if !IsZero(t.err) {
		t.msg = fmt.Sprintf(msg, args...)
	}
	return t
}

func (t *Err) Caller(depth int) *Err {
	if !IsZero(t.err) {
		t.caller = funcCaller(depth)
	}
	return t
}

func (t *Err) M(k string, v interface{}) *Err {
	if t.m == nil {
		t.m = make(map[string]interface{})
	}

	t.m[k] = v
	return t
}

func (t *Err) Tag(tag string) *Err {
	if !IsZero(t.err) {
		t.tag = tag
	}
	return t
}
