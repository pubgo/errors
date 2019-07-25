package internal

import "time"

// default config
var Cfg = struct {
	MaxObj      uint8
	MaxRetryDur time.Duration
}{
	MaxObj:      15,             // the max length of m keys, default=15
	MaxRetryDur: time.Hour * 24, // the max retry duration, default=one day
}
