package ng

import (
	"github.com/colt3k/nglog"
	"github.com/colt3k/nglog/enum"
	"io"
	"strconv"
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
		lgr.ClrERR = ColorFmt(clr)
	}
}
func WarnColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrWARN = ColorFmt(clr)
	}
}
func InfoColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrINFO = ColorFmt(clr)
	}
}
func DebugColor(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrDEBUG = ColorFmt(clr)
	}
}
func DebugLvl2Color(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrDEBUGL2 = ColorFmt(clr)
	}
}
func DebugLvl3Color(clr ColorAttr) LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrDEBUGL3 = ColorFmt(clr)
	}
}
func ColorsOn() LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrDEFAULT = ColorFmt(FgWhite)
		lgr.ClrERR = ColorFmt(FgRed)
		lgr.ClrWARN = ColorFmt(FgYellow)
		lgr.ClrINFO = ColorFmt(FgBlue)
		lgr.ClrDEBUG = std.ClrDEFAULT
		lgr.ClrDEBUGL2 = std.ClrDEFAULT
		lgr.ClrDEBUGL3 = std.ClrDEFAULT
	}
}

func HiColorsOn() LogOption {
	return func(lgr *StdLogger) {
		lgr.ClrDEFAULT = ColorBrightFmt(HiBlack)
		lgr.ClrERR = ColorBrightFmt(HiRed)
		lgr.ClrWARN = ColorBrightFmt(HiYellow)
		lgr.ClrINFO = ColorBrightFmt(HiBlue)
		lgr.ClrDEBUG = std.ClrDEFAULT
		lgr.ClrDEBUGL2 = std.ClrDEFAULT
		lgr.ClrDEBUGL3 = std.ClrDEFAULT
	}
}

func ColorFmt(color ColorAttr) string {
	//fmt.Println("called to get color ", color)
	return "\u001b[" + strconv.Itoa(int(color)) + "m"
}

func ColorBrightFmt(color ColorAttr) string {
	return "\u001b[" + strconv.Itoa(int(color)) + ";1m"
}
func Formatter(f nglog.Layout) LogOption {
	return func(lgr *StdLogger) {
		lgr.Formatter = f
	}
}

func Appenders(a ...Appender) LogOption {
	return func(lgr *StdLogger) {
		lgr.appenders = a
	}
}
