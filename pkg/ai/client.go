package ai

// Este pacote re-exporta os tipos do internal/ai para uso em pkg/
// Seguindo a convenção de Go onde internal/ não deve ser acessado diretamente

import (
	"github.com/JuanCS-Dev/typecraft/internal/ai"
)

// Client é um alias para internal/ai.Client
type Client = ai.Client

// DesignSystem é um alias para internal/ai.DesignSystem
type DesignSystem = ai.DesignSystem

// NewClient cria um novo cliente AI
func NewClient(apiKey, model string, maxTokens int, temperature float32) *Client {
	return ai.NewClient(apiKey, model, maxTokens, temperature)
}
