package main

import (
	"os"
	"testing"
)

func TestFindDifferences(t *testing.T) {
	// Create test input files
	input1Content := "apple\nbanana\ncherry\ndate\n"
	input2Content := "banana\ndate\nfig\ngrape\n"

	err := os.WriteFile("test_input1.txt", []byte(input1Content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input1: %v", err)
	}
	defer os.Remove("test_input1.txt")

	err = os.WriteFile("test_input2.txt", []byte(input2Content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input2: %v", err)
	}
	defer os.Remove("test_input2.txt")

	// Run the function
	err = findDifferences("test_input1.txt", "test_input2.txt", "test_output1.txt", "test_output2.txt")
	if err != nil {
		t.Fatalf("findDifferences failed: %v", err)
	}
	defer os.Remove("test_output1.txt")
	defer os.Remove("test_output2.txt")

	// Read output files
	output1, err := os.ReadFile("test_output1.txt")
	if err != nil {
		t.Fatalf("Failed to read output1: %v", err)
	}

	output2, err := os.ReadFile("test_output2.txt")
	if err != nil {
		t.Fatalf("Failed to read output2: %v", err)
	}

	// Expected outputs
	expectedOutput1 := "apple\ncherry\n"
	expectedOutput2 := "fig\ngrape\n"

	// Compare results
	if string(output1) != expectedOutput1 {
		t.Errorf("Output1 incorrect\nexpected: %q\ngot: %q", expectedOutput1, string(output1))
	}

	if string(output2) != expectedOutput2 {
		t.Errorf("Output2 incorrect\nexpected: %q\ngot: %q", expectedOutput2, string(output2))
	}
}
