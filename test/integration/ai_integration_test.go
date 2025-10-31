package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/ai"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIClient_RealAnalysis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Carrega .env
	_ = godotenv.Load("../../.env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	client := ai.NewClient(apiKey, "gpt-4o", 2000, 0.3)

	tests := []struct {
		name     string
		text     string
		validate func(*testing.T, string)
	}{
		{
			name: "fiction analysis",
			text: `
				Era uma vez, em um reino distante, uma princesa chamada Isabella.
				Ela vivia em um castelo alto, cercado por florestas encantadas.
				
				— Hoje é o dia — disse ela, olhando pela janela.
				
				O príncipe chegaria ao entardecer. Seus olhos azuis brilhavam
				de expectativa enquanto ela escolhia o vestido perfeito para o baile.
				
				Mas algo inesperado estava prestes a acontecer...
			`,
			validate: func(t *testing.T, jsonResp string) {
				analysis, err := ai.ParseAnalysisResponse(jsonResp, "test-fiction-123")
				require.NoError(t, err)
				require.NotNil(t, analysis)
				
				// Deve detectar ficção
				assert.Contains(t, []string{"fiction", "narrative", "story"}, analysis.Genre)
				
				// Tom deve ser narrativo/conversational
				assert.NotEmpty(t, analysis.Tone.Primary)
				
				// Pipeline simples ou standard (não LaTeX para ficção)
				assert.Contains(t, []string{"simple", "standard", "html"}, analysis.RecommendedPipeline)
				
				// Não deve ter muita matemática
				assert.False(t, analysis.HasMath)
				
				t.Logf("✅ Fiction Analysis:")
				t.Logf("   Genre: %s (%.2f)", analysis.Genre, analysis.GenreConfidence)
				t.Logf("   Tone: %s (formality: %.2f)", analysis.Tone.Primary, analysis.Tone.Formality)
				t.Logf("   Pipeline: %s", analysis.RecommendedPipeline)
				t.Logf("   Special: Math=%v Code=%v Images=%v", analysis.HasMath, analysis.HasCode, analysis.HasImages)
			},
		},
		{
			name: "technical analysis",
			text: `
				# Introdução aos Algoritmos de Ordenação
				
				## Bubble Sort
				
				O algoritmo Bubble Sort possui complexidade O(n²) no pior caso.
				
				` + "```python" + `
				def bubble_sort(arr):
				    n = len(arr)
				    for i in range(n):
				        for j in range(0, n-i-1):
				            if arr[j] > arr[j+1]:
				                arr[j], arr[j+1] = arr[j+1], arr[j]
				` + "```" + `
				
				A análise de complexidade pode ser expressa como:
				T(n) = n(n-1)/2 comparações
				
				| Algoritmo | Melhor Caso | Pior Caso | Espaço |
				|-----------|-------------|-----------|---------|
				| Bubble    | O(n)        | O(n²)     | O(1)    |
				| Quick     | O(n log n)  | O(n²)     | O(log n)|
			`,
			validate: func(t *testing.T, jsonResp string) {
				analysis, err := ai.ParseAnalysisResponse(jsonResp, "test-tech-456")
				require.NoError(t, err)
				require.NotNil(t, analysis)
				
				// Deve detectar conteúdo técnico
				assert.Contains(t, []string{"technical", "computer_science", "programming"}, analysis.Genre)
				
				// Deve detectar código
				assert.True(t, analysis.HasCode)
				
				// Densidade técnica alta
				assert.Greater(t, analysis.Complexity.TechnicalDensity, 0.0)
				
				// Pipeline pode ser standard ou latex
				assert.Contains(t, []string{"standard", "latex"}, analysis.RecommendedPipeline)
				
				t.Logf("✅ Technical Analysis:")
				t.Logf("   Genre: %s (%.2f)", analysis.Genre, analysis.GenreConfidence)
				t.Logf("   Code: %v", analysis.HasCode)
				t.Logf("   Math: %v", analysis.HasMath)
				t.Logf("   Pipeline: %s", analysis.RecommendedPipeline)
				t.Logf("   Should use LaTeX: %v", analysis.ShouldUseLaTeX())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			jsonResp, err := client.AnalyzeText(ctx, tt.text)
			require.NoError(t, err)
			require.NotEmpty(t, jsonResp)

			t.Logf("Raw AI Response:\n%s", jsonResp)

			tt.validate(t, jsonResp)
		})
	}
}
