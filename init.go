package errors

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

const callDepth = 2

func assertFn(fn interface{}) {
	T(IsNil(fn), "the func is nil")

	_v := reflect.TypeOf(fn)
	T(_v.Kind() != reflect.Func, "func type error(%s)", _v.String())
}

var goPath = build.Default.GOPATH
var srcDir = fmt.Sprintf("%s%s", filepath.Join(goPath, "src"), string(os.PathSeparator))
var modDir = fmt.Sprintf("%s%s", filepath.Join(goPath, "pkg", "mod"), string(os.PathSeparator))

func funcCaller(callDepth int) string {
	fn, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return "no func caller"
	}

	ma := strings.Split(runtime.FuncForPC(fn).Name(), ".")
	return strings.TrimPrefix(strings.TrimPrefix(fmt.Sprintf("%s:%d:%s", file, line, ma[len(ma)-1]), srcDir), modDir)
}

func IsNil(p interface{}) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = false
		}
	}()

	if p == nil {
		return true
	}

	if !reflect.ValueOf(p).IsValid() {
		return true
	}

	return reflect.ValueOf(p).IsNil()
}
