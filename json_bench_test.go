package logger

import (
	"testing"
)

func BenchmarkMarshal(b *testing.B) {
	data := map[string]interface{}{
		"username": "john_doe",
		"age":      30,
		"balance":  1234.56,
		"active":   true,
	}

	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	data := []byte(`{"username":"john_doe","age":30,"balance":1234.56,"active":true}`)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}
