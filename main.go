package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// To organize stats
type Stats struct {
	TotalLinesFile1 int
	TotalLinesFile2 int
	UniqueInFile1   int
	UniqueInFile2   int
	ProcessingTime  time.Duration
	ReadingTime     time.Duration
	WritingTime     time.Duration
}

// Algorithmic change from maps to struct of map+slice
type StringSet struct {
	data map[string]bool // For O(1) lookup
	keys []string        // A slice to maintain insertion order
}

// A slice i.e. ["apple", "banana"]
// append("cherry") -> ["apple", "banana", "cherry"]

// Functions of the struct
func (s *StringSet) Add(str string) {
	if !s.data[str] {
		s.data[str] = true
		s.keys = append(s.keys, str)
	}
}
func (s *StringSet) Sort() {
	sort.Strings(s.keys)
}
func (s *StringSet) Len() int {
	return len(s.keys)
}

// Main function to find differences
func findDifferences(file1path, file2path, output1path, output2path string) (*Stats, error) {
	stats := &Stats{}
	startTime := time.Now()

	// Initialize the sets
	file1Set := &StringSet{
		data: make(map[string]bool, 10000),
		keys: make([]string, 0, 10000),
	}
	file2Set := &StringSet{
		data: make(map[string]bool, 10000),
		keys: make([]string, 0, 10000),
	}

	// Read files
	readStart := time.Now()

	// Read first file
	if err := readFileIntoSet(file1path, file1Set); err != nil {
		return nil, err
	}
	stats.TotalLinesFile1 = file1Set.Len()

	// Read second file
	if err := readFileIntoSet(file2path, file2Set); err != nil {
		return nil, err
	}

	// Recording the stats
	stats.TotalLinesFile2 = file2Set.Len()
	stats.ReadingTime = time.Since(readStart)

	// Sort both of the sets
	file1Set.Sort()
	file2Set.Sort()

	// Create the output files with buffered writers
	output1, err := os.Create(output1path)
	if err != nil {
		return nil, err
	}
	defer output1.Close()
	writer1 := bufio.NewWriter(output1)
	defer writer1.Flush()

	output2, err := os.Create(output2path)
	if err != nil {
		return nil, err
	}
	defer output2.Close()
	writer2 := bufio.NewWriter(output2)
	defer writer2.Flush()

	// Below we write differences using strings.Builder
	writeStart := time.Now()

	var builder strings.Builder
	builder.Grow(1024 * 1024) // Pre-allocating 1MB

	// Find strings unique to file1
	for _, line := range file1Set.keys { // Loop through slice of strings
		if !file2Set.data[line] { // Quick lookup in map
			// Here line was not found in file2Set, it's unique
			builder.WriteString(line)
			builder.WriteByte('\n')
			stats.UniqueInFile1++

			// Flush if builder gets too large
			if builder.Len() > 1000000 { // // If builder has more than 1MB of data
				writer1.WriteString(builder.String()) // Write everything to file
				builder.Reset()                       // Clear the builder
			}
		}
	}
	// Example of how the loop works
	// file1Set.keys is our slice ["apple", "banana", "cherry"]
	// file2Set.data is our map {"banana": true, "date": true}
	// First iteration:
	// line = "apple"
	// !file2Set.data["apple"] is true  (apple not in map)
	// Second iteration:
	// line = "banana"
	// !file2Set.data["banana"] is false (banana is in map)
	// This is the final cleanup, if we don't get to 1MB, and that
	// is the final iteration we just write the remaning lines
	if builder.Len() > 0 {
		writer1.WriteString(builder.String())
		builder.Reset()
	}

	// Find strings unique to file2
	for _, line := range file2Set.keys {
		if !file1Set.data[line] {
			builder.WriteString(line)
			builder.WriteByte('\n')
			stats.UniqueInFile2++

			// Flush if builder gets too large
			if builder.Len() > 1000000 { // 1MB
				writer2.WriteString(builder.String())
				builder.Reset()
			}
		}
	}
	// Example of how the loop works
	// file1Set.keys is our slice ["apple", "banana", "cherry"]
	// file2Set.data is our map {"banana": true, "date": true}
	// First iteration:
	// line = "apple"
	// !file2Set.data["apple"] is true  (apple not in map)
	// Second iteration:
	// line = "banana"
	// !file2Set.data["banana"] is false (banana is in map)
	if builder.Len() > 0 {
		writer2.WriteString(builder.String())
	}

	// Recording the stats
	stats.WritingTime = time.Since(writeStart)
	stats.ProcessingTime = time.Since(startTime)

	return stats, nil
}

func readFileIntoSet(filepath string, set *StringSet) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Pre-allocate the buffer
	scanner.Buffer(make([]byte, 256*1024), 1024*1024) // 256KB initial, 1MB max

	// Read line by line
	for scanner.Scan() {
		// Add the line to the set
		set.Add(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

func main() {
	// Start time tracking
	startTime := time.Now()

	stats, err := findDifferences("large_input1.txt", "large_input2.txt",
		"only_in_file1.txt", "only_in_file2.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nProcessing Statistics:")
	fmt.Printf("Total lines in file 1: %d\n", stats.TotalLinesFile1)
	fmt.Printf("Total lines in file 2: %d\n", stats.TotalLinesFile2)
	fmt.Printf("Unique lines in file 1: %d\n", stats.UniqueInFile1)
	fmt.Printf("Unique lines in file 2: %d\n", stats.UniqueInFile2)
	fmt.Printf("Reading time: %v\n", stats.ReadingTime)
	fmt.Printf("Writing time: %v\n", stats.WritingTime)
	fmt.Printf("Total processing time: %v\n", stats.ProcessingTime)
	fmt.Printf("Total program time: %v\n", time.Since(startTime))
}
