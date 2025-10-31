package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentAnalyzer_Analyze(t *testing.T) {
	analyzer := NewContentAnalyzer()

	t.Run("basic fiction content", func(t *testing.T) {
		content := `
# The Mystery of the Lost Key

Chapter 1: The Beginning

Detective Sarah walked into the dark room. The murder scene was exactly as she expected.
Every clue pointed to the same suspect. The investigation would be complex.

She examined the evidence carefully. The victim's story was tragic.
Her character had always been mysterious, full of secrets.
`

		analysis, err := analyzer.Analyze(content)
		require.NoError(t, err)
		require.NotNil(t, analysis)

		assert.Greater(t, analysis.WordCount, 0)
		assert.Greater(t, analysis.SentenceCount, 0)
		assert.NotEmpty(t, analysis.PrimaryGenre)
		
		// Deve detectar mystery como gênero principal
		assert.Contains(t, []string{"mystery", "fiction"}, analysis.PrimaryGenre)
		
		// Deve ter alguma complexidade calculada
		assert.GreaterOrEqual(t, analysis.Complexity, 0.0)
		assert.LessOrEqual(t, analysis.Complexity, 1.0)
	})

	t.Run("technical content", func(t *testing.T) {
		content := `
# System Architecture

## Introduction

The algorithm implements a data processing system using a client-server architecture.
The framework provides an interface for database operations.
Implementation details include network protocol analysis and security measures.

The methodology involves systematic analysis of the process.
Therefore, we conclude that the system architecture is optimal.
`

		analysis, err := analyzer.Analyze(content)
		require.NoError(t, err)

		// Deve detectar como técnico ou acadêmico
		assert.Contains(t, []string{"technical", "academic"}, analysis.PrimaryGenre)
		
		// Deve ter alta densidade técnica
		assert.Greater(t, analysis.TechnicalDensity, 0.1)
		
		// Deve ter alta formalidade
		assert.Greater(t, analysis.Formality, 0.3)
	})

	t.Run("content with images and equations", func(t *testing.T) {
		content := `
# Math Tutorial

The equation for a line is:

$$y = mx + b$$

Here's a diagram:

![Graph](graph.png)

Another equation:

$E = mc^2$
`

		analysis, err := analyzer.Analyze(content)
		require.NoError(t, err)

		assert.Equal(t, 1, analysis.ImageCount)
		assert.Equal(t, 2, analysis.EquationCount)
		assert.True(t, analysis.HasComplexMath() == false) // < 10 equations
	})

	t.Run("romance content with sentiment", func(t *testing.T) {
		content := `
# A Love Story

Chapter 1

Their love was wonderful and beautiful. Her heart filled with joy.
The passion between them was amazing and perfect.
Their relationship was fantastic, a brilliant romance that made them both happy.
`

		analysis, err := analyzer.Analyze(content)
		require.NoError(t, err)

		// Deve detectar romance
		assert.Contains(t, []string{"romance", "fiction"}, analysis.PrimaryGenre)
		
		// Deve ter sentimento positivo
		assert.Greater(t, analysis.SentimentScore, 0.0)
	})

	t.Run("empty content", func(t *testing.T) {
		content := ""

		analysis, err := analyzer.Analyze(content)
		require.NoError(t, err)
		
		assert.Equal(t, 0, analysis.WordCount)
		assert.Equal(t, 0, analysis.SentenceCount)
	})
}

func TestContentAnalysis_IsAcademic(t *testing.T) {
	tests := []struct {
		name     string
		analysis ContentAnalysis
		want     bool
	}{
		{
			name: "high academic tone",
			analysis: ContentAnalysis{
				Tone: ToneProfile{Academic: 0.7},
			},
			want: true,
		},
		{
			name: "high formality and technical",
			analysis: ContentAnalysis{
				Formality:        0.8,
				TechnicalDensity: 0.6,
				Tone:             ToneProfile{Academic: 0.5},
			},
			want: true,
		},
		{
			name: "not academic",
			analysis: ContentAnalysis{
				Formality:        0.3,
				TechnicalDensity: 0.2,
				Tone:             ToneProfile{Academic: 0.3},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.analysis.IsAcademic()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContentAnalysis_HasComplexMath(t *testing.T) {
	tests := []struct {
		name          string
		equationCount int
		want          bool
	}{
		{"no equations", 0, false},
		{"few equations", 5, false},
		{"many equations", 15, true},
		{"exactly 10", 10, false},
		{"exactly 11", 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &ContentAnalysis{
				EquationCount: tt.equationCount,
			}
			got := analysis.HasComplexMath()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContentAnalysis_HasRichMedia(t *testing.T) {
	tests := []struct {
		name       string
		imageRatio float64
		want       bool
	}{
		{"no images", 0.0, false},
		{"few images", 0.05, false},
		{"many images", 0.15, true},
		{"exactly threshold", 0.1, false},
		{"just above threshold", 0.11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &ContentAnalysis{
				ImageRatio: tt.imageRatio,
			}
			got := analysis.HasRichMedia()
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkContentAnalyzer_Analyze(b *testing.B) {
	analyzer := NewContentAnalyzer()
	
	// Generate 50k words content
	content := generateLargeContent(50000)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := analyzer.Analyze(content)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func generateLargeContent(wordCount int) string {
	words := []string{
		"the", "detective", "investigated", "crime", "scene",
		"evidence", "suspect", "character", "story", "plot",
		"mystery", "chapter", "analysis", "system", "data",
	}
	
	result := ""
	for i := 0; i < wordCount; i++ {
		result += words[i%len(words)] + " "
		if i%15 == 14 {
			result += ". "
		}
		if i%100 == 99 {
			result += "\n\n"
		}
	}
	
	return result
}
