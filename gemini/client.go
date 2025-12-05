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
	ProductURL string `json:"product_url"`
	ImageURL   string `json:"image_url"`
	ShopName   string `json:"shop_name"`
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
	// 注: 出力は必ず"JSONのみ"とし、説明文や注釈を一切含めないでください。
	//       product_url,image_url,shop_name をキーとして必ず含め、存在しない場合は null を返してください。
	prompt := `You are a web scraping assistant equipped with Google Search.
Analyze the provided HTML content.
If the HTML is incomplete or the image URL is missing, USE GOOGLE SEARCH to find the official Booth product page for this item.

Requirements:
- RETURN ONLY valid JSON and NOTHING ELSE (no explanations, no surrounding text, no code fences).
- JSON must contain keys: "product_url", "image_url", "shop_name". Use null for missing values.
- When identifying the product page, PRIORITIZE canonical Booth item URLs that include the "/items/" path and the numeric ID (examples: https://booth.pm/ja/items/5949310 or https://USERNAME.booth.pm/items/5949310).
- If multiple candidate item URLs are found, choose the one that clearly corresponds to the product (title, snippet, or meta tags). If uncertain, return the single best candidate.
- If no item URL is found but a shop page exists, set "product_url" to null and set "shop_name" to the shop display name or username.
- If a product page exists but is private/unavailable, still attempt to return the canonical item URL if it can be identified from search results or page metadata.
- Prefer the shop display name (author name shown on the page, or meta author) for "shop_name". If not available, use the username from the product URL subdomain.

Task:
1. Identify the Product URL (format: https://USERNAME.booth.pm/items/XXXXXXX)
2. Identify the Main product image URL.
3. Identify the shop/author username or display name that publishes the product.

Return the result as JSON with this exact structure:
{
	"product_url": "the product URL or null",
	"image_url": "the main product image URL or null",
	"shop_name": "the shop or author display name or username or null"
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

	// Try direct JSON parse first
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		fmt.Printf("Initial JSON parse error: %v\n", err)

		// Attempt: extract the first {...} JSON substring from response
		start := strings.Index(jsonStr, "{")
		end := strings.LastIndex(jsonStr, "}")
		if start != -1 && end != -1 && end > start {
			candidate := jsonStr[start : end+1]
			candidate = strings.TrimSpace(candidate)
			if err2 := json.Unmarshal([]byte(candidate), &result); err2 == nil {
				fmt.Printf("Recovered JSON from substring\n")
			} else {
				fmt.Printf("Failed to parse candidate JSON substring: %v\n", err2)
				// Fallback: try to extract URL/shop/image by regex from free text
				// product URL
				productURL := ""
				imageURL := ""
				shopName := ""

				// find booth item URL
				reURL := regexp.MustCompile(`https?://[^\s'"<>]*booth\.pm[^\s'"<>]*items[^\s'"<>]*`)
				if m := reURL.FindString(jsonStr); m != "" {
					productURL = m
				}

				// find image url
				reImg := regexp.MustCompile(`https?://booth\.pximg\.net/[a-zA-Z0-9/._-]+\.(?:jpg|jpeg|png|webp)`)
				if m := reImg.FindString(jsonStr); m != "" {
					imageURL = m
				}

				// try to find a CJK shop/display name in the text
				reCJK := regexp.MustCompile(`([\p{Han}\p{Hiragana}\p{Katakana}ー々]{2,})`)
				if sm := reCJK.FindString(jsonStr); sm != "" {
					shopName = sm
				}

				// fallback: infer username from productURL subdomain
				if shopName == "" && productURL != "" {
					if u, err3 := url.Parse(productURL); err3 == nil {
						host := u.Hostname()
						parts := strings.Split(host, ".")
						if len(parts) >= 3 {
							shopName = parts[0]
						}
					}
				}

				result = BoothProductInfo{
					ProductURL: productURL,
					ImageURL:   imageURL,
					ShopName:   shopName,
				}
			}
		} else {
			// No JSON substring; apply regex extraction directly
			productURL := ""
			imageURL := ""
			shopName := ""

			reURL := regexp.MustCompile(`https?://[^\s'"<>]*booth\.pm[^\s'"<>]*items[^\s'"<>]*`)
			if m := reURL.FindString(jsonStr); m != "" {
				productURL = m
			}

			reImg := regexp.MustCompile(`https?://booth\.pximg\.net/[a-zA-Z0-9/._-]+\.(?:jpg|jpeg|png|webp)`)
			if m := reImg.FindString(jsonStr); m != "" {
				imageURL = m
			}

			reCJK := regexp.MustCompile(`([\p{Han}\p{Hiragana}\p{Katakana}ー々]{2,})`)
			if sm := reCJK.FindString(jsonStr); sm != "" {
				shopName = sm
			}

			if shopName == "" && productURL != "" {
				if u, err3 := url.Parse(productURL); err3 == nil {
					host := u.Hostname()
					parts := strings.Split(host, ".")
					if len(parts) >= 3 {
						shopName = parts[0]
					}
				}
			}

			result = BoothProductInfo{
				ProductURL: productURL,
				ImageURL:   imageURL,
				ShopName:   shopName,
			}
		}
	}

	// Ensure shop name fallback from product url if still empty
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
