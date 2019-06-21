package errors

import (
	"encoding/json"
	"fmt"
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
	defer Handle()

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

func (t *Err) copy() *Err {
	err := *t
	return &err
}

func (t *Err) Err() error {
	return t.err
}

func (t *Err) Error() string {
	return t.err.Error()
}

func (t *Err) Tag() string {
	return t.tag
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
	err = t.err
	t.err = nil
	return
}

func (t *Err) tTag(tag string) string {
	tag = If(tag == "", t.tag, tag).(string)
	t.tag = ""
	return tag
}

func (t *Err) P() {
	P(t.StackTrace())
}

func newM() *M {
	return &M{m: make(map[string]interface{})}
}

type M struct {
	msg    string
	tag    string
	caller string
	m      map[string]interface{}
}

func (t *M) M(k string, v interface{}) *M {
	t.m[k] = v
	return t
}

func (t *M) Msg(format string, args ...interface{}) *M {
	t.msg = fmt.Sprintf(format, args...)
	return t
}

func (t *M) Tag(tag string) *M {
	t.tag = tag
	return t
}

func (t *M) Caller(depth int) *M {
	t.caller = funcCaller(depth)
	return t
}
