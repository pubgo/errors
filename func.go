package errors

import (
	"reflect"
	"time"
)

func _If(b bool, t, f interface{}) interface{} {

	if b {
		return t
	}

	return f
}

func fnOf(fn interface{}, args ...interface{}) func() []reflect.Value {
	assertFn(fn)

	t := reflect.ValueOf(fn)
	return func() []reflect.Value {
		var vs []reflect.Value
		for i, p := range args {
			var _v reflect.Value
			if IsZero(p) {
				if t.Type().IsVariadic() {
					i = 0
				}
				_v = reflect.New(t.Type().In(i)).Elem()
			} else {
				_v = reflect.ValueOf(p)
			}

			vs = append(vs, _v)
		}
		return t.Call(vs)
	}
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
