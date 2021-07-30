package ng

import (
	"bytes"
	"fmt"
	"github.com/colt3k/nglog"
	"github.com/colt3k/nglog/enum"
	"github.com/colt3k/nglog/internal/pkg/types"
	"github.com/colt3k/nglog/internal/pkg/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

type TextLayout struct {
	//Force Color
	ForceColor bool
	// Force disabling colors.
	DisableColors bool
	// Whether the logger's out is to a terminal
	isTerminal bool
	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string // The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting   bool
	QuoteEmptyVal    bool
	DisableTimestamp bool
	sync.Once
}

func (f *TextLayout) Colors(enable bool) {
	f.DisableColors = enable
}
func (f *TextLayout) DisableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *TextLayout) EnableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *TextLayout) Description() string {
	var byt bytes.Buffer
	byt.WriteString("TextLayout:")
	byt.WriteString("\nForceColor: " + strconv.FormatBool(f.ForceColor))
	byt.WriteString("\nDisableColors: " + strconv.FormatBool(f.DisableColors))
	byt.WriteString("\nTimestampFormat: " + f.TimestampFormat)
	byt.WriteString("\nDisableSorting: " + strconv.FormatBool(f.DisableSorting))
	byt.WriteString("\nQuoteEmptyVal: " + strconv.FormatBool(f.QuoteEmptyVal))
	byt.WriteString("\nDisableTimestamp: " + strconv.FormatBool(f.DisableTimestamp))
	return byt.String()
}
func (f *TextLayout) init(entry nglog.Msg) {
	if entry.MsgLogger() != nil {
		f.isTerminal = !NotTerminal
	}
}
func (f *TextLayout) Format(entry nglog.Msg, disableColor bool) ([]byte, error) {
	var b *bytes.Buffer

	b = &bytes.Buffer{}

	f.Do(func() { f.init(entry) })

	isColored := (f.ForceColor || f.isTerminal) && !f.DisableColors && !disableColor

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, entry.MsgFields(), timestampFormat)
	} else {
		if !f.DisableTimestamp {
			util.AppendKeyValue(b, "time", entry.MsgTime().Format(timestampFormat), f.QuoteEmptyVal)
		}
		if entry.MsgLevel() > enum.NONE {
			util.AppendKeyValue(b, "level", entry.MsgLevel().String(), f.QuoteEmptyVal)
		}
		if entry.MessageStr() != "" {
			util.AppendKeyValue(b, "msg", entry.MessageStr(), f.QuoteEmptyVal)
		}
		for _, d := range entry.MsgFields() {
			for k, v := range d {
				util.AppendKeyValue(b, k, v, f.QuoteEmptyVal)
			}
		}

	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextLayout) printColored(b *bytes.Buffer, entry nglog.Msg, keys []types.Fields, timestampFormat string) {
	var levelColorOption string

	switch entry.MsgLevel() {
	case enum.DBGL3:
		levelColorOption = entry.MsgLogger().ColorDEBUGL3()
	case enum.DBGL2:
		levelColorOption = entry.MsgLogger().ColorDEBUGL2()
	case enum.DEBUG:
		levelColorOption = entry.MsgLogger().ColorDEBUG()
	case enum.INFO:
		levelColorOption = entry.MsgLogger().ColorINFO()
	case enum.WARN:
		levelColorOption = entry.MsgLogger().ColorWARN()
	case enum.ERROR, enum.FATAL, enum.FATALNOEXIT:
		levelColorOption = entry.MsgLogger().ColorERR()
	default:
		levelColorOption = entry.MsgLogger().ColorDEFAULT()
	}
	var levelText = "\t "
	if entry.MsgLevel() > enum.NONE {
		if len(entry.MsgLevel().String()) == 4 {
			levelText = strings.ToUpper(entry.MsgLevel().String())[0:4] + " "
		} else {
			levelText = strings.ToUpper(entry.MsgLevel().String())[0:5]
		}
	}

	if f.DisableTimestamp {
		if len(strings.TrimSpace(levelText)) > 0 {
			fmt.Fprintf(b, "%s%s%s %-44s ", levelColorOption, levelText, CLRRESET, entry.MessageStr())
		} else {
			fmt.Fprintf(b, "%-44s ", entry.MessageStr())
		}
	} else {
		if len(strings.TrimSpace(levelText)) > 0 {
			fmt.Fprintf(b, "%s%s%s [%s] %-44s ", levelColorOption, levelText, CLRRESET, entry.MsgTime().Format(timestampFormat), entry.MessageStr())
		} else {
			fmt.Fprintf(b, "[%s] %-44s ", entry.MsgTime().Format(timestampFormat), entry.MessageStr())
		}
	}

	for _, d := range keys {
		for k, v := range d {
			fmt.Fprintf(b, " %s%s%s =", levelColorOption, k, CLRRESET)
			util.AppendValue(b, v, f.QuoteEmptyVal)
		}
	}
	//for _, k := range keys {
	//	v := entry.Fields[k]
	//	fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColorOption, k)
	//	util.AppendValue(b, v, f.QuoteEmptyVal)
	//}
}

func prefixFieldClashes(data types.Fields) {
	if t, ok := data[ParamTime]; ok {
		data["param."+ParamTime] = t
	}

	if m, ok := data[ParamMsg]; ok {
		data["param."+ParamMsg] = m
	}

	if l, ok := data[ParamLevel]; ok {
		data["param."+ParamLevel] = l
	}
}
