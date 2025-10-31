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

// EnhanceTypography melhora a tipografia do texto usando IA
func (c *Client) EnhanceTypography(text string) (string, error) {
	ctx := context.Background()
	
	prompt := fmt.Sprintf(`Melhore a tipografia do seguinte texto, aplicando:
- Aspas tipográficas corretas ("" e '')
- Travessões em vez de hífens quando apropriado
- Espaçamento adequado
- Reticências corretas (…)

Retorne APENAS o texto melhorado, sem explicações.

Texto:
%s`, text)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       c.model,
			MaxTokens:   c.maxTokens,
			Temperature: 0.3,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return text, err
	}

	if len(resp.Choices) == 0 {
		return text, fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// DesignSystem representa um sistema de design gerado
type DesignSystem struct {
	CSS         string
	Description string
	Colors      map[string]string
	Fonts       []string
}

// GenerateDesignSystem gera um sistema de design baseado em prompt
func (c *Client) GenerateDesignSystem(prompt string) (*DesignSystem, error) {
	ctx := context.Background()
	
	systemPrompt := `Você é um designer editorial especializado em tipografia e design de livros.
Gere um sistema de design em CSS que seja elegante, legível e apropriado para impressão.
Retorne apenas o CSS válido, sem explicações ou markdown.`

	fullPrompt := fmt.Sprintf(`Crie um design system em CSS para um livro com as seguintes características:
%s

Inclua:
- Variáveis CSS para cores e tipografia
- Estilos para @page (paginação)
- Tipografia responsiva e elegante
- Estilos para capítulos, citações, listas
- Detalhes refinados (drop caps, ornamentos)

Formato A5 (148mm x 210mm), margens adequadas para impressão.`, prompt)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       c.model,
			MaxTokens:   2000,
			Temperature: 0.7,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fullPrompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate design system: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	css := resp.Choices[0].Message.Content
	
	// Remove markdown code blocks se presentes
	css = strings.TrimPrefix(css, "```css")
	css = strings.TrimPrefix(css, "```")
	css = strings.TrimSuffix(css, "```")
	css = strings.TrimSpace(css)

	return &DesignSystem{
		CSS:         css,
		Description: prompt,
	}, nil
}
