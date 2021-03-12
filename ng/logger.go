package ng

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/colt3k/nglog/internal/pkg/enum"
	"github.com/colt3k/nglog/internal/pkg/util"
)

const DefaultTimestampFormat = time.RFC3339

var (
	once  ResyncOnce
	std   = NewLogger()
	DBGL3 = enum.DBGL3
	DBGL2 = enum.DBGL2
	DEBUG = enum.DEBUG
	FATAL = enum.FATAL
	FATALNE = enum.FATALNOEXIT
	ERROR = enum.ERROR
	WARN  = enum.WARN
	INFO  = enum.INFO
	NONE  = enum.NONE
)

type StdLogger struct {
	Formatter    Layout
	appenders    []Appender
	Out          io.Writer
	level        enum.LogLevel
	depth        int
	flags        enum.Flags
	ColorDEFAULT string
	ColorERR     string
	ColorWARN    string
	ColorINFO    string
	ColorDEBUG   string
	ColorDEBUGL2 string
	ColorDEBUGL3 string
	// Reusable empty entry
	entryPool sync.Pool
	// Used to sync writing to the log. Locking is enabled by Default
	MU  util.MutexWrap
	Now time.Time
}

/*
NewLogger can be used instead of the default to change settings
*/
func NewLogger(opts ...LogOption) *StdLogger {

	var t *StdLogger

	once.Do(func() {
		t = new(StdLogger)
		t.Now = time.Now() // sets when first created
		t.Formatter = &TextLayout{ForceColor: true}
		//Set colors by default, still can override with opts
		t.ColorDEFAULT = ColorFmt(FgWhite)
		t.ColorERR = ColorFmt(FgRed)
		t.ColorWARN = ColorFmt(FgYellow)
		t.ColorINFO = ColorFmt(FgBlue)
		t.ColorDEBUG = t.ColorDEFAULT
		t.ColorDEBUGL2 = t.ColorDEFAULT
		t.ColorDEBUGL3 = t.ColorDEFAULT
		t.depth = 5 //Default how deep to go in order to find caller
		for _, opt := range opts {
			opt(t)
		}
		// Default logger
		if t.level == 0 {
			t.level = enum.INFO
		}

		if t.Out == nil {
			t.Out = os.Stdout
		}

		if t.appenders == nil {
			t.appenders = []Appender{NewConsoleAppender("*")}
		}

		t.SetFlags(t.flags)

	})

	return t
}

// newEntry retrieve a LogMsg from the pool or if none are available create one
func (l *StdLogger) newEntry(fields bool) *LogMsg {
	lmsg, ok := l.entryPool.Get().(*LogMsg)
	if ok {
		if !fields {
			lmsg.Fields = make([]Fields, 0)
		}
		return lmsg
	}

	return NewEntry(l)
}

func (l *StdLogger) releaseEntry(entry *LogMsg) {
	l.entryPool.Put(entry)
}

func Modify(opts ...LogOption) {
	for _, opt := range opts {
		opt(std)
	}
}
func ShowAppenders() []string {
	appendrs := make([]string,0)
	for _,j := range std.appenders {
		appendrs = append(appendrs, j.Name())
	}
	return appendrs
}
func Logger() *StdLogger {
	return std
}
func DebugL3ln(args ...interface{}) {
	Logln(enum.DBGL3, args...)
}
func DebugL2(format string, args ...interface{}) {
	Logf(enum.DBGL3, format, args...)
}
func DebugX2ln(args ...interface{}) {
	Logln(enum.DBGL2, args...)
}
func DebugX2(format string, args ...interface{}) {
	Logf(enum.DBGL2, format, args...)
}
func Debugln(args ...interface{}) {
	Logln(enum.DEBUG, args...)
}
func Debug(format string, args ...interface{}) {
	Logf(enum.DEBUG, format, args...)
}
func Infoln(args ...interface{}) {
	Logln(enum.INFO, args...)
}
func Info(format string, args ...interface{}) {
	Logf(enum.INFO, format, args...)
}
func Error(format string, args ...interface{}) {
	Logf(enum.ERROR, format, args...)
}
func Errorln(args ...interface{}) {
	Logln(enum.ERROR, args...)
}
func Logln(lvl enum.LogLevel, args ...interface{}) {
	if lvl <= std.level {
		std.Logln(lvl, args...)
	}
}
func Logf(lvl enum.LogLevel, format string, args ...interface{}) {
	if lvl <= std.level {
		std.Logf(lvl, format, args...)
	}
}
func Flags() int {
	return std.Flags()
}

