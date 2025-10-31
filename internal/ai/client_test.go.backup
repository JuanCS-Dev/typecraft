package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", "gpt-4", 2000, 0.3)
	
	assert.NotNil(t, client)
	assert.NotNil(t, client.client)
	assert.Equal(t, "gpt-4", client.model)
	assert.Equal(t, 2000, client.maxTokens)
	assert.Equal(t, float32(0.3), client.temperature)
}

func TestClient_AnalyzeText_EmptyText(t *testing.T) {
	client := NewClient("test-key", "gpt-4", 1000, 0.3)
	
	_, err := client.AnalyzeText(nil, "")
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// Nota: Testes de integração real com OpenAI estão em test/integration/
