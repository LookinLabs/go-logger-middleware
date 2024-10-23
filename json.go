package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// KeyValuePair represents a key-value pair.
type KeyValuePair struct {
	Key   string
	Value interface{}
}

// Marshal is a custom JSON marshalling function that preserves the order of keys.
func Marshal(pairs []KeyValuePair) ([]byte, error) {
	// Estimate buffer size to reduce reallocations
	buf := bytes.NewBuffer(make([]byte, 0, len(pairs)*32))
	buf.WriteByte('{')
	first := true
	for _, pair := range pairs {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.WriteString(`"` + pair.Key + `":`)
		switch v := pair.Value.(type) {
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
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return nil
}
