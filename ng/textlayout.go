package ng

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"strconv"

	"github.com/colt3k/nglog/internal/pkg/enum"
	"github.com/colt3k/nglog/internal/pkg/util"
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

type Layout interface {
	Format(*LogMsg, bool) ([]byte, error)
	Description() string
	Colors(bool)
	DisableTimeStamp()
	EnableTimeStamp()
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
func (f *TextLayout) init(entry *LogMsg) {
	if entry.Logger != nil {
		f.isTerminal = !NotTerminal
	}
}
func (f *TextLayout) Format(entry *LogMsg, disableColor bool) ([]byte, error) {
	var b *bytes.Buffer

	b = &bytes.Buffer{}

	f.Do(func() { f.init(entry) })

	isColored := (f.ForceColor || f.isTerminal) && !f.DisableColors && !disableColor

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, entry.Fields, timestampFormat)
	} else {
		if !f.DisableTimestamp {
			util.AppendKeyValue(b, "time", entry.Time.Format(timestampFormat), f.QuoteEmptyVal)
		}
		if entry.Level > enum.NONE {
			util.AppendKeyValue(b, "level", entry.Level.String(), f.QuoteEmptyVal)
		}
		if entry.Message != "" {
			util.AppendKeyValue(b, "msg", entry.Message, f.QuoteEmptyVal)
		}
		for _, d := range entry.Fields {
			for k, v := range d {
				util.AppendKeyValue(b, k, v, f.QuoteEmptyVal)
			}
		}

	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextLayout) printColored(b *bytes.Buffer, entry *LogMsg, keys []Fields, timestampFormat string) {
	var levelColorOption string

	switch entry.Level {
	case enum.DBGL3:
		levelColorOption = entry.Logger.ColorDEBUGL3
	case enum.DBGL2:
		levelColorOption = entry.Logger.ColorDEBUGL2
	case enum.DEBUG:
		levelColorOption = entry.Logger.ColorDEBUG
	case enum.INFO:
		levelColorOption = entry.Logger.ColorINFO
	case enum.WARN:
		levelColorOption = entry.Logger.ColorWARN
	case enum.ERROR, enum.FATAL, enum.FATALNOEXIT:
		levelColorOption = entry.Logger.ColorERR
	default:
		levelColorOption = entry.Logger.ColorDEFAULT
	}
	var levelText = "\t "
	if entry.Level > enum.NONE {
		if len(entry.Level.String()) == 4 {
			levelText = strings.ToUpper(entry.Level.String())[0:4] + " "
		} else {
			levelText = strings.ToUpper(entry.Level.String())[0:5]
		}
	}

	if f.DisableTimestamp {
		if len(strings.TrimSpace(levelText)) > 0 {
			fmt.Fprintf(b, "%s%s%s %-44s ", levelColorOption, levelText, CLRRESET, entry.Message)
		} else {
			fmt.Fprintf(b, "%-44s ", entry.Message)
		}
	} else {
		if len(strings.TrimSpace(levelText)) > 0 {
			fmt.Fprintf(b, "%s%s%s [%s] %-44s ", levelColorOption, levelText, CLRRESET, entry.Time.Format(timestampFormat), entry.Message)
		} else {
			fmt.Fprintf(b, "[%s] %-44s ", entry.Time.Format(timestampFormat), entry.Message)
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

func prefixFieldClashes(data Fields) {
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
