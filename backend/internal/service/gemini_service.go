package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GeminiService struct {
	apiKey     string
	model      string
	baseURL    string
	timeout    time.Duration
	systemPrompt string
}

type GeminiMessage struct {
	Role  string `json:"role"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

type GeminiRequest struct {
	Contents         []GeminiMessage `json:"contents"`
	GenerationConfig struct {
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"maxOutputTokens"`
		TopP        float64 `json:"topP"`
		TopK        int     `json:"topK"`
	} `json:"generationConfig"`
	SafetySettings []struct {
		Category  string `json:"category"`
		Threshold string `json:"threshold"`
	} `json:"safetySettings"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewGeminiService() *GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = "AIzaSyCHUXz2cqq1kiBPzuN2AnazG29w2hvG1BQ" // Default from config
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-pro"
	}

	baseURL := os.Getenv("GEMINI_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta/models/"
	}

	timeout := 30 * time.Second
	if timeoutStr := os.Getenv("GEMINI_TIMEOUT"); timeoutStr != "" {
		if parsed, err := time.ParseDuration(timeoutStr + "s"); err == nil {
			timeout = parsed
		}
	}

	// Load system prompt from file or use default
	systemPrompt := loadSystemPrompt()

	return &GeminiService{
		apiKey:       apiKey,
		model:        model,
		baseURL:      baseURL,
		timeout:      timeout,
		systemPrompt: systemPrompt,
	}
}

func loadSystemPrompt() string {
	// Try to load from file
	paths := []string{
		"/data/manual_entrenamiento.txt",
		"./data/manual_entrenamiento.txt",
		"../data/manual_entrenamiento.txt",
	}

	for _, path := range paths {
		if content, err := os.ReadFile(path); err == nil {
			return string(content)
		}
	}

	// Default prompt
	return "Eres Plot AI, el asistente digital de Plot Center. Ayuda a los usuarios con órdenes de trabajo y consultas del taller. Responde siempre en español de forma amigable y profesional."
}

func (s *GeminiService) SendMessage(userMessage string, history []struct {
	Role    string
	Content string
}) (string, error) {
	url := fmt.Sprintf("%s%s:generateContent?key=%s", s.baseURL, s.model, s.apiKey)

	// Build messages
	messages := []GeminiMessage{}

	// Add system prompt as first message
	if s.systemPrompt != "" {
		messages = append(messages, GeminiMessage{
			Role: "user",
			Parts: []struct {
				Text string `json:"text"`
			}{{Text: s.systemPrompt}},
		})
		messages = append(messages, GeminiMessage{
			Role: "model",
			Parts: []struct {
				Text string `json:"text"`
			}{{Text: "Entendido, estoy listo para ayudar."}},
		})
	}

	// Add conversation history
	for _, msg := range history {
		messages = append(messages, GeminiMessage{
			Role: msg.Role,
			Parts: []struct {
				Text string `json:"text"`
			}{{Text: msg.Content}},
		})
	}

	// Add current user message
	messages = append(messages, GeminiMessage{
		Role: "user",
		Parts: []struct {
			Text string `json:"text"`
		}{{Text: userMessage}},
	})

	// Build request
	req := GeminiRequest{
		Contents: messages,
		GenerationConfig: struct {
			Temperature float64 `json:"temperature"`
			MaxTokens   int     `json:"maxOutputTokens"`
			TopP        float64 `json:"topP"`
			TopK        int     `json:"topK"`
		}{
			Temperature: 0.7,
			MaxTokens:   2000,
			TopP:        0.95,
			TopK:        40,
		},
		SafetySettings: []struct {
			Category  string `json:"category"`
			Threshold string `json:"threshold"`
		}{
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
			{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
			{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
		},
	}

	// Marshal request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{Timeout: s.timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp GeminiResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil && errorResp.Error != nil {
			return "", fmt.Errorf("gemini API error: %s", errorResp.Error.Message)
		}
		return "", fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	// Extract message
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

