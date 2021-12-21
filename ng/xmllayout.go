package ng

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
)

type keyXML string

// FieldMap allows customization of the key names for default fields.
type XMLParamMap map[keyXML]string
type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// Default key names for the default fields
const (
	XMLParamMsg   = "msg"
	XMLParamLevel = "level"
	XMLParamTime  = "time"
)

func (f XMLParamMap) Find(key keyXML) string {
	if k, ok := f[key]; ok {
		return k
	}
	return string(key)
}

type XMLLayout struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool
	DisableQuoting   bool

	ParamMap XMLParamMap
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
	f.DisableTimestamp = false
}
func (f *XMLLayout) DisableTextQuoting() {
	f.DisableQuoting = true
}
func (f *XMLLayout) EnableTextQuoting() {
	f.DisableQuoting = false
}
func (f *XMLLayout) Format(entry *LogMsg, disableColor bool) ([]byte, error) {
	data := make(Fields, len(entry.Fields)+3)

	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data[f.ParamMap.Find(XMLParamTime)] = entry.Time.Format(timestampFormat)
	}
	data[f.ParamMap.Find(XMLParamMsg)] = entry.Message
	data[f.ParamMap.Find(XMLParamLevel)] = entry.Level.String()

	serialized, err := xml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to XML\n%+v", err)
	}
	return append(serialized, '\n'), nil
}

// StringMap marshals into XML.
func (m Fields) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, i := range m {
		switch v := i.(type) {
		case int:
			e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: strconv.Itoa(v)})
		case string:
			e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
		default:
			fmt.Println("unknown type")
		}
	}

	return e.EncodeToken(start.End())
}
