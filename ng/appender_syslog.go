package ng

import (
	"fmt"
	"log/syslog"
	"strings"
)

//******************* CONSOLE APPENDER ********************
type SyslogAppender struct {
	*OutAppender
}

func NewSyslogAppender(filter, programName string) (*SyslogAppender, error) {
	oa := newOutAppender(filter, "")
	t := new(SyslogAppender)
	t.OutAppender = oa
	logwriter, err := syslog.New(syslog.LOG_NOTICE, programName)
	if err != nil {
		return nil, fmt.Errorf("issue creating syslog appender\n%+v", err)
	}
	t.Out = logwriter

	return t, nil
}
func (c *SyslogAppender) Name() string {
	if len(c.name) > 0 {
		return c.name
	}
	return fmt.Sprintf("%T", c)
}
func (c *SyslogAppender) DisableColor() bool {
	return c.disableColor
}
func (c *SyslogAppender) Applicable(filter string) bool {
	if c.filter == "*" {
		return true
	}
	if strings.Index(filter, c.filter) > -1 {
		return true
	}
	return false
}

func (c *SyslogAppender) Process(msg []byte) {
	c.Out.Write(msg)
}
