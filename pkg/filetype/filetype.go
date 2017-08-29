package filetype

import "encoding/json"

func IsJSON(in string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(in), &js) == nil
}
