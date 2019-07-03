package errors

import (
	"reflect"
	"time"
)

func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

func Default(t, f interface{}) interface{} {
	if IsZero(reflect.ValueOf(t)) {
		return f
	}
	return t
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

func FnOf(fn interface{}) (reflect.Value, bool, reflect.Value) {
	defer Handle()()

	_fn := reflect.ValueOf(fn)

	assertFn(_fn)

	var variadicType reflect.Value
	var isVariadic = _fn.Type().IsVariadic()
	if isVariadic {
		variadicType = reflect.New(_fn.Type().In(_fn.Type().NumIn() - 1).Elem()).Elem()
	}

	return _fn, isVariadic, variadicType
}
