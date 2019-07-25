package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

func P(s string, d ...interface{}) {
	fmt.Print(s)
	for _, i := range d {
		if i == nil || IsNone(i) {
			continue
		}

		dt, err := json.MarshalIndent(i, "", "\t")
		Wrap(err, "P json MarshalIndent error")
		fmt.Println(string(dt))
	}
}

func GetOkOnce(fn func() bool) func() bool {
	var _isOk = false
	var _okOnce sync.Once
	return func() bool {
		_okOnce.Do(func() {
			_isOk = fn()
		})
		return _isOk
	}
}

func LoadEnvFile(envPath string) {
	_p, err := filepath.EvalSymlinks(envPath)
	dt, err := ioutil.ReadFile(_p)
	Wrap(err, "file open error")
	for _, env := range strings.Split(string(dt), "\n") {
		envA := strings.Split(env, "=")
		if len(envA) == 2 {
			Panic(os.Setenv(envA[0], envA[1]))
		}
	}
}
