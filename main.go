package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func findDifferences(file1path, file2path, output1path, output2path string) error {
	// A map like this has O(1) complexity
	// It's constant time to check if a key exists
	// Regardless of the size of the map
	file1Lines := make(map[string]bool)
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

	// Read second file and compare
	file2Lines := make(map[string]bool)
	file2, err := os.Open(file2path)
	if err != nil {
		return fmt.Errorf("error opening file2: %v", err)
	}
	defer file2.Close()

	// Read all the lines from file2
	// scanner is like a cursor that moves through the file
	// line by line
	scanner = bufio.NewScanner(file2)
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

	// Create the output files
	output1, err := os.Create(output1path)
	if err != nil {
		return fmt.Errorf("error creating output1: %v", err)
	}
	defer output1.Close()

	output2, err := os.Create(output2path)
	if err != nil {
		return fmt.Errorf("error creating output2: %v", err)
	}
	defer output2.Close()

	// Find the strings unique to file1
	for line := range file1Lines {
		if !file2Lines[line] {
			fmt.Fprintln(output1, line)
		}
	}

	// Find string unique to file2
	for line := range file2Lines {
		if !file1Lines[line] {
			fmt.Fprintln(output2, line)
		}
	}

	return nil
}

func main() {
	// Start time tracking
	startTime := time.Now()

	err := findDifferences("test-set-1/input1.txt", "test-set-1/input2.txt", "test-set-1/only_in_file1.txt", "test-set-1/only_in_file2.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = findDifferences("test-set-2/input1.txt", "test-set-2/input2.txt", "test-set-2/only_in_file1.txt", "test-set-2/only_in_file2.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Calculate execution time
	duration := time.Since(startTime)

	fmt.Println("Files processed successfully")
	fmt.Printf("Execution time: %v\n", duration)
}
