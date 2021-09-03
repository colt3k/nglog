package ng

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/colt3k/nglog"
	"github.com/colt3k/nglog/internal/pkg/types"
	"strconv"
)


// Default key names for the default fields
const (
	XMLParamMsg   = "msg"
	XMLParamLevel = "level"
	XMLParamTime  = "time"
)

type XMLLayout struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	ParamMap types.XMLParamMap
}

func (f *XMLLayout) Description() string {
	var byt bytes.Buffer
	byt.WriteString("XMLLayout:")
	byt.WriteString("\nTimestampFormat: " + f.TimestampFormat)
	byt.WriteString("\nDisableTimestamp: " + strconv.FormatBool(f.DisableTimestamp))
	return byt.String()
}
func (f *XMLLayout) Colors(enable bool) {

}
func (f *XMLLayout) DisableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *XMLLayout) EnableTimeStamp() {
	f.DisableTimestamp = true
}
func (f *XMLLayout) Format(entry nglog.Msg, disableColor bool) ([]byte, error) {
	data := make(types.Fields, len(entry.MsgFields())+3)

	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data[f.ParamMap.Find(XMLParamTime)] = entry.MsgTime().Format(timestampFormat)
	}
	data[f.ParamMap.Find(XMLParamMsg)] = entry.MessageStr()
	data[f.ParamMap.Find(XMLParamLevel)] = entry.MsgLevel().String()

	serialized, err := xml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to XML\n %+v", err)
	}
	return append(serialized, '\n'), nil
	//return serialized, nil
}
