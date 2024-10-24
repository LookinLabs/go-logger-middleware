package main

import (
	"encoding/json"
	"testing"

	logger "github.com/lookinlabs/go-logger-middleware"

	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

// Small struct for benchmarking
var smallData = map[string]interface{}{
	"username": "john_doe",
	"age":      30,
	"balance":  1234.56,
	"active":   true,
	"email":    "john_doe@example.com",
	"address":  "123 Main St",
	"phone":    "123-456-7890",
	"city":     "Metropolis",
	"country":  "Freedonia",
}

// Medium struct for benchmarking
var mediumData = map[string]interface{}{
	"username":    "john_doe",
	"age":         30,
	"balance":     1234.56,
	"active":      true,
	"email":       "john_doe@example.com",
	"address":     "123 Main St",
	"phone":       "123-456-7890",
	"city":        "Metropolis",
	"country":     "Freedonia",
	"zip":         "12345",
	"state":       "CA",
	"preferences": map[string]interface{}{"newsletter": true, "sms": false},
	"history":     []interface{}{"login", "purchase", "logout"},
	"last_login":  "2023-10-01T12:34:56Z",
	"created_at":  "2023-01-01T00:00:00Z",
	"updated_at":  "2023-10-01T12:34:56Z",
	"role":        "user",
	"status":      "active",
	"bio":         "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	"website":     "https://example.com",
	"twitter":     "@johndoe",
	"linkedin":    "https://linkedin.com/in/johndoe",
	"github":      "https://github.com/johndoe",
	"projects":    []interface{}{"project1", "project2", "project3"},
}

// Large struct for benchmarking
var largeData = map[string]interface{}{
	"username":    "john_doe",
	"age":         30,
	"balance":     1234.56,
	"active":      true,
	"email":       "john_doe@example.com",
	"address":     "123 Main St",
	"phone":       "123-456-7890",
	"city":        "Metropolis",
	"country":     "Freedonia",
	"zip":         "12345",
	"state":       "CA",
	"preferences": map[string]interface{}{"newsletter": true, "sms": false},
	"history":     []interface{}{"login", "purchase", "logout"},
	"last_login":  "2023-10-01T12:34:56Z",
	"created_at":  "2023-01-01T00:00:00Z",
	"updated_at":  "2023-10-01T12:34:56Z",
	"role":        "user",
	"status":      "active",
	"bio":         "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	"website":     "https://example.com",
	"twitter":     "@johndoe",
	"linkedin":    "https://linkedin.com/in/johndoe",
	"github":      "https://github.com/johndoe",
	"projects":    []interface{}{"project1", "project2", "project3"},
	"company":     "Example Corp",
	"position":    "Software Engineer",
	"department":  "Engineering",
	"manager":     "Jane Smith",
	"office":      "HQ",
	"skills":      []interface{}{"Go", "Python", "JavaScript"},
	"languages":   []interface{}{"English", "Spanish"},
	"education":   []interface{}{"BSc Computer Science", "MSc Software Engineering"},
	"certifications": []interface{}{
		"Certified Kubernetes Administrator",
		"Certified Go Developer",
	},
	"awards": []interface{}{
		"Employee of the Month",
		"Best Innovator Award",
	},
	"publications": []interface{}{
		"Go Concurrency Patterns",
		"Building Scalable Microservices",
	},
	"conferences": []interface{}{
		"GopherCon 2023",
		"KubeCon 2023",
	},
	"volunteer_work": []interface{}{
		"Open Source Contributor",
		"Community Mentor",
	},
	"interests": []interface{}{
		"Open Source",
		"Cloud Computing",
		"Machine Learning",
	},
	"personal_projects": []interface{}{
		"Home Automation System",
		"Personal Blog",
	},
	"favorite_books": []interface{}{
		"Clean Code",
		"The Pragmatic Programmer",
	},
	"favorite_movies": []interface{}{
		"Inception",
		"The Matrix",
	},
	"favorite_music": []interface{}{
		"Rock",
		"Jazz",
	},
	"favorite_food": []interface{}{
		"Pizza",
		"Sushi",
	},
	"favorite_sports": []interface{}{
		"Soccer",
		"Basketball",
	},
}

// Benchmarking custom JSON marshalling and unmarshalling
func BenchmarkMarshal_Logger_Small(b *testing.B) {
	pairs := logger.MapToKeyValuePairs(smallData)
	for i := 0; i < b.N; i++ {
		_, err := logger.Marshal(pairs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Logger_Medium(b *testing.B) {
	pairs := logger.MapToKeyValuePairs(mediumData)
	for i := 0; i < b.N; i++ {
		_, err := logger.Marshal(pairs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Logger_Large(b *testing.B) {
	pairs := logger.MapToKeyValuePairs(largeData)
	for i := 0; i < b.N; i++ {
		_, err := logger.Marshal(pairs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Logger_Small(b *testing.B) {
	pairs := logger.MapToKeyValuePairs(smallData)
	data, _ := logger.Marshal(pairs)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := logger.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Logger_Medium(b *testing.B) {
	pairs := logger.MapToKeyValuePairs(mediumData)
	data, _ := logger.Marshal(pairs)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := logger.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Logger_Large(b *testing.B) {
	pairs := logger.MapToKeyValuePairs(largeData)
	data, _ := logger.Marshal(pairs)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := logger.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarking encoding/json

func BenchmarkMarshal_Stdlib_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(smallData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Stdlib_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(mediumData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Stdlib_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(largeData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Stdlib_Small(b *testing.B) {
	data, _ := json.Marshal(smallData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Stdlib_Medium(b *testing.B) {
	data, _ := json.Marshal(mediumData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Stdlib_Large(b *testing.B) {
	data, _ := json.Marshal(largeData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarking goccy/go-json

func BenchmarkMarshal_Goccy_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := gojson.Marshal(smallData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Goccy_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := gojson.Marshal(mediumData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Goccy_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := gojson.Marshal(largeData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Goccy_Small(b *testing.B) {
	data, _ := gojson.Marshal(smallData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := gojson.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Goccy_Medium(b *testing.B) {
	data, _ := gojson.Marshal(mediumData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := gojson.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Goccy_Large(b *testing.B) {
	data, _ := gojson.Marshal(largeData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := gojson.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarking json-iterator/go

func BenchmarkMarshal_Jsoniter_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := jsoniter.Marshal(smallData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Jsoniter_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := jsoniter.Marshal(mediumData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal_Jsoniter_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := jsoniter.Marshal(largeData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Jsoniter_Small(b *testing.B) {
	data, _ := jsoniter.Marshal(smallData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Jsoniter_Medium(b *testing.B) {
	data, _ := jsoniter.Marshal(mediumData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal_Jsoniter_Large(b *testing.B) {
	data, _ := jsoniter.Marshal(largeData)
	var result map[string]interface{}

	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}
