package main

import (
	"aslm/booth"
	"aslm/config"
	"aslm/db"
	"aslm/gemini"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// FileItem represents a file or directory
type FileItem struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"` // "file" or "folder"
	Path     string   `json:"path"`
	Source   string   `json:"source"`   // "booth", "gumroad", etc.
	Url      string   `json:"url"`      // Link to product page
	ImageUrl string   `json:"imageUrl"` // Thumbnail URL
	Tags     []string `json:"tags"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := db.InitDB(); err != nil {
		fmt.Printf("Error initializing DB: %v\n", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// UpdateProduct updates the product information
func (a *App) UpdateProduct(path string, url string, imageUrl string, tags []string) error {
	return db.UpdateProduct(path, url, imageUrl, tags)
}

// ListFiles returns a list of files and directories in the given path
func (a *App) ListFiles(path string) ([]FileItem, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []FileItem
	for _, entry := range entries {
		itemType := "file"
		if entry.IsDir() {
			itemType = "folder"
		}

		fullPath := filepath.Join(path, entry.Name())

		// Register product if not exists (auto-discovery)
		// In a real app, maybe we only do this for folders or specific files
		if itemType == "folder" {
			_ = db.RegisterProduct(fullPath, entry.Name())
		}

		// Get info from DB
		info, _ := db.GetProductInfo(fullPath)

		item := FileItem{
			Name: entry.Name(),
			Type: itemType,
			Path: fullPath,
		}

		if info != nil {
			item.Url = info.Url
			item.ImageUrl = info.ImageUrl
			item.Tags = info.Tags
			// Source detection logic could go here based on URL or tags
			if info.Url != "" {
				// Simple heuristic for demo
				if strings.Contains(info.Url, "booth.pm") {
					item.Source = "booth"
				} else if strings.Contains(info.Url, "gumroad.com") {
					item.Source = "gumroad"
				}
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// GetParentProduct returns the nearest parent product info for a given path
func (a *App) GetParentProduct(path string) (*db.ProductInfo, error) {
	return db.GetParentProductInfo(path)
}

// BoothInfo represents information from Booth
type BoothInfo struct {
	ProductURL string `json:"productUrl"`
	ImageURL   string `json:"imageUrl"`
}

// AutoFetchBoothInfo automatically fetches product info from Booth based on folder name
func (a *App) AutoFetchBoothInfo(folderName string) (*BoothInfo, error) {
	// Import booth package at the top of the file is required
	info, err := booth.SearchBooth(folderName)
	if err != nil {
		return nil, err
	}

	return &BoothInfo{
		ProductURL: info.ProductURL,
		ImageURL:   info.ImageURL,
	}, nil
}

// FetchBoothImageFromURL fetches thumbnail image from a Booth product URL
func (a *App) FetchBoothImageFromURL(productURL string) (string, error) {
	info, err := booth.ExtractBoothProductInfo(productURL)
	if err != nil {
		return "", err
	}

	return info.ImageURL, nil
}

// GetGeminiApiKey retrieves the saved Gemini API key
func (a *App) GetGeminiApiKey() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	return cfg.GeminiAPIKey, nil
}

// SaveGeminiApiKey saves the Gemini API key
func (a *App) SaveGeminiApiKey(apiKey string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		// Create default config if doesn't exist
		cfg = &config.Config{
			HomePath:     "D:/VRChatAssetPack",
			GeminiAPIKey: apiKey,
		}
	} else {
		cfg.GeminiAPIKey = apiKey
	}
	return config.SaveConfig(cfg)
}

// FetchBoothInfoWithGemini uses Gemini API to fetch Booth product info
func (a *App) FetchBoothInfoWithGemini(folderName string) (*BoothInfo, error) {
	// Get API key
	apiKey, err := a.GetGeminiApiKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	// Extract clean search query
	cleanQuery := booth.ExtractSearchQuery(folderName)
	if cleanQuery == "" {
		cleanQuery = folderName
	}

	// Fetch Booth search page HTML
	searchURL := fmt.Sprintf("https://booth.pm/ja/search/%s", cleanQuery)
	fmt.Printf("Searching Booth: %s\n", searchURL)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch search page: %v\n", err)
		return nil, fmt.Errorf("failed to fetch search page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Fetched HTML size: %d bytes\n", len(body))

	// Use Gemini to extract product info
	info, err := gemini.ExtractBoothProductInfo(string(body), apiKey)
	if err != nil {
		fmt.Printf("Failed to extract with Gemini: %v\n", err)
		return nil, fmt.Errorf("failed to extract with Gemini: %w", err)
	}

	fmt.Printf("Gemini extracted info: %+v\n", info)

	return &BoothInfo{
		ProductURL: info.ProductURL,
		ImageURL:   info.ImageURL,
	}, nil
}
