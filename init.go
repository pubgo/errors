package errors

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const callDepth = 2

func assertFn(fn interface{}) {
	T(IsZero(fn), "the func is nil")

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

func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return true
	}

	kind := val.Kind()
	switch kind {
	case reflect.String:
		return val.Len() == 0
	case reflect.Bool:
		return val.Bool() == false
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.Interface, reflect.Slice, reflect.Map:
		return val.IsNil()
	case reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if !IsZero(val.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		if t, ok := val.Interface().(time.Time); ok {
			return t.IsZero()
		} else {
			valid := val.FieldByName("Valid")
			if valid.IsValid() {
				va, ok := valid.Interface().(bool)
				return ok && !va
			}

			return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
		}
	default:
		return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
	}
}
