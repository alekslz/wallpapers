package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	baseURL    = "https://boards.4chan.org/wg/"
	imagesDir  = "resource/images"
	maxPages   = 10
	maxThreads = 15
)

func main() {
	log.Println("Starting 4chan /wg/ wallpaper scraper...")

	// Create images directory if it doesn't exist
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		log.Fatalf("Error creating images directory: %v", err)
	}

	// Scrape images
	imageURLs := scrapeImageURLs()
	log.Printf("Found %d image URLs to download\n", len(imageURLs))

	// Download images
	downloadImages(imageURLs)

	log.Println("Scraping complete!")
}

// scrapeImageURLs scrapes image URLs from 4chan /wg/ board
func scrapeImageURLs() []string {
	var imageURLs []string
	threadURLs := scrapeThreadURLs()

	log.Printf("Found %d threads to scrape\n", len(threadURLs))

	for i, threadURL := range threadURLs {
		if i >= maxThreads {
			break
		}

		log.Printf("Scraping thread %d/%d: %s\n", i+1, len(threadURLs), threadURL)
		urls := scrapeImagesFromThread(threadURL)
		imageURLs = append(imageURLs, urls...)

		// Be nice to 4chan's servers
		time.Sleep(1 * time.Second)
	}

	return imageURLs
}

// scrapeThreadURLs gets all thread URLs from the first N pages of /wg/
func scrapeThreadURLs() []string {
	var threadURLs []string
	threadRegex := regexp.MustCompile(`href="(/wg/thread/\d+)"`)

	pages := []string{""}
	for i := 2; i <= maxPages; i++ {
		pages = append(pages, fmt.Sprintf("%d", i))
	}

	for _, page := range pages {
		url := baseURL + page
		log.Printf("Fetching page: %s\n", url)

		html, err := fetchURL(url)
		if err != nil {
			log.Printf("Error fetching page %s: %v\n", url, err)
			continue
		}

		matches := threadRegex.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 1 {
				threadURL := "https://boards.4chan.org" + match[1]
				threadURLs = append(threadURLs, threadURL)
			}
		}

		// Be nice to 4chan's servers
		time.Sleep(1 * time.Second)
	}

	return threadURLs
}

// scrapeImagesFromThread extracts image URLs from a thread
func scrapeImagesFromThread(threadURL string) []string {
	var imageURLs []string

	html, err := fetchURL(threadURL)
	if err != nil {
		log.Printf("Error fetching thread %s: %v\n", threadURL, err)
		return imageURLs
	}

	// Match image URLs from 4cdn.org
	// Pattern matches: //i.4cdn.org/wg/1234567890.jpg or .png
	imageRegex := regexp.MustCompile(`//i\.4cdn\.org/wg/\d+\.(jpg|png|gif)`)
	matches := imageRegex.FindAllString(html, -1)

	for _, match := range matches {
		imageURL := "https:" + match
		imageURLs = append(imageURLs, imageURL)
	}

	return imageURLs
}

// fetchURL fetches content from a URL
func fetchURL(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// downloadImages downloads images to the images directory
func downloadImages(imageURLs []string) {
	downloaded := 0
	skipped := 0

	for i, imageURL := range imageURLs {
		// Extract filename from URL
		parts := strings.Split(imageURL, "/")
		filename := parts[len(parts)-1]
		filepath := filepath.Join(imagesDir, filename)

		// Skip if file already exists
		if _, err := os.Stat(filepath); err == nil {
			skipped++
			continue
		}

		log.Printf("Downloading %d/%d: %s\n", i+1, len(imageURLs), filename)

		if err := downloadFile(imageURL, filepath); err != nil {
			log.Printf("Error downloading %s: %v\n", imageURL, err)
			continue
		}

		downloaded++

		// Be nice to 4chan's servers
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Downloaded: %d, Skipped (already exist): %d\n", downloaded, skipped)
}

// downloadFile downloads a file from a URL and saves it to disk
func downloadFile(url, filepath string) error {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
