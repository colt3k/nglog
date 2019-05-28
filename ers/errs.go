package ers

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	lvl = 4
)

type Error interface {
	Err(e error, skipTrace bool, msg ...string) bool
	NoPrintErr(e error) bool
	NotErr(e error, msg ...string) bool
	NotErrSkipTrace(e error, skipTrace bool, msg ...string) bool
	StopErr(e error, msg ...string)
	WarnErr(e error, msg ...string) bool
}

func FindCaller(skip int) (string, string, int) {
	return Trace(skip)
}
func caller(skip int) (uintptr, string, int) {

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	log.Println(file)
	//Strip off path except last part and file name
	tmp := strings.Split(file, string(filepath.Separator))
	f2 := tmp[len(tmp)-2:]
	file = strings.Join(f2, string(filepath.Separator))

	return pc, file, line
}

func Trace(skip int) (string, string, int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.File, frame.Function, frame.Line
}
