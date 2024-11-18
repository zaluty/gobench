package main

import (
	"fmt"
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	fmt.Println("Running test...")
	// Add your test assertions here
	if 1+1 != 2 {
		t.Error("Basic math failed")
	}
}

func TestAnotherExample(t *testing.T) {
	fmt.Println("Running another test...")
	// Add more test assertions
	if len("hello") != 5 {
		t.Error("String length check failed")
	}
}

func TestMain(m *testing.M) {
	fmt.Println("Starting tests...")
	// Run any setup code here

	code := m.Run()

	// Run any cleanup code here
	fmt.Println("All tests completed!")

	os.Exit(code)
}
