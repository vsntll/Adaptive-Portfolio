package main

import (
	"encoding/json"
	"fmt"
	"os"
	"resume-backend/parser"
)
func main() {
	fmt.Println("Starting Resume Parser...")

	// 1. Parse the PDF
	data, err := parser.ParsePDF("parser/resume.pdf")
	if err != nil {
		fmt.Printf("Error parsing PDF: %v\n", err)
		return
	}

	// 2. Convert to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	// 3. Save to Frontend Data Directory
	// Ensure this path matches your actual structure
	outputPath := "../frontend/src/data/info.json"
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("✅ Successfully generated info.json!")
	fmt.Printf("Name detected: %s\n", data.Name)
}