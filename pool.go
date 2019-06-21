package assert

import (
	"sync"
)

var _kerr = &sync.Pool{
	New: func() interface{} {
		return &Err{}
	},
}

func kerrGet() *Err {
	defer Handle()

	return _kerr.Get().(*Err)
}

func kerrPut(m *Err) {
	defer Handle()

	_kerr.Put(m)
}
