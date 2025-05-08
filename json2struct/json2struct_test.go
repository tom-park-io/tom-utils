package json2struct

import (
	"fmt"
	"os"
	"testing"
)

// go test -v --run TestConvertJSONToStructor
func TestConvertJSONToStructor(t *testing.T) {

	inputFile := "data/test.json"

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create generator with configuration
	config := Config{
		StructPrefix: "Generated", // Default prefix for all structs
		IndentSize:   4,           // Standard Go indentation
	}
	generator := NewGenerator(config)

	// Generate structs
	if err := generator.GenerateFromJSON(data); err != nil {
		fmt.Printf("Error generating structs: %v\n", err)
		os.Exit(1)
	}

	// Write to file with random name
	filename, err := generator.WriteToFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated structs and saved to %s\n", filename)

}
