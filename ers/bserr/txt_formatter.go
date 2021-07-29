package bserr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/colt3k/nglog/internal/pkg/enum"

	"golang.org/x/crypto/ssh/terminal"

	log "github.com/colt3k/nglog/ng"
)

const defaultTimestampFormat = time.RFC3339

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	log.TextLayout

	isTerminal bool

	FormatValue func(value interface{}) string
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}
func (f *TextFormatter) Terminal(term bool) {
	f.isTerminal = term
}
func (f *TextFormatter) init(entry *log.LogMsg) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
	if f.FormatValue == nil {
		f.FormatValue = f.defaultFormat
	}
}
func (f *TextFormatter) Colors(enable bool) {
	f.DisableColors = enable
}
func (f *TextFormatter) DisableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *TextFormatter) EnableTimeStamp() {
	f.DisableTimestamp = true
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *log.LogMsg, disableColor bool) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	isColored := (f.isTerminal) && !f.DisableColors && !disableColor

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, timestampFormat)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		}
		if entry.Level > log.NONE {
			f.appendKeyValue(b, "level", entry.Level.String())
		}
		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message)
		}

		keys := make([]string, 0)
		for _, d := range entry.Fields {
			for k := range d {
				keys = append(keys, k)
			}
			break
		}
		if !f.DisableSorting {
			sort.Strings(keys)
		}
		for _, d := range entry.Fields {

			if entry.Level == log.ERROR {
				fmt.Fprintf(b, "\n\t")
				f.appendValue(b, fmt.Sprintf("%v", d["Method"]))
				fmt.Fprintf(b, "\n\t\t")
				f.appendValue(b, fmt.Sprintf("%v", d["File"]))

			} else {
				for _, k := range keys {
					f.appendKeyValue(b, k, d[k])
				}
			}
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *log.LogMsg, timestampFormat string) {
	var levelColorOption string
	switch entry.Level {
	case enum.DEBUG:
		levelColorOption = entry.MsgLogger().ColorDEBUG()
	case enum.INFO:
		levelColorOption = entry.MsgLogger().ColorINFO()
	case enum.WARN:
		levelColorOption = entry.MsgLogger().ColorWARN()
	case enum.ERROR, enum.FATAL:
		levelColorOption = entry.MsgLogger().ColorERR()
	default:
		levelColorOption = entry.MsgLogger().ColorDEFAULT()
	}

	var levelText = "\t "
	if entry.Level > log.NONE {
		if len(entry.Level.String()) == 4 {
			levelText = strings.ToUpper(entry.Level.String())[0:4] + " "
		} else {
			levelText = strings.ToUpper(entry.Level.String())[0:5]
		}
	}

	if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%sm%s\x1b[0m %-44s ", levelColorOption, levelText, entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%sm%s\x1b[0m[%s] %-44s ", levelColorOption, levelText, entry.Time.Format(timestampFormat), entry.Message)
	}

	keys := make([]string, 0)
	for _, d := range entry.Fields {
		for k := range d {
			keys = append(keys, k)
		}
		break
	}
	if !f.DisableSorting {
		sort.Strings(keys)
	}

	for _, d := range entry.Fields {

		if entry.Level == log.ERROR {

			fmt.Fprintf(b, "\n\t")
			f.appendValue(b, fmt.Sprintf("%v", d["Method"]))

			fmt.Fprintf(b, "\n\t\t")
			f.appendValue(b, fmt.Sprintf("%v", d["File"]))
		} else {
			for _, k := range keys {
				fmt.Fprintf(b, "\n\t\x1b[%sm%s\x1b[0m=", levelColorOption, k)
				f.appendValue(b, d[k])
			}
		}
	}
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	if key == "level" || key == "msg" {
		b.WriteString("\n" + key)
		b.WriteByte(':')
	} else {
		b.WriteString("\n\t" + key)
		b.WriteString(":\n\t\t")
	}

	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	b.WriteString(f.FormatValue(value))
}

func (f *TextFormatter) defaultFormat(value interface{}) string {
	s, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%+v", value)
	}

	if (f.QuoteEmptyVal && len(s) == 0) || NeedsQuoting(s) {
		return fmt.Sprintf("%s", s)
	}

	return s
}

func NeedsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func prefixFieldClashes(data log.Fields) {
	if t, ok := data[log.ParamTime]; ok {
		data["param.time"] = t
	}

	if m, ok := data[log.ParamMsg]; ok {
		data["param."+log.ParamMsg] = m
	}

	if l, ok := data[log.ParamLevel]; ok {
		data["param."+log.ParamLevel] = l
	}
}
