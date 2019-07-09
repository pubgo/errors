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
	MaxObj:      15,             // the max length of m keys is 15
	MaxRetryDur: time.Hour * 24, // the max retry duration is one day
}

var ErrTags = struct {
	UnknownTypeCode string
}{
	"errors_unknown_type",
}

var _errTags = make(map[string]bool)

func ErrCodeRegistry(err string) {
	if _, ok := _errTags[err]; ok {
		T(ok, "%s has existed", err)
	}
	_errTags[err] = true
}

func GetErrTags() map[string]bool {
	return _errTags
}

func init() {
	ErrCodeRegistry(ErrTags.UnknownTypeCode)
}
