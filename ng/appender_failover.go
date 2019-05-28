package ng

import (
	"fmt"
	"time"
)

//******************* FAIL OVER APPENDER ********************
type FailoverAppender struct {
	RetryIntervalSeconds int // default 60 sec
	filter               string
	primary              Appender
	primaryEnable        bool
	primaryTimeout       time.Time
	failovers            []Appender
	disableColor         bool
}

func NewFailoverAppender(primary Appender, failovers []Appender) *FailoverAppender {
	t := new(FailoverAppender)
	t.primary = primary
	t.failovers = failovers
	t.RetryIntervalSeconds = 60

	//ENABLE NEXT TWO, TO TESTING FAILOVER
	//t.primaryEnable = false
	//t.primaryTimeout = time.Now()

	return t
}
func (f *FailoverAppender) Name() string {
	return fmt.Sprintf("%T", f)
}

func (f *FailoverAppender) Applicable(msg string) bool {

	if f.primaryEnabled() {
		return f.primary.Applicable(msg)
	} else {
		for _, d := range f.failovers {
			return d.Applicable(msg)
		}
	}
	return false
}
func (f *FailoverAppender) Process(msg []byte) {
	if f.primaryEnabled() {
		f.primary.Process(msg)
	} else {
		for _, d := range f.failovers {
			d.Process(msg)
		}
	}
}

func (f *FailoverAppender) primaryEnabled() bool {

	if f.primaryEnable || (!f.primaryEnable && (time.Since(f.primaryTimeout).Seconds() > float64(f.RetryIntervalSeconds))) {
		return true
	} else {
		return false
	}
}
func (f *FailoverAppender) DisableColor() bool {
	return f.disableColor
}
