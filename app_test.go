package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadAndSaveIni verifies that the LoadIni and SaveIni methods operate
// correctly.  It writes a sample INI file to a temporary directory, loads it
// into the App, modifies one of the values, writes the modified data to a
// new file and asserts that the changes were persisted.  It also tests
// that CurrentFile and GetData return the expected values after loading.
func TestLoadAndSaveIni(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.ini")
	outputPath := filepath.Join(dir, "output.ini")
	content := "[Section1]\nKey1=Value1\nKey2=Value2\n\n[Section2]\nAnotherKey=AnotherValue\n"
	if err := os.WriteFile(inputPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp ini file: %v", err)
	}
	// Create a new App and load the sample INI file
	app := NewApp("")
	data, err := app.LoadIni(inputPath)
	if err != nil {
		t.Fatalf("LoadIni returned error: %v", err)
	}
	// Verify that CurrentFile is updated
	if app.CurrentFile() != inputPath {
		t.Fatalf("CurrentFile expected %s, got %s", inputPath, app.CurrentFile())
	}
	// Verify that GetData returns the same data map
	returned := app.GetData()
	if len(returned) != len(data) {
		t.Fatalf("GetData returned incorrect number of sections: %d", len(returned))
	}
	if v := returned["Section1"]["Key1"]; v != "Value1" {
		t.Fatalf("unexpected value for Section1.Key1: %s", v)
	}
	// Modify a value and save to a new file
	data["Section1"]["Key1"] = "Modified"
	if err := app.SaveIni(outputPath, data); err != nil {
		t.Fatalf("SaveIni returned error: %v", err)
	}
	// Read the saved file and assert that the modifications were written
	outBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read saved ini file: %v", err)
	}
	outContent := string(outBytes)
	if !strings.Contains(outContent, "Modified") {
		t.Fatalf("saved file does not contain modified value: %s", outContent)
	}
}
