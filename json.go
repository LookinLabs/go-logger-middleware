package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Marshal is a custom JSON marshalling function.
func Marshal(value interface{}) ([]byte, error) {
	data, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{}")
	}

	// Estimate buffer size to reduce reallocations
	buf := bytes.NewBuffer(make([]byte, 0, len(data)*32))
	buf.WriteByte('{')
	first := true
	for key, val := range data {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.WriteString(`"` + key + `":`)
		switch v := val.(type) {
		case string:
			buf.WriteString(`"` + v + `"`)
		case int:
			buf.WriteString(fmt.Sprintf("%d", v))
		case float64:
			buf.WriteString(fmt.Sprintf("%f", v))
		default:
			buf.WriteString(`null`)
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return nil
}
