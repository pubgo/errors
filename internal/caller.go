package internal

import (
	"github.com/rs/zerolog/log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// func caller depth, default=2
const callDepth = 2

//var srcDir = filepath.Join(build.Default.GOPATH, "src") + string(os.PathSeparator)
//var modDir = filepath.Join(build.Default.GOPATH, "pkg", "mod") + string(os.PathSeparator)

func FuncCaller(callDepth int) string {
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

func GetCallerFromFn(fn reflect.Value) string {
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
	return buf.String()
}

var _bytesPool = &sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}
