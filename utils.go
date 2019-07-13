package errors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func IsNone(val interface{}) bool {
	return val == nil || IsZero(reflect.ValueOf(val))
}

func IsZero(val reflect.Value) bool {
	defer Throw(func() {})

	if !val.IsValid() {
		return true
	}

	switch val.Kind() {
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

func P(d ...interface{}) {
	defer Throw(func() {})

	for _, i := range d {
		if IsZero(reflect.ValueOf(i)) {
			continue
		}

		dt, err := json.MarshalIndent(i, "", "\t")
		Wrap(err, "P json MarshalIndent error")
		fmt.Println(string(dt))
	}
}

var srcDir = filepath.Join(build.Default.GOPATH, "src") + string(os.PathSeparator)
var modDir = filepath.Join(build.Default.GOPATH, "pkg", "mod") + string(os.PathSeparator)

func funcCaller(callDepth int) string {
	fn, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		log.Error().Msg("no func caller error")
		return "no func caller"
	}

	var buf = _bytesPool.Get().(*strings.Builder)
	defer _bytesPool.Put(buf)
	defer buf.Reset()

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")

	ma := strings.Split(runtime.FuncForPC(fn).Name(), ".")
	buf.WriteString(ma[len(ma)-1])
	return buf.String()
}

func getCallerFromFn(fn reflect.Value) string {
	_fn := fn.Pointer()
	_e := runtime.FuncForPC(_fn)
	file, line := _e.FileLine(_fn)

	var buf = _bytesPool.Get().(*strings.Builder)
	defer _bytesPool.Put(buf)
	defer buf.Reset()

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")

	ma := strings.Split(_e.Name(), ".")
	buf.WriteString(ma[len(ma)-1])
	return strings.TrimPrefix(strings.TrimPrefix(buf.String(), srcDir), modDir)
}

var _bytesPool = &sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}
