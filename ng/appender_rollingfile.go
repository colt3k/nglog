package ng

import (
	"fmt"
	"strings"
)

type RollingFileAppender struct {
	*FileAppender
	triggerPolicy TriggerPolicy // Triggering Policy to roll file
	strategy      Strategy      // Strategy of archives
	filePattern   string        // 	i.e. logs/$${date:yyyy-MM}/app-%d{yyyy-MM-dd-HH}-%i.log.gz
}

func NewRollingFileAppender(filter, fileName, name string, bufferSize int, trigger TriggerPolicy, strategy Strategy) (*RollingFileAppender,error) {
	fa,err := NewFileAppender(filter, fileName, name, bufferSize)
	if err != nil {
		return nil, err
	}
	t := new(RollingFileAppender)
	t.FileAppender = fa
	t.triggerPolicy = trigger
	t.strategy = strategy

	return t,nil
}

func (r *RollingFileAppender) Name() string {
	if len(r.name) > 0 {
		return r.name
	}
	return fmt.Sprintf("%T", r)
}

func (r *RollingFileAppender) Applicable(filter string) bool {
	if r.filter == "*" {
		return true
	}
	if strings.Index(filter, r.filter) > -1 {
		return true
	}
	return false
}

// TODO add logging for rolling
func (r *RollingFileAppender) Process(msg []byte) {
	if r.triggerPolicy != nil {
		b, err := r.triggerPolicy.Rotate(r.fileName, r.strategy)
		if err != nil {
			tmpMsg := fmt.Sprintf("issue rotating log file %v \n", err)
			t := string(msg)
			msg = []byte(tmpMsg+t)
		}
		if b {
			// re-establish connection to new log file
		}
	}

	if r.buffered {
		r.buf.Write(msg)
		if r.immediateFlush {
			r.buf.Flush()
		}
	} else {
		r.Out.Write(msg)
	}
}
func (r *RollingFileAppender) DisableColor() bool {
	return r.disableColor
}