func Level() enum.LogLevel {
	return std.Level()
}
func IsDebugL3() bool {
	return std.Level() == enum.DBGL3
}
func IsDebugX2() bool {
	return std.Level() == enum.DBGL2
}
func IsDebug() bool {
	return std.Level() == enum.DEBUG
}
func IsWarn() bool {
	return std.Level() == enum.WARN
}
func IsInfo() bool {
	return std.Level() == enum.INFO
}
func IsError() bool {
	return std.Level() == enum.ERROR
}
func IsFatal() bool {
	return std.Level() == enum.FATAL
}
func IsFatalNE() bool {
	return std.Level() == enum.FATALNOEXIT
}
func IsNone() bool {
	return std.Level() == enum.NONE
}
func Print(args ...interface{}) {
	std.Print(args...)
}
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}
func PrintTypeOfValue(arg interface{}) {
	std.PrintTypeOfValue(arg)
}
func PrintStructWithFieldNames(arg interface{}) {
	std.PrintStructWithFieldNames(arg)
}
func PrintGoSyntaxOfValue(arg interface{}) {
	std.PrintGoSyntaxOfValue(arg)
}
func Println(args ...interface{}) {
	std.Println(args...)
}
func SetLevel(level enum.LogLevel) {
	std.SetLevel(level)
}
func SetFlags(flg enum.Flags) {
	std.SetFlags(flg)
}
func DisableTimestamp() {
	std.Formatter.DisableTimeStamp()
}
func EnableTimestamp() {
	std.Formatter.DisableTimeStamp()
}
func SetFormatter(formatter Layout) {
	std.SetFormatter(formatter)
}
func ShowConfig() {
	std.ShowOptions()
}
func (l *StdLogger) ShowOptions() {
	var buf bytes.Buffer
	buf.WriteString(l.Formatter.Description())

	buf.WriteString(l.ColorDEFAULT+" DefaultColor:")
	buf.WriteString(l.ColorERR+" ErrorColor:")
	buf.WriteString(l.ColorWARN+" WarnColor:")
	buf.WriteString(l.ColorINFO+" InfoColor:")
	buf.WriteString(l.ColorDEBUG+" DebugColor:")
	buf.WriteString(l.ColorDEBUGL2+" DebugL2Color:")
	buf.WriteString(l.ColorDEBUGL3+" DebugL3Color:")
	entry := l.newEntry(false)
	entry.LogEnt(enum.DEBUG, "", l.Caller(), false, buf.String())
	l.releaseEntry(entry)
}
func (l *StdLogger) Flags() int {
	return log.Flags()
}
func (l *StdLogger) Level() enum.LogLevel {
	return l.level
}
func (l *StdLogger) Print(args ...interface{}) {
	entry := l.newEntry(false)
	entry.LogEnt(enum.NONE, "", l.Caller(), false, args...)
	l.releaseEntry(entry)
}
func (l *StdLogger) Printf(format string, args ...interface{}) {
	entry := l.newEntry(false)
	entry.LogEnt(enum.NONE, format, l.Caller(), false, args...)
	l.releaseEntry(entry)
}
func (l *StdLogger) PrintTypeOfValue(arg interface{}) {
	entry := l.newEntry(false)
	entry.LogEnt(enum.NONE, "%T", l.Caller(), false, arg)
	l.releaseEntry(entry)
}
func (l *StdLogger) PrintStructWithFieldNames(arg interface{}) {
	entry := l.newEntry(false)
	entry.LogEnt(enum.NONE, "%+v", l.Caller(), false, arg)
	l.releaseEntry(entry)
}
func (l *StdLogger) PrintGoSyntaxOfValue(arg interface{}) {
	entry := l.newEntry(false)
	entry.LogEnt(enum.NONE, "%#v", l.Caller(), false, arg)
	l.releaseEntry(entry)
}
func (l *StdLogger) Println(args ...interface{}) {
	entry := l.newEntry(false)
	entry.LogEnt(enum.NONE, "", l.Caller(), true, args...)
	l.releaseEntry(entry)
}
func (l *StdLogger) SetLevel(level enum.LogLevel) {
	l.level = level
}

