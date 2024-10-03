package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Store shortened URLs in memory
var urlStore = make(map[string]string)

// Characters for generating random short codes
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano()) // Initialize random seed
}

// Function to generate a random string for the short URL
func generateShortURL(length int) string {
	var shortURL strings.Builder
	for i := 0; i < length; i++ {
		randomChar := charset[rand.Intn(len(charset))]
		shortURL.WriteByte(randomChar)
	}
	return shortURL.String()
}

// Handler to shorten a URL
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL from the request query parameters (e.g., /shorten?url=http://example.com)
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Generate a short code
	shortURL := generateShortURL(6)

	// Store the short URL and original URL in the map
	urlStore[shortURL] = originalURL

	// Respond with the shortened URL
	shortenedLink := fmt.Sprintf("http://localhost:8080/%s", shortURL)
	fmt.Fprintf(w, "Shortened URL: %s", shortenedLink)
}

// Handler to redirect to the original URL
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the short code from the URL path (e.g., /abc123)
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	// Lookup the original URL in the map
	originalURL, exists := urlStore[shortURL]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	// Route to handle shortening URLs
	http.HandleFunc("/shorten", shortenURLHandler)

	// Route to handle redirecting short URLs
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Server running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
