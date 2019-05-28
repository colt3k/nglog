package ng

import (
	"io"
	"os"
)

/*

Mail
Storage
RollingFile
Socket
SSL
Syslog
MQ
Rewrite
HTTP

----- DONE BELOW HERE
Fail over wraps others to be used as a backup i.e. File and fail over to Console
Console	(default)
File
*/

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
	//Layout		default format "%m%n"
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
