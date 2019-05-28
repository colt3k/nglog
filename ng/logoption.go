package ng

import (
	"io"
	"strconv"

	"github.com/colt3k/nglog/internal/pkg/enum"
)

type LogOption func(h *StdLogger)

// Logger returns a MDHashOption that sets the logger for the hash.
func LogLevel(l enum.LogLevel) LogOption {
	return func(lgr *StdLogger) {
		lgr.level = l
	}
}

func LogOut(i io.Writer) LogOption {
	return func(lgr *StdLogger) {
		lgr.Out = i
	}
}
func CallDepth(i int) LogOption {
	return func(lgr *StdLogger) {
		lgr.depth = i
	}
}
func SetFlgs(i enum.Flags) LogOption {
	return func(lgr *StdLogger) {
		lgr.flags = i
	}
}
func ErrColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorERR = ColorFmt(clr)
	}
}
func WarnColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorWARN = ColorFmt(clr)
	}
}
func InfoColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorINFO = ColorFmt(clr)
	}
}
func DebugColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEBUG = ColorFmt(clr)
	}
}
func ColorsOn() LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEFAULT = ColorFmt(FgWhite)
		lgr.ColorERR = ColorFmt(FgRed)
		lgr.ColorWARN = ColorFmt(FgYellow)
		lgr.ColorINFO = ColorFmt(FgBlue)
		lgr.ColorDEBUG = std.ColorDEFAULT
	}
}

func HiColorsOn() LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEFAULT = ColorBrightFmt(HiBLACK)
		lgr.ColorERR = ColorBrightFmt(HiRED)
		lgr.ColorWARN = ColorBrightFmt(HiYELLOW)
		lgr.ColorINFO = ColorBrightFmt(HiBLUE)
		lgr.ColorDEBUG = std.ColorDEFAULT
	}
}

func ColorFmt(color ColorAttr) string {
	//fmt.Println("called to get color ", color)
	return "\u001b[" + strconv.Itoa(int(color)) + "m"
}

func ColorBrightFmt(color ColorAttr) string {
	return "\u001b[" + strconv.Itoa(int(color)) + ";1m"
}
func Formatr(f Layout) LogOption {
	return func(lgr *StdLogger) {
		lgr.Formatter = f
	}
}

func Appenders(a ...Appender) LogOption {
	return func(lgr *StdLogger) {
		lgr.appenders = a
	}
}
