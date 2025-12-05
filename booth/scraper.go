package booth

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// BoothInfo represents information scraped from Booth
type BoothInfo struct {
	ProductURL string
	ImageURL   string
	ShopName   string
}

// ExtractSearchQuery extracts a clean search query from folder name
// Removes brackets, version numbers, and special characters
func ExtractSearchQuery(folderName string) string {
	// Remove content within brackets: 【】[] () （）
	re := regexp.MustCompile(`[\[\]【】\(\)（）]`)
	cleaned := re.ReplaceAllString(folderName, " ")

	// Remove version patterns like v1.0, ver1.0, _v1, etc.
	versionRe := regexp.MustCompile(`(?i)(_?v(er)?\.?\d+(\.\d+)*)`)
	cleaned = versionRe.ReplaceAllString(cleaned, " ")

	// Remove common keywords that don't help search
	keywords := []string{"booth", "gumroad", "avatar", "アバター", "unity", "vrc", "vrchat"}
	for _, keyword := range keywords {
		keywordRe := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		cleaned = keywordRe.ReplaceAllString(cleaned, " ")
	}

	// Remove underscores and multiple spaces
	cleaned = strings.ReplaceAll(cleaned, "_", " ")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

	// Trim whitespace
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

// SearchBooth searches Booth for a product and returns the first result's URL and thumbnail
func SearchBooth(query string) (*BoothInfo, error) {
	// Extract clean search query from folder name
	cleanQuery := ExtractSearchQuery(query)
	if cleanQuery == "" {
		cleanQuery = query // Fallback to original if cleaning removed everything
	}

	// URL encode the query and construct search URL
	searchURL := fmt.Sprintf("https://booth.pm/ja/search/%s", url.QueryEscape(cleanQuery))

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request with proper headers
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	html := string(body)

	// Extract first product URL
	// Booth product URLs follow the pattern: https://booth.pm/ja/items/XXXXXXX or https://USERNAME.booth.pm/items/XXXXXXX
	productURLPattern := regexp.MustCompile(`https://[a-zA-Z0-9_-]+\.booth\.pm/items/\d+`)
	productURLs := productURLPattern.FindAllString(html, -1)

	if len(productURLs) == 0 {
		return nil, fmt.Errorf("no products found for query: %s", query)
	}

	productURL := productURLs[0]

	// Try to infer shop name from product URL subdomain (e.g. username.booth.pm)
	var shopName string
	if u, err := url.Parse(productURL); err == nil {
		host := u.Hostname()
		// host might be like username.booth.pm
		parts := strings.Split(host, ".")
		if len(parts) >= 3 {
			shopName = parts[0]
		}
	}

	// Extract thumbnail image URL
	// Look for image URL near the product URL
	// Booth uses data-src or src attributes for images
	imageURLPattern := regexp.MustCompile(`https://booth\.pximg\.net/[a-zA-Z0-9/._-]+\.(?:jpg|jpeg|png|webp)`)
	imageURLs := imageURLPattern.FindAllString(html, -1)

	var imageURL string
	if len(imageURLs) > 0 {
		// Get the first thumbnail (usually the main product image)
		imageURL = imageURLs[0]
		// Remove size suffix and get original or larger size
		imageURL = regexp.MustCompile(`/c/\d+x\d+_[a-z0-9]+/`).ReplaceAllString(imageURL, "/")
	}

	return &BoothInfo{
		ProductURL: productURL,
		ImageURL:   imageURL,
		ShopName:   shopName,
	}, nil
}

// ExtractBoothProductInfo extracts product info from a Booth product page URL
func ExtractBoothProductInfo(productURL string) (*BoothInfo, error) {
	// Validate URL
	if !strings.Contains(productURL, "booth.pm/items/") {
		return nil, fmt.Errorf("invalid Booth product URL")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", productURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	html := string(body)

	// Extract main product image from market-item-detail-item-image class
	var imageURL string
	var shopName string

	// Look for the market-item-detail-item-image class and extract image URL
	// Pattern: find div/element with class containing 'market-item-detail-item-image'
	// and then find img tag with src or data-src attribute
	classPattern := regexp.MustCompile(`class="[^"]*market-item-detail-item-image[^"]*"[^>]*>[\s\S]*?<img[^>]+(?:data-src|src)="([^"]+)"`)
	matches := classPattern.FindStringSubmatch(html)

	if len(matches) > 1 {
		imageURL = matches[1]
	}
	// Try to extract shop name from meta author or canonical info
	authorPattern := regexp.MustCompile(`<meta name="author" content="([^"]+)"`)
	aMatches := authorPattern.FindStringSubmatch(html)
	if len(aMatches) > 1 {
		shopName = strings.TrimSpace(aMatches[1])
	}
	if shopName == "" {
		// Fallback: parse host from productURL
		if u, err := url.Parse(productURL); err == nil {
			host := u.Hostname()
			parts := strings.Split(host, ".")
			if len(parts) >= 3 {
				shopName = parts[0]
			}
		}
	} else {
		// Fallback: Look for img tag inside market-item-detail-item-image more loosely
		fallbackPattern := regexp.MustCompile(`market-item-detail-item-image[\s\S]{0,500}?<img[^>]+(?:data-src|src)="([^"]+)"`)
		fallbackMatches := fallbackPattern.FindStringSubmatch(html)
		if len(fallbackMatches) > 1 {
			imageURL = fallbackMatches[1]
		} else {
			// Second fallback: try og:image meta tag
			ogPattern := regexp.MustCompile(`<meta property="og:image" content="([^"]+)"`)
			ogMatches := ogPattern.FindStringSubmatch(html)
			if len(ogMatches) > 1 {
				imageURL = ogMatches[1]
			}
		}
	}

	return &BoothInfo{
		ProductURL: productURL,
		ImageURL:   imageURL,
		ShopName:   shopName,
	}, nil
}
