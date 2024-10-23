package logger

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   []KeyValuePair
		want    string
		wantErr bool
	}{
		{
			name: "valid key-value pairs",
			input: []KeyValuePair{
				{"name", "John"},
				{"age", 30},
			},
			want:    `{"name":"John","age":30}`,
			wantErr: false,
		},
		{
			name: "valid key-value pairs with float",
			input: []KeyValuePair{
				{"name", "John"},
				{"age", 30.5},
			},
			want:    `{"name":"John","age":30.500000}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:  "valid JSON",
			input: `{"name":"John","age":30}`,
			want: map[string]interface{}{
				"name": "John",
				"age":  float64(30), // json.Unmarshal decodes numbers as float64
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"name":"John","age":30`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]interface{}
			err := Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equal(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two maps
func equal(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
