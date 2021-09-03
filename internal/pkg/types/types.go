package types

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

type keyXML string
// FieldMap allows customization of the key names for default fields.
type XMLParamMap map[keyXML]string
type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}
func (f XMLParamMap) Find(key keyXML) string {
	if k, ok := f[key]; ok {
		return k
	}
	return string(key)
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