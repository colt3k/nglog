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
func DebugLvl2Color(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEBUGL2 = ColorFmt(clr)
	}
}
func DebugLvl3Color(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEBUGL3 = ColorFmt(clr)
	}
}
func ColorsOn() LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEFAULT = ColorFmt(FgWhite)
		lgr.ColorERR = ColorFmt(FgRed)
		lgr.ColorWARN = ColorFmt(FgYellow)
		lgr.ColorINFO = ColorFmt(FgBlue)
		lgr.ColorDEBUG = std.ColorDEFAULT
		lgr.ColorDEBUGL2 = std.ColorDEFAULT
		lgr.ColorDEBUGL3 = std.ColorDEFAULT
	}
}
func ColorsOff() LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEFAULT = ""
		lgr.ColorERR = ""
		lgr.ColorWARN = ""
		lgr.ColorINFO = ""
		lgr.ColorDEBUG = ""
		lgr.ColorDEBUGL2 = ""
		lgr.ColorDEBUGL3 = ""
	}
}

func HiColorsOn() LogOption {
	return func(lgr *StdLogger) {
		lgr.ColorDEFAULT = ColorBrightFmt(HiBlack)
		lgr.ColorERR = ColorBrightFmt(HiRed)
		lgr.ColorWARN = ColorBrightFmt(HiYellow)
		lgr.ColorINFO = ColorBrightFmt(HiBlue)
		lgr.ColorDEBUG = std.ColorDEFAULT
		lgr.ColorDEBUGL2 = std.ColorDEFAULT
		lgr.ColorDEBUGL3 = std.ColorDEFAULT
	}
}

func ColorFmt(color ColorAttr) string {
	//fmt.Println("called to get color ", color)
	return "\u001b[" + strconv.Itoa(int(color)) + "m"
}

func ColorBrightFmt(color ColorAttr) string {
	return "\u001b[" + strconv.Itoa(int(color)) + ";1m"
}
func Formatter(f Layout) LogOption {
	return func(lgr *StdLogger) {
		lgr.Formatter = f
	}
}

func Appenders(a ...Appender) LogOption {
	return func(lgr *StdLogger) {
		lgr.appenders = a
	}
}
