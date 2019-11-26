package ng

import (
	"io"
	"os"
)

// Appender
type Appender interface {
	Name() string
	//package or name to find in message, in Java this was a package
	Applicable(filter string) bool
	Process(msg []byte)
	DisableColor() bool
}

type OutAppender struct {
	name         string
	Out          io.Writer
	filter       string
	disableColor bool
}

func newOutAppender(filter, name string) *OutAppender {
	t := new(OutAppender)
	t.disableColor = false
	if t.Out == nil {
		t.Out = os.Stdout
	}
	t.filter = filter
	if len(name) > 0 {
		t.name = name
	}
	return t
}
