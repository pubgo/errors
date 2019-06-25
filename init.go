package errors

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const callDepth = 2

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

func init() {
	log.Logger = log.Output(zerolog.NewConsoleWriter()).With().Caller().Logger()
}
