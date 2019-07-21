package internal

import (
	"encoding/json"
	"fmt"
	"strings"
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

func (t *_Err) reset() {
	t.Tag = ""
	t.M = nil
	t.Err = nil
	t.Msg = ""
	t.Caller = nil
	t.Sub = nil
}

func (t *_Err) String() string {
	defer t.reset()

	buf := &strings.Builder{}
	defer buf.Reset()

	Wrap(json.NewEncoder(buf).Encode(t), "_Err json marshal error")
	return buf.String()
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
	t.caller = nil
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
	defer t.reset()

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

	_filter := func(k string) bool {
		for _, _k := range []string{"handle.go", "testing/testing.go", "src/runtime", "testing/testing.go", "go/src/reflect"} {
			if strings.Contains(k, _k) {
				return true
			}
		}
		return false
	}

	buf.WriteString("========================================================================================================================\n")
	for i := len(_errs) - 1; i > -1; i-- {
		if len(_errs[i].Caller) < 1 {
			continue
		}

		buf.WriteString(fmt.Sprintf(
			"[%s]: %s %s\n"+
				"[ %s ]: %s\n"+
				"[ %s ]: %#v\n"+
				"[ %s ]: %s\n"+
				"[  %s  ]: %#v\n",
			Yellow("Debug"), time.Now().Format("2006/01/02 - 15:04:05"), _errs[i].Caller[0],
			Green("Msg"), _errs[i].Msg,
			Red("Err"), _errs[i].Err,
			Blue("Tag"), _errs[i].Tag,
			Magenta("M"), _errs[i].M))

		for _, k := range _errs[i].Caller[1:] {
			if _filter(k) {
				continue
			}

			buf.WriteString(time.Now().Format("[Debug] 2006/01/02 - 15:04:05 "))
			buf.WriteString(fmt.Sprintln(k))
		}
	}
	buf.WriteString("========================================================================================================================")
	return buf.String()
}

func (t *Err) isNil() bool {
	return t == nil || t.err == nil || IsNone(t)
}

func (t *Err) Done() {
	if !t.isNil() {
		panic(t)
	}
}

func (t *Err) _msg(msg string, args ...interface{}) *Err {
	if !t.isNil() {
		t.msg = fmt.Sprintf(msg, args...)
	}

	return t
}

func (t *Err) Caller(caller string) *Err {
	if !t.isNil() {
		t.caller = append(t.caller, caller)
	}

	return t
}

func (t *Err) M(k string, v interface{}) *Err {
	if !t.isNil() {
		if t.m == nil {
			t.m = make(map[string]interface{}, Cfg.MaxObj)
		}

		t.m[k] = v
	}

	return t
}

func (t *Err) SetTag(tag string) *Err {
	if !t.isNil() {
		t.tag = tag
	}

	return t
}

func (t *Err) Tag() string {
	return t.tag
}
