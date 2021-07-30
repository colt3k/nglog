package nglog

import (
	"github.com/colt3k/nglog/enum"
	"github.com/colt3k/nglog/internal/pkg/types"
	"time"
)

type Msg interface {
	WithFields(fields []types.Fields) Msg
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	LogEnt(level enum.LogLevel, format, caller string, rtn bool, args ...interface{})
	MsgLogger() Logger
	MsgFields() []types.Fields
	MsgLevel() enum.LogLevel
	MessageStr() string
	MsgTime() time.Time
}
