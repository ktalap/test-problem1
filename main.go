package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func findDifferences(file1path, file2path, output1path, output2path string) error {
	startTime := time.Now()

	// A map like this has O(1) complexity
	// It's constant time to check if a key exists
	// Regardless of the size of the map
	file1Lines := make(map[string]bool, 10000) // pre-allocated capacity
	file1, err := os.Open(file1path)
	if err != nil {
		return fmt.Errorf("error opening file1: %v", err)
	}
	// Making sure the file is closed on the exit of the function
	defer file1.Close()

	// Read all lines from file1
	// scanner is like a cursor that moves through the file
	// line by line
	scanner := bufio.NewScanner(file1)
	// Increase scanner buffer size
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	// This loop reads file line by line
	for scanner.Scan() {
		// scanner.Text gets the current line as string
		// all we care about is existence of lines
		// here we are writing boolean values to the map
		// we defined map[string]bool
		file1Lines[scanner.Text()] = true
		// returns false if it's the end of the file or error
	}
	// By the end of the loop we have something similar to this
	// file1Lines = {
	// 	"banana": true,
	// 	"date": true,
	// 	"fig": true
	// }

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file1: %v", err)
	}

	// Read second file and compare
	file2Lines := make(map[string]bool, 1000) // pre-allocated capacity
	file2, err := os.Open(file2path)
	if err != nil {
		return fmt.Errorf("error opening file2: %v", err)
	}
	defer file2.Close()

	// Read all the lines from file2
	// scanner is like a cursor that moves through the file
	// line by line
	scanner = bufio.NewScanner(file2)
	// Increase scanner buffer size
	scanner.Buffer(buf, 1024*1024)
	// This loop reads file line by line
	for scanner.Scan() {
		// scanner.Text gets the current line as string
		// all we care about is existence of lines
		// here we are writing boolean values to the map
		// we defined map[string]bool
		file2Lines[scanner.Text()] = true
		// returns false if it's the end of the file or error
	}
	// By the end of the loop we have something similar to this
	// file2Lines = {
	// 	"banana": true,
	// 	"date": true,
	// 	"fig": true
	// }

	// Create the output files with buffered writers
	output1, err := os.Create(output1path)
	if err != nil {
		return fmt.Errorf("error creating output1: %v", err)
	}
	defer output1.Close()
	writer1 := bufio.NewWriter(output1)
	defer writer1.Flush()

	output2, err := os.Create(output2path)
	if err != nil {
		return fmt.Errorf("error creating output2: %v", err)
	}
	defer output2.Close()
	writer2 := bufio.NewWriter(output2)
	defer writer2.Flush()

	// Process in batches
	const batchSize = 1000
	batch := make([]string, 0, batchSize)

	// Find the strings unique to file1
	for line := range file1Lines {
		if !file2Lines[line] {
			// Here we have lines that are unique to file 1
			batch = append(batch, line)
			// A check for batch size
			if len(batch) >= batchSize {
				// Combining new strings with newlines between
				fmt.Fprint(output1, strings.Join(batch, "\n")+"\n")
				batch = batch[:0]
			}
		}
	}
	// Final cleanup, and saving whatever elements went past the
	// batch size
	if len(batch) > 0 {
		fmt.Fprint(output1, strings.Join(batch, "\n")+"\n")
		batch = batch[:0]
	}

	// Find string unique to file2
	for line := range file2Lines {
		if !file1Lines[line] {
			batch = append(batch, line)
			if len(batch) >= batchSize {
				fmt.Fprint(output2, strings.Join(batch, "\n")+"\n")
				batch = batch[:0]
			}
		}
	}
	if len(batch) > 0 {
		fmt.Fprint(output2, strings.Join(batch, "\n")+"\n")
	}

	fmt.Printf("Total time: %v\n", time.Since(startTime))
	return nil
}

func main() {
	// Start time tracking
	startTime := time.Now()

	err := findDifferences("large_input1.txt", "large_input2.txt", "only_in_file1.txt", "only_in_file2.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Calculate execution time
	duration := time.Since(startTime)

	fmt.Println("Files processed successfully")
	fmt.Printf("Execution time: %v\n", duration)
}
