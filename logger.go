package nglog

/*
the enumeration application is in https://go.collins-tech.com/coltek/gogenerate.git
Build: go build -o enumeration cmd/enumeration/main.go
Move: mv enumeration $GOPATH/bin
*/
import (
	"github.com/colt3k/nglog/enum"
	"github.com/colt3k/nglog/internal/pkg/types"
)

//go:generate enumeration -pkg enum -type LogLevel -list NONE,FATAL,FATALNOEXIT,ERROR,WARN,INFO,DEBUG,DBGL2,DBGL3
//go:generate enumeration -pkg enum -type Flags -list Ldate,Ltime,Lmicroseconds,Llongfile,Lshortfile,LUTC,LstdFlags -listval shift

type Logger interface {
	Flags() int
	Level() enum.LogLevel
	// Logln lvl can be a string of one of the levels above or an enum value of the same
	Logln(lvl interface{}, args ...interface{})
	// Logf lvl can be a string of one of the levels above or an enum value of the same
	Logf(lvl interface{}, format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	PrintTypeOfValue(arg interface{})
	PrintStructWithFieldNames(arg interface{})
	PrintStructWithFieldNamesIndent(arg interface{}, indent bool)
	PrintGoSyntaxOfValue(arg interface{})
	Println(args ...interface{})
	SetLevel(level enum.LogLevel)
	SetFlags(flg enum.Flags)
	SetFormatter(formatter Layout)
	WithFields(fields []types.Fields) Msg
	ColorDEBUGL3() string
	ColorDEBUGL2() string
	ColorDEBUG() string
	ColorINFO() string
	ColorWARN() string
	ColorERR() string
	ColorDEFAULT() string
}
