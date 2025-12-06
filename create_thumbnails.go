package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const (
	sourceDir    = "resource/images"
	thumbnailDir = "resource/imagesmall"
	maxWidth     = 300
	maxHeight    = 300
)

func main() {
	log.Println("Creating thumbnails...")

	// Create thumbnail directory if it doesn't exist
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		log.Fatalf("Error creating thumbnail directory: %v", err)
	}

	// Get all images from source directory
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Fatalf("Error reading source directory: %v", err)
	}

	created := 0
	skipped := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			continue
		}

		sourcePath := filepath.Join(sourceDir, file.Name())
		thumbPath := filepath.Join(thumbnailDir, file.Name())

		// Skip if thumbnail already exists
		if _, err := os.Stat(thumbPath); err == nil {
			skipped++
			continue
		}

		log.Printf("Creating thumbnail: %s\n", file.Name())

		if err := createThumbnail(sourcePath, thumbPath); err != nil {
			log.Printf("Error creating thumbnail for %s: %v\n", file.Name(), err)
			continue
		}

		created++
	}

	log.Printf("Created: %d thumbnails, Skipped: %d (already exist)\n", created, skipped)
}

// createThumbnail creates a thumbnail from source image
func createThumbnail(sourcePath, destPath string) error {
	// Open source image
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Decode image
	img, format, err := image.Decode(sourceFile)
	if err != nil {
		return err
	}

	// Calculate thumbnail dimensions maintaining aspect ratio
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	ratio := float64(width) / float64(height)
	var thumbWidth, thumbHeight int

	if width > height {
		thumbWidth = maxWidth
		thumbHeight = int(float64(maxWidth) / ratio)
	} else {
		thumbHeight = maxHeight
		thumbWidth = int(float64(maxHeight) * ratio)
	}

	// Create thumbnail image
	thumbnail := image.NewRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))

	// Resize using high-quality algorithm
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, bounds, draw.Over, nil)

	// Save thumbnail
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Encode based on original format
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(destFile, thumbnail, &jpeg.Options{Quality: 85})
	case "png":
		return png.Encode(destFile, thumbnail)
	case "gif":
		return gif.Encode(destFile, thumbnail, nil)
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}
