package errors

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const callDepth = 2

var srcDir = filepath.Join(build.Default.GOPATH, "src") + string(os.PathSeparator)
var modDir = filepath.Join(build.Default.GOPATH, "pkg", "mod") + string(os.PathSeparator)

func funcCaller(callDepth int) string {
	fn, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return "no func caller"
	}

	ma := strings.Split(runtime.FuncForPC(fn).Name(), ".")

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(".")
	buf.WriteString(ma[len(ma)-1])
	return strings.TrimPrefix(strings.TrimPrefix(buf.String(), srcDir), modDir)
}

func init() {
	log.Logger = log.Output(zerolog.NewConsoleWriter()).With().Caller().Logger()
}
