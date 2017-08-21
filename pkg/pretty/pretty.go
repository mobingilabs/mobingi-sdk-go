package pretty

import (
	"bytes"
	"encoding/json"
)

var Pad int = 2

func Indent(count int) string {
	pad := ""
	for i := 0; i < count; i++ {
		pad += " "
	}

	return pad
}

// JSON returns a prettified JSON string of `v`.
func JSON(v interface{}, indent int) string {
	var out bytes.Buffer

	pad := Indent(indent)
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	err = json.Indent(&out, b, "", pad)
	if err != nil {
		return err.Error()
	}

	return out.String()
}
