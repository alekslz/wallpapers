package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InitData scans the images directory and creates a JSON file with all image metadata
func InitData() {
	imagesDir := "resource/images"
	dataFile := "data/images.json"

	// Check if images directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		log.Printf("Images directory %s does not exist. Creating empty data file.\n", imagesDir)
		createEmptyDataFile(dataFile)
		return
	}

	// Scan images directory
	files, err := os.ReadDir(imagesDir)
	if err != nil {
		log.Fatalf("Error reading images directory: %v", err)
	}

	var images []ImageFile
	id := 1

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Filter for image files
		ext := filepath.Ext(file.Name())
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			images = append(images, ImageFile{
				ID:       id,
				Name:     file.Name(),
				Liked:    0,
				Disliked: 0,
			})
			id++
		}
	}

	// Save to JSON file
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		log.Fatalf("Error writing data file: %v", err)
	}

	fmt.Printf("Successfully initialized data file with %d images\n", len(images))
}

func createEmptyDataFile(path string) {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	if err := os.WriteFile(path, []byte("[]"), 0644); err != nil {
		log.Fatalf("Error creating empty data file: %v", err)
	}

	fmt.Println("Created empty data file")
}

func main() {
	InitData()
}
