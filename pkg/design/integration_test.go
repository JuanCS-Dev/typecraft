package design

import (
	"testing"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndToEndDesignSystem testa o fluxo completo
func TestEndToEndDesignSystem(t *testing.T) {
	// Setup
	contentAnalyzer := analyzer.NewContentAnalyzer()
	fontSuggester, err := NewFontSuggester()
	require.NoError(t, err)
	colorGenerator := NewColorGenerator()

	t.Run("complete design flow for mystery novel", func(t *testing.T) {
		// 1. Manuscrito de teste
		manuscript := `
# The Midnight Detective

Chapter 1: The Case Begins

Detective Sarah Harrison stood in the rain, staring at the locked door. 
The murder scene was exactly as she had expected - dark, ominous, and full of clues.
Every piece of evidence pointed to the same suspect, but the investigation would be complex.

She examined the crime scene carefully. The victim's story was tragic, 
a tale of fear and terrible secrets. Her character had always been mysterious, 
full of hidden motives that would take a skilled detective to uncover.

The clue at the door frame suggested the murderer had planned everything meticulously.
`

		// 2. Analisar conteúdo
		analysis, err := contentAnalyzer.Analyze(manuscript)
		require.NoError(t, err)
		
		t.Logf("✅ Content Analysis:")
		t.Logf("   Genre: %s (secondary: %s)", analysis.PrimaryGenre, analysis.SecondaryGenre)
		t.Logf("   Complexity: %.2f", analysis.Complexity)
		t.Logf("   Sentiment: %.2f", analysis.SentimentScore)
		t.Logf("   Word Count: %d", analysis.WordCount)
		t.Logf("   Academic Tone: %.2f", analysis.Tone.Academic)

		// Verificar detecção de gênero
		assert.Contains(t, []string{"mystery", "fiction"}, analysis.PrimaryGenre)
		
		// 3. Sugerir fontes
		fontPairs, err := fontSuggester.Suggest(analysis)
		require.NoError(t, err)
		require.NotEmpty(t, fontPairs)
		
		t.Logf("\n✅ Font Suggestions:")
		for i, pair := range fontPairs {
			t.Logf("   %d. %s + %s (%.2f)", i+1, pair.Body, pair.Heading, pair.Score)
			t.Logf("      Mood: %s | %s", pair.Mood, pair.Rationale)
		}
		
		// Verificar qualidade das sugestões
		assert.GreaterOrEqual(t, fontPairs[0].Score, fontPairs[1].Score)
		
		// 4. Gerar paleta de cores
		palette, err := colorGenerator.Generate(analysis)
		require.NoError(t, err)
		
		t.Logf("\n✅ Color Palette:")
		t.Logf("   Primary: %s", palette.Primary.Hex)
		t.Logf("   Secondary: %s", palette.Secondary.Hex)
		t.Logf("   Accent: %s", palette.Accent.Hex)
		t.Logf("   Background: %s", palette.Background.Hex)
		t.Logf("   Text: %s", palette.Text.Hex)
		t.Logf("   Mood: %s | Harmony: %s", palette.Metadata.Mood, palette.Metadata.Harmony)
		t.Logf("   Contrast: %.2f (%s)", palette.Metadata.Contrast, palette.Metadata.Accessibility)
		t.Logf("   Rationale: %s", palette.Metadata.Rationale)
		
		// Verificar acessibilidade
		assert.GreaterOrEqual(t, palette.Metadata.Contrast, 4.5)
		assert.Contains(t, []string{"AAA", "AA"}, palette.Metadata.Accessibility)
		
		// Mood deve ser dark para mystery com sentimento negativo
		if analysis.SentimentScore < -0.2 {
			assert.Equal(t, "dark", palette.Metadata.Mood)
		}
	})

	t.Run("complete design flow for technical book", func(t *testing.T) {
		manuscript := `
# System Architecture Design

## Introduction

The algorithm implements a data processing system using a client-server architecture.
The framework provides an interface for database operations and network protocol analysis.

Implementation details include security measures and systematic process optimization.
The methodology involves comprehensive analysis of system performance and scalability.

Therefore, we conclude that the architecture meets all technical requirements.
The function provides optimal resource utilization through efficient algorithm design.
`

		// Análise
		analysis, err := contentAnalyzer.Analyze(manuscript)
		require.NoError(t, err)
		
		t.Logf("✅ Technical Content Analysis:")
		t.Logf("   Genre: %s", analysis.PrimaryGenre)
		t.Logf("   Technical Density: %.2f", analysis.TechnicalDensity)
		t.Logf("   Formality: %.2f", analysis.Formality)
		
		// Deve detectar como technical ou academic
		assert.Contains(t, []string{"technical", "academic"}, analysis.PrimaryGenre)
		assert.Greater(t, analysis.TechnicalDensity, 0.2)
		
		// Fontes
		fontPairs, err := fontSuggester.Suggest(analysis)
		require.NoError(t, err)
		
		t.Logf("\n✅ Technical Font Suggestions:")
		t.Logf("   Best: %s + %s", fontPairs[0].Body, fontPairs[0].Heading)
		
		// Cores
		palette, err := colorGenerator.Generate(analysis)
		require.NoError(t, err)
		
		t.Logf("\n✅ Technical Color Palette:")
		t.Logf("   Harmony: %s", palette.Metadata.Harmony)
		t.Logf("   Accessibility: %s", palette.Metadata.Accessibility)
		
		// Technical/Academic deve ter boa acessibilidade
		assert.Contains(t, []string{"AAA", "AA"}, palette.Metadata.Accessibility)
		assert.Contains(t, []string{"complementary", "analogous"}, palette.Metadata.Harmony)
	})

	t.Run("complete design flow for romance", func(t *testing.T) {
		manuscript := `
# A Love Beyond Time

Chapter 1

Their love was wonderful and beautiful, a perfect romance that transcended everything.
Her heart filled with joy and passion as their relationship blossomed.
The embrace was amazing, a brilliant moment of pure happiness and desire.

Their love story was fantastic, full of tender moments and romantic gestures.
Every kiss was perfect, every glance filled with love and longing.
`

		analysis, err := contentAnalyzer.Analyze(manuscript)
		require.NoError(t, err)
		
		t.Logf("✅ Romance Content Analysis:")
		t.Logf("   Genre: %s", analysis.PrimaryGenre)
		t.Logf("   Sentiment: %.2f (positive)", analysis.SentimentScore)
		
		// Deve detectar romance
		assert.Contains(t, []string{"romance", "fiction"}, analysis.PrimaryGenre)
		
		// Deve ter sentimento positivo
		assert.Greater(t, analysis.SentimentScore, 0.3)
		
		// Fontes
		fontPairs, err := fontSuggester.Suggest(analysis)
		require.NoError(t, err)
		
		t.Logf("\n✅ Romance Font Suggestions:")
		t.Logf("   Best: %s + %s (%s)", fontPairs[0].Body, fontPairs[0].Heading, fontPairs[0].Mood)
		
		// Cores
		palette, err := colorGenerator.Generate(analysis)
		require.NoError(t, err)
		
		t.Logf("\n✅ Romance Color Palette:")
		t.Logf("   Mood: %s", palette.Metadata.Mood)
		t.Logf("   Colors: %s, %s, %s", palette.Primary.Hex, palette.Secondary.Hex, palette.Accent.Hex)
		
		// Deve ter mood bright ou vibrant
		assert.Contains(t, []string{"bright", "vibrant"}, palette.Metadata.Mood)
	})
}

// BenchmarkCompleteDesignFlow benchmark do fluxo completo
func BenchmarkCompleteDesignFlow(b *testing.B) {
	contentAnalyzer := analyzer.NewContentAnalyzer()
	fontSuggester, _ := NewFontSuggester()
	colorGenerator := NewColorGenerator()

	manuscript := `
# The Detective's Last Case

Detective Morgan examined the crime scene carefully. Every clue pointed
to a complex investigation ahead. The victim's story was mysterious,
full of secrets that would take skill to uncover. The murder weapon
lay hidden, waiting to reveal its dark tale.
`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Complete flow
		analysis, _ := contentAnalyzer.Analyze(manuscript)
		_, _ = fontSuggester.Suggest(analysis)
		_, _ = colorGenerator.Generate(analysis)
	}
}
