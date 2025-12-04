package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/genai"
)

// BoothProductInfo represents extracted Booth product information
type BoothProductInfo struct {
	ProductURL string `json:"product_url"`
	ImageURL   string `json:"image_url"`
}

// ExtractBoothProductInfo uses Gemini API to extract product information from HTML
func ExtractBoothProductInfo(html string, apiKey string) (*BoothProductInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize client with API key
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	// Note: genai.Client might not have a Close method, checking documentation or sample code implies it's not needed or handled differently.
	// If it has Close, we should defer it. Assuming it doesn't based on user sample.

	prompt := `You are a web scraping assistant. Extract the following information from this Booth product page HTML:

1. Product URL (the canonical URL of the product, usually in format: https://USERNAME.booth.pm/items/XXXXXXX)
2. Main product image URL (the primary/main product thumbnail image)

Return the result as JSON with this exact structure:
{
  "product_url": "the product URL",
  "image_url": "the main product image URL"
}

HTML content:
` + html

	fmt.Println("Sending request to Gemini API (google.golang.org/genai)...")

	// Configure generation options
	temperature := float32(0.1)
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		Temperature:      &temperature,
	}

	// Use gemini-1.5-flash-001 as it is a specific version
	modelName := "gemini-1.5-flash-001"

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), config)
	if err != nil {
		fmt.Printf("Gemini API error: %v\n", err)
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text result
	jsonStr := resp.Text()
	fmt.Printf("Gemini raw response: %s\n", jsonStr)

	// Clean up markdown code blocks if present
	if len(jsonStr) > 7 && jsonStr[:7] == "```json" {
		jsonStr = jsonStr[7:]
		if len(jsonStr) > 3 && jsonStr[len(jsonStr)-3:] == "```" {
			jsonStr = jsonStr[:len(jsonStr)-3]
		}
	}

	var result BoothProductInfo
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("JSON parse error: %v\n", err)
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	return &result, nil
}
