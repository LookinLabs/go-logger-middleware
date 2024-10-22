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

// Unmarshal is a custom JSON unmarshalling function.
func Unmarshal(data []byte, value interface{}) error {
	result := make(map[string]interface{})
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	mapPtr, ok := value.(*map[string]interface{})
	if !ok {
		return fmt.Errorf("expected *map[string]interface{}")
	}
	*mapPtr = result
	return nil
}
