package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"google.golang.org/genai"
)

// BoothProductInfo represents extracted Booth product information
type BoothProductInfo struct {
	ProductURL string `json:"productUrl"`
	ImageURL   string `json:"imageUrl"`
	ShopName   string `json:"shopName"`
}

// ExtractBoothProductInfo uses Gemini API with Google Search Grounding
func ExtractBoothProductInfo(html string, apiKey string) (*BoothProductInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second) // 検索時間を考慮し少し長めに
	defer cancel()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Improved Prompt
	prompt := `You are a specialized data extraction backend.
Your task is to identify a Booth.pm product from the provided HTML/Text snippet and return structured data.

Process:
1. Analyze the input text/HTML.
2. If the Product URL or Image URL is missing or incomplete, USE GOOGLE SEARCH to find the official Booth product page.
   - Search Query Hint: Look for product titles, folder names, or unique identifiers in the input.
3. Prioritize canonical URLs: https://booth.pm/ja/items/ID or https://subdomain.booth.pm/items/ID.

Constraints:
- OUTPUT FORMAT: Raw JSON only. NO Markdown code blocks, no explanations.
- MISSING VALUES: Use an empty string "" for any missing fields. DO NOT use null.
- JSON KEYS: "productUrl", "imageUrl", "shopName".

Field Definitions:
- productUrl: The canonical URL of the item page.
- imageUrl: The URL of the main product image (usually booth.pximg.net).
- shopName: The display name of the shop/author. If unknown, infer from the subdomain.

Input Data:
` + html

	fmt.Println("Sending request to Gemini API (google.golang.org/genai) with Grounding...")

	tools := []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}

	temperature := float32(0.1)
	config := &genai.GenerateContentConfig{
		Temperature: &temperature,
		Tools:       tools,
		// SystemInstruction: 検索ツールを使う場合、SystemInstructionよりもPromptに直接書いたほうが指示が通りやすい傾向があります
	}

	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text result
	jsonStr := resp.Text()

	// デバッグ用: 生レスポンスの確認（開発中のみ有効にすることを推奨）
	// fmt.Printf("Gemini raw response: %s\n", jsonStr)

	// Clean up Markdown code blocks if strictly generated
	jsonStr = strings.TrimSpace(jsonStr)
	if strings.HasPrefix(jsonStr, "```") {
		// Remove first line (```json or ```)
		if i := strings.Index(jsonStr, "\n"); i != -1 {
			jsonStr = jsonStr[i+1:]
		}
		// Remove last line (```)
		if strings.HasSuffix(jsonStr, "```") {
			jsonStr = strings.TrimSuffix(jsonStr, "```")
		}
		jsonStr = strings.TrimSpace(jsonStr)
	}

	var result BoothProductInfo

	// Attempt Standard Parsing
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Printf("Initial JSON parse error: %v. Attempting heuristic extraction.\n", err)
		// Parse失敗時のフォールバック処理（正規表現抽出）は維持することを推奨します
		// 元のコードにある正規表現ロジックをここに配置してください。
		// プロンプト改善により、ここに来る確率は減りますが、Groundingの出力は予測しきれないため安全策として残します。

		return heuristicFallback(jsonStr) // ※リファクタリングとして別関数への切り出しを推奨
	}

	// Post-processing check
	// ShopNameが空でProductURLがある場合の推論ロジックは維持
	if result.ShopName == "" && result.ProductURL != "" {
		if u, err := url.Parse(result.ProductURL); err == nil {
			host := u.Hostname()
			parts := strings.Split(host, ".")
			if len(parts) >= 3 {
				result.ShopName = parts[0]
			}
		}
	}

	return &result, nil
}

// heuristicFallback はJSONパース失敗時の救済措置です（元のコードのロジックを移動）
func heuristicFallback(text string) (*BoothProductInfo, error) {
	result := &BoothProductInfo{}

	// Product URL extraction
	reURL := regexp.MustCompile(`https?://[^\s'"<>]*booth\.pm[^\s'"<>]*items[^\s'"<>]*`)
	if m := reURL.FindString(text); m != "" {
		result.ProductURL = m
	}

	// Image URL extraction
	reImg := regexp.MustCompile(`https?://booth\.pximg\.net/[a-zA-Z0-9/._-]+\.(?:jpg|jpeg|png|webp)`)
	if m := reImg.FindString(text); m != "" {
		result.ImageURL = m
	}

	// Shop Name extraction (simplified CJK regex)
	reCJK := regexp.MustCompile(`([\p{Han}\p{Hiragana}\p{Katakana}ー々]{2,})`)
	if sm := reCJK.FindString(text); sm != "" {
		result.ShopName = sm
	}

	// Fallback shopname logic
	if result.ShopName == "" && result.ProductURL != "" {
		if u, err := url.Parse(result.ProductURL); err == nil {
			host := u.Hostname()
			parts := strings.Split(host, ".")
			if len(parts) >= 3 {
				result.ShopName = parts[0]
			}
		}
	}

	return result, nil
}
