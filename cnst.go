package errors

import (
	"reflect"
	"time"
)

var Cfg = struct {
	Debug        bool
	MaxObj       uint8
	MaxRetryDur  time.Duration
	MaxRetryTime uint64
}{
	Debug:        true,
	MaxObj:       15,
	MaxRetryDur:  time.Hour * 24,
	MaxRetryTime: 0,
}

var errType = reflect.TypeOf(&Err{})

var ErrTag = struct {
	UnknownErr string
}{
	UnknownErr: "unknown",
}
