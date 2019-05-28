package ng

import (
	"fmt"
	"strings"
)

//******************* CONSOLE APPENDER ********************
type ConsoleAppender struct {
	*OutAppender
}

func NewConsoleAppender(filter string) *ConsoleAppender {
	oa := newOutAppender(filter, "")
	t := new(ConsoleAppender)
	t.OutAppender = oa

	return t
}
func (c *ConsoleAppender) Name() string {
	if len(c.name) > 0 {
		return c.name
	}
	return fmt.Sprintf("%T", c)
}
func (c *ConsoleAppender) DisableColor() bool {
	return c.disableColor
}
func (c *ConsoleAppender) Applicable(filter string) bool {
	if c.filter == "*" {
		return true
	}
	if strings.Index(filter, c.filter) > -1 {
		return true
	}
	return false
}

func (c *ConsoleAppender) Process(msg []byte) {
	c.Out.Write(msg)
}
