package bserr

import (
	"strconv"
	"strings"

	"github.com/colt3k/nglog/ers"
	log "github.com/colt3k/nglog/ng"
)

var (
	Level  = 4
	stders = NewErr()
)

type ErrStack struct {
	fname  string
	method string
	line   int
}
type BSErr struct {
	ignore  []string
	current string
}

func NewErr() ers.Error {
	tmp := new(BSErr)
	tmp.ignore = []string{"nglog/ers/errs.go",
		"nglog/ers/bserr/errs.go", "runtime.goexit"}
	return tmp
}

func Err(e error, msg ...string) bool {
	return stders.Err(e, false, msg...)
}
func ErrSkipTrace(e error, msg ...string) bool {
	return stders.Err(e, true, msg...)
}
func NoPrintErr(e error) bool {
	return stders.NoPrintErr(e)
}
func NotErr(e error, msg ...string) bool {
	return stders.NotErr(e, msg...)
}
func NotErrSkipTrace(e error, msg ...string) bool {
	return stders.NotErrSkipTrace(e, true, msg...)
}
func StopErr(e error, msg ...string) {
	stders.StopErr(e, msg...)
}
func WarnErr(e error, msg ...string) bool {
	return stders.WarnErr(e, msg...)
}

//Err check for error and output message if one is passed in
func (er *BSErr) Err(e error, skipTrace bool, msg ...string) bool {
	if e != nil && !skipTrace {
		log.DisableTimestamp()
		errStack := er.findIssue()
		flds := make([]log.Fields, 0)
		for _, d := range errStack {
			f := d.fname + ":" + strconv.Itoa(d.line)
			fldsTmp := log.Fields{"Method": d.method, "File": f}
			flds = append(flds, fldsTmp)
		}
		entry := log.WithFields(flds)
		if len(msg) > 0 {
			entry.Error(msg[0], e)
		} else {
			entry.Error(e)
		}
		log.EnableTimestamp()
		return true
	} else if e != nil && len(msg) > 0 {
		log.Logln(log.ERROR, msg, e)
		return true
	}
	return false
}

//NoPrintErr is this an Error if so don't print, used in if/else
func (er *BSErr) NoPrintErr(e error) bool {
	if e != nil {
		return true
	}
	return false
}

//NotErr check for NO error and output message if one is passed in
func (er *BSErr) NotErr(e error, msg ...string) bool {
	return er.NotErrSkipTrace(e, false, msg...)
}

//NotErr check for NO error and output message if one is passed in
func (er *BSErr) NotErrSkipTrace(e error, skipTrace bool, msg ...string) bool {
	if e != nil && !skipTrace {

		log.DisableTimestamp()
		errStack := er.findIssue()
		flds := make([]log.Fields, 0)
		for _, d := range errStack {
			f := d.fname + ":" + strconv.Itoa(d.line)
			fldsTmp := log.Fields{"Method": d.method, "File": f}
			flds = append(flds, fldsTmp)
		}
		entry := log.WithFields(flds)
		if len(msg) > 0 {
			entry.Error(msg[0], e)
		} else {
			entry.Error(e)
		}
		log.EnableTimestamp()
		return false
	} else if e != nil && len(msg) > 0 {
		log.Logln(log.ERROR, msg, e)
		return false
	}
	return true
}

//StopErr check for error, output message if one is passed in and stop execution
func (er *BSErr) StopErr(e error, msg ...string) {
	if e != nil && msg == nil {

		log.DisableTimestamp()
		errStack := er.findIssue()
		flds := make([]log.Fields, 0)
		for _, d := range errStack {
			f := d.fname + ":" + strconv.Itoa(d.line)
			fldsTmp := log.Fields{"Method": d.method, "File": f}
			flds = append(flds, fldsTmp)
		}
		entry := log.WithFields(flds)
		entry.Fatal(e)
		log.EnableTimestamp()
	} else if e != nil && len(msg) > 0 {
		log.Logln(log.FATAL, msg, e)
	}
}

func (er *BSErr) WarnErr(e error, msg ...string) bool {
	if e != nil {
		if len(msg) > 0 {
			log.Logln(log.WARN, msg, e.Error())
		} else {
			log.Logln(log.WARN, e.Error())
		}
		return true
	}
	return false
}

func (er *BSErr) findIssue() []*ErrStack {
	level := Level
	found := false
	fname, method, ln := ers.FindCaller(Level)

	stack := make([]*ErrStack, 0)

	for _, d := range er.ignore {
		if strings.HasSuffix(fname, d) {
			found = true
		}
	}
	for {

		if !found {
			break
		}
		found = false
		skip := false
		for _, d := range er.ignore {
			if strings.HasSuffix(fname, d) {
				found = true
				skip = true
			} else if strings.HasSuffix(method, d) {
				found = true
				skip = true
			}
		}
		if len(fname) > 0 && !skip {
			t := &ErrStack{
				fname:  fname,
				method: method,
				line:   ln,
			}
			stack = append(stack, t)
		}
		if len(fname) != 0 {
			found = true
		}
		if found {
			level = level + 1
			fname, method, ln = ers.FindCaller(level)
		}

	}

	return stack
}
