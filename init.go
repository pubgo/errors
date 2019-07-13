package errors

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"reflect"
	"sync"
)

var _valuePool = sync.Pool{
	New: func() interface{} {
		return []reflect.Value{}
	},
}

func valueGet() []reflect.Value {
	return _valuePool.Get().([]reflect.Value)
}

func valuePut(v []reflect.Value) {
	v = v[:0]
	_valuePool.Put(v)
}

func init() {
	log.Logger = log.Output(zerolog.NewConsoleWriter()).With().Caller().Logger()
}
