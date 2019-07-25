package internal

import (
	"os"
)

var IsDebug = GetOkOnce(func() bool {
	debug := os.Getenv("debug")
	return debug == "true" || debug == "t" || debug == "1" || debug == "ok"
})

var IsSkipErrorFile = GetOkOnce(func() bool {
	skipErrorFile := os.Getenv("skip_error_file")
	return skipErrorFile == "true" || skipErrorFile == "t" || skipErrorFile == "1" || skipErrorFile == "ok"
})

func InitDebug() {
	Wrap(os.Setenv("debug", "true"), "set debug env error")
}

func InitSkipErrorFile() {
	Wrap(os.Setenv("skip_error_file", "true"), "set debug env error")
}