func (l *StdLogger) Caller() string {
	//_, function, _ := util.FindCaller(l.depth)
	//log.Println("depth:", l.depth, "Function:", function)
	_, function, _ := util.FindIssue(l.depth)
	return function
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields []Fields) *LogMsg {
	return std.WithFields(fields)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (l *StdLogger) WithFields(fields []Fields) *LogMsg {
	entry := l.newEntry(true)
	defer l.releaseEntry(entry)
	return entry.WithFields(fields)
}

/**
SetFlags
Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
Ltime                         // the time in the local time zone: 01:23:23
Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
Llongfile                     // full file name and line number: /a/b/c/d.go:23
Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
LstdFlags     = Ldate | Ltime // initial values for the standard logger
*/
func (l *StdLogger) SetFlags(flg enum.Flags) {
	log.SetFlags(int(flg))
}
func (l *StdLogger) SetFormatter(formatter Layout) {
	l.Formatter = formatter
}

func (l *StdLogger) Logln(lvl enum.LogLevel, args ...interface{}) {

	switch lvl {
	case enum.NONE:
		if l.level >= enum.NONE {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.FATAL:
		if l.level >= enum.FATAL {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
			os.Exit(1)
		}
	case enum.FATALNOEXIT:
		if l.level >= enum.FATALNOEXIT {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.ERROR:
		if l.level >= enum.ERROR {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.WARN:
		if l.level >= enum.WARN {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.INFO:
		if l.level >= enum.INFO {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.DEBUG:
		if l.level >= enum.DEBUG {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.DBGL2:
		if l.level >= enum.DBGL2 {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	case enum.DBGL3:
		if l.level >= enum.DBGL3 {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, "", l.Caller(), true, args...)
			l.releaseEntry(entry)
		}
	default:
		entry := l.newEntry(false)
		entry.LogEnt(lvl, "", l.Caller(), true, args...)
		l.releaseEntry(entry)
	}
}

func (l *StdLogger) Logf(lvl enum.LogLevel, format string, args ...interface{}) {

	switch lvl {
	case enum.NONE:
		if l.level >= enum.NONE {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.FATAL:
		if l.level >= enum.FATAL {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.FATALNOEXIT:
		if l.level >= enum.FATALNOEXIT {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.ERROR:
		if l.level >= enum.ERROR {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.WARN:
		if l.level >= enum.WARN {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.INFO:
		if l.level >= enum.INFO {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.DEBUG:
		if l.level >= enum.DEBUG {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.DBGL2:
		if l.level >= enum.DBGL2 {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	case enum.DBGL3:
		if l.level >= enum.DBGL3 {
			entry := l.newEntry(false)
			entry.LogEnt(lvl, format, l.Caller(), false, args...)
			l.releaseEntry(entry)
		}
	default:
		entry := l.newEntry(false)
		entry.LogEnt(lvl, format, l.Caller(), false, args...)
		l.releaseEntry(entry)
	}
}
