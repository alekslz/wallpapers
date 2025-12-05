package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var dataStore *DataStore

const (
	itemsPerGroup = 70
	dataFilePath  = "data/images.json"
)

func main() {
	// Initialize data store
	dataStore = NewDataStore(dataFilePath)

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal(err)
	}

	// Ensure download directory exists
	if err := os.MkdirAll("download", 0755); err != nil {
		log.Fatal(err)
	}

	// Set up routes
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/count", handleCount)
	http.HandleFunc("/api/images", handleImages)
	http.HandleFunc("/api/like", handleLike)
	http.HandleFunc("/api/dislike", handleDislike)
	http.HandleFunc("/api/download", handleDownload)

	// Serve static files
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("resource"))))
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("download"))))

	// Start server
	port := "8080"
	log.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// serveIndex serves the main HTML page
func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

// handleCount returns the total number of groups for pagination
func handleCount(w http.ResponseWriter, r *http.Request) {
	images, err := dataStore.LoadImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter images with disliked <= 10
	filtered := 0
	for _, img := range images {
		if img.Disliked <= 10 {
			filtered++
		}
	}

	totalGroups := int(math.Ceil(float64(filtered) / float64(itemsPerGroup)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalGroups)
}

// handleImages returns paginated images
func handleImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groupNoStr := r.FormValue("group_no")
	groupNo, err := strconv.Atoi(groupNoStr)
	if err != nil || groupNo < 0 {
		http.Error(w, "Invalid group number", http.StatusBadRequest)
		return
	}

	offset := groupNo * itemsPerGroup
	images, _, err := dataStore.GetImagesPaginated(offset, itemsPerGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build HTML response (matching PHP output)
	var html strings.Builder
	for _, img := range images {
		html.WriteString(fmt.Sprintf(
			`<div class='div-img'><img class='img' src='resource/imagesmall/%s'></img><div class='checkboxes'><label for='checkbox'><input type='checkbox'>P.F.D</label><a class='like'>Like</a><a class='dislike'>Dislike</a></div></div>`,
			img.Name,
		))
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html.String()))
}

// handleLike increments the like counter for an image
func handleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	imagePath := r.FormValue("like")
	if imagePath == "" {
		http.Error(w, "Missing image path", http.StatusBadRequest)
		return
	}

	// Extract filename from path
	imageName := strings.TrimPrefix(imagePath, "resource/images/")

	if err := dataStore.UpdateLike(imageName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleDislike increments the dislike counter for an image
func handleDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	imagePath := r.FormValue("dislike")
	if imagePath == "" {
		http.Error(w, "Missing image path", http.StatusBadRequest)
		return
	}

	// Extract filename from path
	imageName := strings.TrimPrefix(imagePath, "resource/images/")

	if err := dataStore.UpdateDislike(imageName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleDownload creates a ZIP file with selected images
func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	objStr := r.FormValue("obj")
	if objStr == "" {
		http.Error(w, "Missing obj parameter", http.StatusBadRequest)
		return
	}

	// Parse the JSON object
	var images map[string]struct {
		Src string `json:"src"`
	}
	if err := json.Unmarshal([]byte(objStr), &images); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create ZIP file
	rand.Seed(time.Now().UnixNano())
	zipName := fmt.Sprintf("%d.zip", rand.Int())
	zipPath := filepath.Join("download", zipName)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	// Add images to ZIP
	for _, img := range images {
		srcPath := img.Src
		fileName := strings.TrimPrefix(srcPath, "resource/images/")

		// Open source file
		sourceFile, err := os.Open(srcPath)
		if err != nil {
			log.Printf("Error opening file %s: %v", srcPath, err)
			continue
		}

		// Create file in ZIP
		zipEntry, err := zipWriter.Create(fileName)
		if err != nil {
			sourceFile.Close()
			log.Printf("Error creating zip entry for %s: %v", fileName, err)
			continue
		}

		// Copy file content
		if _, err := io.Copy(zipEntry, sourceFile); err != nil {
			sourceFile.Close()
			log.Printf("Error copying file %s: %v", fileName, err)
			continue
		}

		sourceFile.Close()
	}

	zipWriter.Close()

	// Read the ZIP file
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send ZIP file to client
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipName))
	w.Header().Set("Content-Length", strconv.Itoa(len(zipData)))
	w.Write(zipData)

	// Clean up ZIP file after sending (optional, you might want to keep them)
	go func() {
		time.Sleep(5 * time.Second)
		os.Remove(zipPath)
	}()
}
