package errors

import (
	"reflect"
)

var Cfg = struct {
	Debug  bool
	MaxObj uint8
}{
	Debug:  true,
	MaxObj: 15,
}

var errType = reflect.TypeOf(&Err{})

var ErrTag = struct {
	UnknownErr string
}{
	UnknownErr: "unknown",
}
