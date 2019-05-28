package ng

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type key string

// FieldMap allows customization of the key names for default fields.
type ParamMap map[key]string

// Default key names for the default fields
const (
	ParamMsg   = "msg"
	ParamLevel = "level"
	ParamTime  = "time"
)

func (f ParamMap) Find(key key) string {
	if k, ok := f[key]; ok {
		return k
	}
	return string(key)
}

type JSONLayout struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	ParamMap ParamMap
}

func (f *JSONLayout) Description() string {
	var byt bytes.Buffer
	byt.WriteString("JSONLayout:")
	byt.WriteString("\nTimestampFormat: " + f.TimestampFormat)
	byt.WriteString("\nDisableTimestamp: " + strconv.FormatBool(f.DisableTimestamp))
	return byt.String()
}
func (f *JSONLayout) Colors(enable bool) {

}
func (f *JSONLayout) DisableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *JSONLayout) EnableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *JSONLayout) Format(entry *LogMsg, disableColor bool) ([]byte, error) {
	data := make(Fields, len(entry.Fields)+3)
	//@TODO uncomment and fix
	//for k, v := range entry.Fields {
	//	switch v := v.(type) {
	//	case error:
	//		// Otherwise errors are ignored by `encoding/json`
	//		// https://github.com/sirupsen/logrus/issues/137
	//		data[k] = v.Error()
	//	default:
	//		data[k] = v
	//	}
	//}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data[f.ParamMap.Find(ParamTime)] = entry.Time.Format(timestampFormat)
	}
	data[f.ParamMap.Find(ParamMsg)] = entry.Message
	data[f.ParamMap.Find(ParamLevel)] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON\n%+v", err)
	}
	return append(serialized, '\n'), nil
}
