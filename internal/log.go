package internal

import (
	"os"
	"sync"
)

var _debug = false
var _debugOnce sync.Once

func IsDebug() bool {
	_debugOnce.Do(func() {
		debug := os.Getenv("debug")
		_debug = debug == "true" || debug == "t" || debug == "1" || debug == "ok"
	})
	return _debug
}

func InitDebug() {
	Wrap(os.Setenv("debug", "true"), "env set error")
}
