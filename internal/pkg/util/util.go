package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/colt3k/nglog/ers"
)

var ignore = []string{"nglog/ng/logger.go"}

func AppendKeyValue(b *bytes.Buffer, key string, value interface{}, quoteEmptyField, disableQuoting bool) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	AppendValue(b, value, quoteEmptyField, disableQuoting)
}

func AppendValue(b *bytes.Buffer, value interface{}, quoteEmptyField, disableQuoting bool) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !needsQuoting(stringVal, quoteEmptyField) || disableQuoting {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
func needsQuoting(text string, quoteEmptyField bool) bool {
	if quoteEmptyField && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}
func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// changed to 3rd party
func CheckIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func FindIssue(startLevel int) (string, string, int) {
	level := startLevel
	found := false
	fname, method, ln := ers.FindCaller(startLevel)

	for _, d := range ignore {
		if strings.HasSuffix(fname, d) {
			found = true
		}
	}
	for {
		if !found {
			break
		}
		//log.Logf(log.DEBUG, "Level: %d File: %s", level, fname)
		found = false
		for _, d := range ignore {
			if strings.HasSuffix(fname, d) {
				found = true
			}
		}
		if found {
			level = level + 1
			fname, method, ln = ers.FindCaller(level)
		}
	}

	return fname, method, ln
}
func FindCaller(skip int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 10; i++ {
		pc, file, line = caller(skip + i)
		if !strings.HasPrefix(file, "nglog") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}
func caller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return pc, file, line
}

func HomeFolder() string {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("home dir not defined %+v", err)
	}
	return home
}
