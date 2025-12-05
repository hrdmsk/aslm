package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func InitDB() error {
	// ユーザーのホームディレクトリを取得して、そこにDBファイルを保存する
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(homeDir, ".aslm")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(appDir, "aslm.db")

	var openErr error
	DB, openErr = sql.Open("sqlite", dbPath)
	if openErr != nil {
		return openErr
	}

	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE NOT NULL,
		name TEXT,
		url TEXT,
		image_url TEXT,
		shop_name TEXT
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS product_tags (
		product_id INTEGER,
		tag_id INTEGER,
		PRIMARY KEY (product_id, tag_id),
		FOREIGN KEY (product_id) REFERENCES products(id),
		FOREIGN KEY (tag_id) REFERENCES tags(id)
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("Error creating tables: %v", err)
		return err
	}

	// Migration: Add image_url and shop_name columns if they don't exist (simple check)
	// In a real app, use a migration tool. Here we just try to add them and ignore error if exists
	DB.Exec(`ALTER TABLE products ADD COLUMN image_url TEXT;`) // Ignore error if column exists
	DB.Exec(`ALTER TABLE products ADD COLUMN shop_name TEXT;`) // Ignore error if column exists

	return nil
}

// RegisterProduct registers a product if it doesn't exist
func RegisterProduct(path string, name string) error {
	// Prevent registering a product inside an existing product folder
	parentInfo, err := GetParentProductInfo(path)
	if err != nil {
		return err
	}
	if parentInfo != nil {
		return fmt.Errorf("cannot register product inside existing product: parent=%s", parentInfo.Path)
	}

	query := `INSERT OR IGNORE INTO products (path, name) VALUES (?, ?)`
	_, err = DB.Exec(query, path, name)
	return err
}

// UpdateProduct updates the product information
func UpdateProduct(path string, url string, imageUrl string, shopName string, tags []string) error {
	// 1. Update product details
	query := `UPDATE products SET url = ?, image_url = ?, shop_name = ? WHERE path = ?`
	if _, err := DB.Exec(query, url, imageUrl, shopName, path); err != nil {
		return err
	}

	// 2. Handle Tags
	// First, get product ID
	var productId int
	if err := DB.QueryRow("SELECT id FROM products WHERE path = ?", path).Scan(&productId); err != nil {
		return err
	}

	// Clear existing tags for this product
	if _, err := DB.Exec("DELETE FROM product_tags WHERE product_id = ?", productId); err != nil {
		return err
	}

	// Add new tags
	for _, tagName := range tags {
		if tagName == "" {
			continue
		}
		if err := AddTag(path, tagName); err != nil {
			log.Printf("Error adding tag %s: %v", tagName, err)
			// Continue adding other tags even if one fails
		}
	}

	return nil
}

// AddTag adds a tag to a product
func AddTag(path string, tagName string) error {
	// 1. Ensure tag exists
	tagQuery := `INSERT OR IGNORE INTO tags (name) VALUES (?)`
	if _, err := DB.Exec(tagQuery, tagName); err != nil {
		return err
	}

	// 2. Get product ID
	var productId int
	if err := DB.QueryRow("SELECT id FROM products WHERE path = ?", path).Scan(&productId); err != nil {
		return err
	}

	// 3. Get tag ID
	var tagId int
	if err := DB.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagId); err != nil {
		return err
	}

	// 4. Link product and tag
	linkQuery := `INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES (?, ?)`
	_, err := DB.Exec(linkQuery, productId, tagId)
	return err
}

type ProductInfo struct {
	Path     string
	Name     string
	Url      string
	ImageUrl string
	ShopName string
	Tags     []string
}

// GetProductInfo retrieves information for a product
func GetProductInfo(path string) (*ProductInfo, error) {
	var info ProductInfo
	info.Path = path

	// Get basic info
	query := `SELECT name, url, image_url, shop_name FROM products WHERE path = ?`
	var url sql.NullString
	var imageUrl sql.NullString
	var shopName sql.NullString

	err := DB.QueryRow(query, path).Scan(&info.Name, &url, &imageUrl)
	if err == sql.ErrNoRows {
		return nil, nil // Not found is not an error here, just return nil
	} else if err != nil {
		return nil, err
	}

	if url.Valid {
		info.Url = url.String
	}
	if imageUrl.Valid {
		info.ImageUrl = imageUrl.String
	}
	if shopName.Valid {
		info.ShopName = shopName.String
	}

	// Get tags
	tagQuery := `
		SELECT t.name 
		FROM tags t
		JOIN product_tags pt ON t.id = pt.tag_id
		JOIN products p ON p.id = pt.product_id
		WHERE p.path = ?
	`
	rows, err := DB.Query(tagQuery, path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}
		info.Tags = append(info.Tags, tagName)
	}

	return &info, nil
}

// GetParentProductInfo finds the nearest parent product info in the database
func GetParentProductInfo(path string) (*ProductInfo, error) {
	currentPath := filepath.Clean(path)

	for {
		// Move to parent directory
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// Reached root
			return nil, nil
		}

		// Try to get product info for parent
		info, err := GetProductInfo(parentPath)
		if err == sql.ErrNoRows || info == nil {
			// Parent not found, continue searching
			currentPath = parentPath
			continue
		}
		if err != nil {
			return nil, err
		}

		// Found a parent with product info
		return info, nil
	}
}
