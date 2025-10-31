package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
}

func NewClient(apiKey, model string, maxTokens int, temperature float32) *Client {
	return &Client{
		client:      openai.NewClient(apiKey),
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}
}

func (c *Client) AnalyzeText(ctx context.Context, text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}

	prompt := buildAnalysisPrompt(text)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       c.model,
			MaxTokens:   c.maxTokens,
			Temperature: c.temperature,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to analyze text: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// ParseAnalysisResponse converte a resposta JSON da IA em estrutura de análise
func ParseAnalysisResponse(jsonResp string, projectID string) (*domain.AIAnalysis, error) {
	// Remove markdown code blocks se presentes
	jsonResp = strings.TrimPrefix(jsonResp, "```json")
	jsonResp = strings.TrimPrefix(jsonResp, "```")
	jsonResp = strings.TrimSuffix(jsonResp, "```")
	jsonResp = strings.TrimSpace(jsonResp)

	var analysis domain.AIAnalysis
	if err := json.Unmarshal([]byte(jsonResp), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	analysis.ProjectID = projectID
	return &analysis, nil
}
