package main

import (
	"aslm/db"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
