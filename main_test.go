package main

import (
	"os"
	"strings"
	"testing"
)

func TestStringSet(t *testing.T) {
	set := &StringSet{
		data: make(map[string]bool),
		keys: make([]string, 0),
	}

	// Test adding strings
	testStrings := []string{"banana", "apple", "cherry"}
	for _, s := range testStrings {
		set.Add(s)
	}

	// Test size
	if len(set.keys) != 3 {
		t.Errorf("Expected 3 items, got %d", len(set.keys))
	}

	// Test duplicate handling
	set.Add("apple")
	if len(set.keys) != 3 {
		t.Errorf("Duplicate added: expected 3 items, got %d", len(set.keys))
	}

	// Test sorting
	set.Sort()
	expected := []string{"apple", "banana", "cherry"}
	for i, v := range expected {
		if set.keys[i] != v {
			t.Errorf("Expected %s at position %d, got %s", v, i, set.keys[i])
		}
	}
}

func TestFindDifferences(t *testing.T) {
	// Create test files
	input1 := []string{"apple", "banana", "cherry", "date"}
	input2 := []string{"banana", "date", "fig", "grape"}

	err := os.WriteFile("test_input1.txt", []byte(strings.Join(input1, "\n")+"\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_input1.txt")

	err = os.WriteFile("test_input2.txt", []byte(strings.Join(input2, "\n")+"\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_input2.txt")

	// Run function
	_, err = findDifferences("test_input1.txt", "test_input2.txt",
		"test_output1.txt", "test_output2.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_output1.txt")
	defer os.Remove("test_output2.txt")

	// Verify results
	output1, err := os.ReadFile("test_output1.txt")
	if err != nil {
		t.Fatal(err)
	}
	expected1 := []string{"apple", "cherry"}
	got1 := strings.Split(strings.TrimSpace(string(output1)), "\n")

	if !stringSliceEqual(got1, expected1) {
		t.Errorf("Output1 mismatch:\nExpected: %v\nGot: %v", expected1, got1)
	}

	output2, err := os.ReadFile("test_output2.txt")
	if err != nil {
		t.Fatal(err)
	}
	expected2 := []string{"fig", "grape"}
	got2 := strings.Split(strings.TrimSpace(string(output2)), "\n")

	if !stringSliceEqual(got2, expected2) {
		t.Errorf("Output2 mismatch:\nExpected: %v\nGot: %v", expected2, got2)
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
