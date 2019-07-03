package errors

import (
	"time"
)

var Cfg = struct {
	Debug       bool
	MaxObj      uint8
	MaxRetryDur time.Duration
}{
	Debug:       true,
	MaxObj:      15,
	MaxRetryDur: time.Hour * 24,
}

var ErrTags = struct {
	UnknownTypeCode uint16
}{
	1000,
}

var _errTags = make(map[uint16]string)

func ErrCodeRegistry(code uint16, err string) {
	if _err, ok := _errTags[code]; ok {
		T(ok, "%d has existed, err(%s)", code, _err)
	}
	_errTags[code] = err
}

func GetErrTags() map[uint16]string {
	return _errTags
}

func GetErrTag(code uint16) string {
	if _dt, ok := _errTags[code]; ok {
		return _dt
	}
	return ""
}

func init() {
	ErrCodeRegistry(ErrTags.UnknownTypeCode, "errors_unknown_type")
}
