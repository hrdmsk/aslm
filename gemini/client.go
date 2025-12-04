package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

// BoothProductInfo represents extracted Booth product information
type BoothProductInfo struct {
	ProductURL string `json:"product_url"`
	ImageURL   string `json:"image_url"`
}

// ExtractBoothProductInfo uses Gemini API with Google Search Grounding
func ExtractBoothProductInfo(html string, apiKey string) (*BoothProductInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize client with API key
	// ドキュメント準拠: https://docs.cloud.google.com/vertex-ai/generative-ai/docs/start/quickstart?hl=ja&usertype=apikey#go-gen-ai-sdk
	// APIキーを使用する場合は ClientConfig にセットします。
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Prompt: HTMLが渡された場合はそれを解析し、不足があれば検索で補うように指示
	prompt := `You are a web scraping assistant equipped with Google Search.
Analyze the provided HTML content.
If the HTML is incomplete or the image URL is missing, USE GOOGLE SEARCH to find the official Booth product page for this item.

Task:
1. Identify the Product URL (format: https://USERNAME.booth.pm/items/XXXXXXX)
2. Identify the Main product image URL.

Return the result as JSON with this exact structure:
{
  "product_url": "the product URL",
  "image_url": "the main product image URL"
}

HTML Snippet:
` + html

	fmt.Println("Sending request to Gemini API (google.golang.org/genai) with Grounding...")

	// Google Search Grounding ツールの定義
	tools := []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}

	// Configure generation options
	temperature := float32(0.1)
	config := &genai.GenerateContentConfig{
		// 【修正】JSONモードとTools(Grounding)は併用できないため、MIMEType指定を削除
		// ResponseMIMEType: "application/json",
		Temperature: &temperature,
		Tools:       tools,
	}

	// 【指定修正】Gemini 2.5 Flash を使用
	// ユーザー指定: 1.5は使用せず、2.5Flashを強制的に使用する
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), config)
	if err != nil {
		// エラー発生時の詳細ログ
		fmt.Printf("Gemini API error for model %s: %v\n", modelName, err)
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text result
	jsonStr := resp.Text()
	fmt.Printf("Gemini raw response: %s\n", jsonStr)

	// JSON cleanup logic
	jsonStr = strings.TrimSpace(jsonStr)
	if strings.HasPrefix(jsonStr, "```") {
		if strings.HasPrefix(jsonStr, "```json") {
			jsonStr = strings.TrimPrefix(jsonStr, "```json")
		} else {
			jsonStr = strings.TrimPrefix(jsonStr, "```")
		}
		// Unconditional TrimSuffix as per linter suggestion
		jsonStr = strings.TrimSuffix(jsonStr, "```")

		jsonStr = strings.TrimSpace(jsonStr)
	}

	var result BoothProductInfo
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("JSON parse error: %v\n", err)
		return nil, fmt.Errorf("failed to parse Gemini response: %w, raw: %s", err, jsonStr)
	}

	return &result, nil
}