package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// BenchmarkResult holds the parsed benchmark result
type BenchmarkResult struct {
	Name  string
	Value float64
}

// CustomTicker is a custom ticker that splits long text into multiple lines
type CustomTicker struct {
	Labels        []string
	MaxLineLength int
}

// Ticks returns the ticks for the axis
func (t CustomTicker) Ticks(min, max float64) []plot.Tick {
	ticks := make([]plot.Tick, len(t.Labels))
	for i, label := range t.Labels {
		lines := splitText(label, t.MaxLineLength)
		ticks[i] = plot.Tick{
			Value: float64(i),
			Label: strings.Join(lines, "\n"),
		}
	}
	return ticks
}

// splitText splits the text into multiple lines based on the maximum line length
func splitText(txt string, maxLineLength int) []string {
	words := strings.Fields(txt)
	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 > maxLineLength {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func main() {
	// Open the benchmark results file
	file, err := os.Open("benchmark_results.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the benchmark results
	results := parseBenchmarkResults(file)

	// Create a new plot
	p := plot.New()
	p.Title.Text = "Benchmark Comparison"
	p.X.Label.Text = "Benchmark"
	p.Y.Label.Text = "Time (ns/op)"
	p.X.Tick.Label.Rotation = 1.5 // Rotate x-axis labels for better readability

	// Add the benchmark results to the plot
	values := make(plotter.Values, len(results))
	labels := make([]string, len(results))
	for i, result := range results {
		values[i] = result.Value
		labels[i] = result.Name
	}

	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		fmt.Println("Error creating bar chart:", err)
		return
	}
	p.Add(bars)

	// Set x-axis labels with custom ticker
	p.X.Tick.Marker = CustomTicker{Labels: labels, MaxLineLength: 10}

	// Save the plot to a PNG file
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "benchmark_graph.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("Benchmark graph saved to benchmark_graph.png")
}

// parseBenchmarkResults parses the benchmark results from the given file
func parseBenchmarkResults(file *os.File) []BenchmarkResult {
	var results []BenchmarkResult
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^Benchmark(\w+)-\d+\s+\d+\s+([\d.]+)\s+ns/op`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			name := matches[1]
			value, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				fmt.Println("Error parsing value:", err)
				continue
			}
			results = append(results, BenchmarkResult{Name: name, Value: value})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return results
}
