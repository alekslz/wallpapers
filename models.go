package main

import (
	"encoding/json"
	"os"
	"sync"
)

// ImageFile represents a wallpaper image with metadata
type ImageFile struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Liked    int    `json:"liked"`
	Disliked int    `json:"disliked"`
}

// DataStore handles JSON file operations
type DataStore struct {
	FilePath string
	mu       sync.RWMutex
}

// NewDataStore creates a new data store
func NewDataStore(filepath string) *DataStore {
	return &DataStore{
		FilePath: filepath,
	}
}

// LoadImages loads all images from JSON file
func (ds *DataStore) LoadImages() ([]ImageFile, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	data, err := os.ReadFile(ds.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []ImageFile{}, nil
		}
		return nil, err
	}

	var images []ImageFile
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// SaveImages saves all images to JSON file
func (ds *DataStore) SaveImages(images []ImageFile) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ds.FilePath, data, 0644)
}

// GetImagesPaginated returns a paginated list of images
func (ds *DataStore) GetImagesPaginated(offset, limit int) ([]ImageFile, int, error) {
	images, err := ds.LoadImages()
	if err != nil {
		return nil, 0, err
	}

	// Filter images with disliked <= 10
	filtered := make([]ImageFile, 0)
	for _, img := range images {
		if img.Disliked <= 10 {
			filtered = append(filtered, img)
		}
	}

	total := len(filtered)

	// Apply pagination
	if offset >= total {
		return []ImageFile{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

// UpdateLike increments the like counter for an image
func (ds *DataStore) UpdateLike(imageName string) error {
	images, err := ds.LoadImages()
	if err != nil {
		return err
	}

	for i := range images {
		if images[i].Name == imageName {
			images[i].Liked++
			break
		}
	}

	return ds.SaveImages(images)
}

// UpdateDislike increments the dislike counter for an image
func (ds *DataStore) UpdateDislike(imageName string) error {
	images, err := ds.LoadImages()
	if err != nil {
		return err
	}

	for i := range images {
		if images[i].Name == imageName {
			images[i].Disliked++
			break
		}
	}

	return ds.SaveImages(images)
}
