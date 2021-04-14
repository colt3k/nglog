package nglog

/*
test
*/
import (
	"github.com/colt3k/nglog/internal/pkg/enum"
	"github.com/colt3k/nglog/ng"
)

//go:generate enumeration -pkg enum -type LogLevel -list NONE,FATAL,ERROR,WARN,INFO,DEBUG,DEBUGX2
//go:generate enumeration -pkg enum -type Flags -list Ldate,Ltime,Lmicroseconds,Llongfile,Lshortfile,LUTC,LstdFlags -listval shift

type Logger interface {
	Flags() int
	Level() enum.LogLevel
	Logln(lvl enum.LogLevel, args ...interface{})
	Logf(lvl enum.LogLevel, format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	PrintTypeOfValue(arg interface{})
	PrintStructWithFieldNames(arg interface{})
	PrintStructWithFieldNamesIndent(arg interface{}, indent bool)
	PrintGoSyntaxOfValue(arg interface{})
	Println(args ...interface{})
	SetLevel(level enum.LogLevel)
	SetFlags(flg enum.Flags)
	SetFormatter(formatter ng.Layout)
	WithFields(fields []ng.Fields) *ng.LogMsg
}
