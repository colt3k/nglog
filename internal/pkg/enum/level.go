// Code generated by go generate
// This file was generated by robots at 2018-05-21 16:26:06.261281984 +0000 UTC
package enum

// LogLevel is an optional LogLevel
type LogLevel int

const (
	NONE LogLevel = 1 + iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
)

var loglevel = [...]string{
	"NONE", "FATAL", "ERROR", "WARN", "INFO", "DEBUG",
}

func (l LogLevel) String() string {
	return loglevel[l-1]
}

/*
Types pulls full list as []string
*/
func (l LogLevel) Types() []string {
	return loglevel[:]
}
