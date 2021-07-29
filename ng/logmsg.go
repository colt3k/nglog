package ng

import (
	"bytes"
	"fmt"
	"github.com/colt3k/nglogint"
	"github.com/colt3k/nglogint/enum"
	"github.com/colt3k/nglogint/types"
	"os"
	"sync"
	"time"

)


var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}



//Create a pool for reuse in order to speed up process of logging
type LogMsg struct {
	Logger *StdLogger
	//None, Fatal, Error, Warn, Info, Debug
	Level enum.LogLevel
	//Message passed into logger
	Message string
	//Fields passed to template i.e. for JSON output
	Fields []types.Fields
	// Time entry was created
	Time time.Time
	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer

	caller string
}

func NewEntry(logger *StdLogger) *LogMsg {
	t := &LogMsg{
		Logger: logger,
		// Default is three fields, give a little extra room
		Fields: make([]types.Fields, 5),
	}
	return t
}

func (entry *LogMsg) MsgLogger() nglogint.Logger {
	return entry.Logger
}
func (entry *LogMsg) MsgFields() []types.Fields {
	return entry.Fields
}
func (entry *LogMsg) MsgLevel() enum.LogLevel {
	return entry.Level
}
func (entry *LogMsg) MessageStr() string {
	return entry.Message
}
func (entry *LogMsg) MsgTime() time.Time {
	return entry.Time
}
// Add a map of fields to the Entry.
func (entry *LogMsg) WithFields(fields []types.Fields) nglogint.Msg {
	t := make([]types.Fields, 0)
	//data := make(Fields, len(entry.Fields)+len(fields))
	for _, d := range entry.Fields {
		t = append(t, d)
	}
	for _, d := range fields {
		t = append(t, d)
	}

	return &LogMsg{Logger: entry.Logger, Fields: t}
}
func (entry *LogMsg) Info(args ...interface{}) {
	if entry.Logger.level >= enum.INFO {
		entry.log(enum.INFO, fmt.Sprint(args...))
	}
}
func (entry *LogMsg) Error(args ...interface{}) {
	if entry.Logger.level >= enum.ERROR {
		entry.log(enum.ERROR, fmt.Sprint(args...))
	}
}
func (entry *LogMsg) Fatal(args ...interface{}) {
	if entry.Logger.level >= enum.FATAL {
		entry.log(enum.FATAL, fmt.Sprint(args...))
	}
	Exit(1)
}

func (entry *LogMsg) LogEnt(level enum.LogLevel, format, caller string, rtn bool, args ...interface{}) {
	entry.caller = caller
	//has format but doesn't perform a return
	// https://golang.org/pkg/fmt/
	if len(format) > 0 && !rtn {
		entry.log(level, fmt.Sprint(fmt.Sprintf(format, args...)))
	} else if len(format) == 0 && !rtn { // no format and no return
		entry.log(level, fmt.Sprint(args...))
	} else if len(format) == 0 && rtn { // no format and a return
		entry.log(level, fmt.Sprint(entry.sprintlnn(args...)))
	}

	if level == enum.FATAL {
		Exit(1)
	}
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *LogMsg) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
func (entry LogMsg) log(level enum.LogLevel, msg string) {
	var buffer *bytes.Buffer
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	//entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	//@TODO Send to Appenders instead
	entry.write()

	entry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	//if level <= PanicLevel {
	//	panic(&entry)
	//}
}

func (entry *LogMsg) write() {
	entry.Logger.MU.Lock()
	defer entry.Logger.MU.Unlock()
	//Process to appenders HERE!!!
	for _, d := range entry.Logger.appenders {
		serialized, err := entry.Logger.Formatter.Format(entry, d.DisableColor())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to obtain reader, \n%+v", err)
		}
		if d.Applicable(entry.caller) {
			d.Process(serialized)
		}
	}
}

func Exit(code int) {
	//runHandlers()
	os.Exit(code)
}
