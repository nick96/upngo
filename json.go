package upngo

import (
	"bytes"
	"encoding/json"
)

// unmarshal is a custom JSON unmarshaller that differs from the standard
// library's JSON.Unmarshal in that it does not allow unknown fields.
func unmarshal(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}
