package errors

import (
	"sync"
)

var _kerr = &sync.Pool{
	New: func() interface{} {
		return &Err{}
	},
}

func kerrGet() *Err {
	defer Handle(func() {})

	return _kerr.Get().(*Err)
}

func kerrPut(m *Err) {
	defer Handle(func() {})

	_kerr.Put(m)
}
