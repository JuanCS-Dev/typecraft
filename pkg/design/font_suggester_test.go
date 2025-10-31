package design

import (
	"testing"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFontSuggester(t *testing.T) {
	suggester, err := NewFontSuggester()
	require.NoError(t, err)
	require.NotNil(t, suggester)
	require.NotNil(t, suggester.db)
}

func TestFontSuggester_Suggest(t *testing.T) {
	suggester, err := NewFontSuggester()
	require.NoError(t, err)

	t.Run("fiction content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre: "fiction",
			Tone: analyzer.ToneProfile{
				Formal:   0.5,
				Academic: 0.3,
			},
			Complexity: 0.5,
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		assert.Len(t, pairs, 3)
		
		// Todos devem ter body e heading
		for _, pair := range pairs {
			assert.NotEmpty(t, pair.Body)
			assert.NotEmpty(t, pair.Heading)
			assert.NotEmpty(t, pair.Mood)
			assert.NotEmpty(t, pair.Rationale)
			assert.Greater(t, pair.Score, 0.0)
		}
	})

	t.Run("mystery content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:   "mystery",
			SentimentScore: -0.3,
			Tone: analyzer.ToneProfile{
				Formal: 0.6,
			},
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		assert.Len(t, pairs, 3)
		
		// Primeiro par deve ter score mais alto
		assert.GreaterOrEqual(t, pairs[0].Score, pairs[1].Score)
		assert.GreaterOrEqual(t, pairs[1].Score, pairs[2].Score)
	})

	t.Run("technical content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:     "technical",
			TechnicalDensity: 0.7,
			Tone: analyzer.ToneProfile{
				Technical: 0.8,
				Formal:    0.7,
			},
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		assert.NotEmpty(t, pairs)
		
		// Technical content deve retornar boas sugestões
		for _, pair := range pairs {
			assert.NotEmpty(t, pair.Body)
			assert.NotEmpty(t, pair.Heading)
		}
	})

	t.Run("academic content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre: "academic",
			Complexity:   0.8,
			Tone: analyzer.ToneProfile{
				Academic: 0.9,
				Formal:   0.8,
			},
			Formality: 0.85,
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		assert.NotEmpty(t, pairs)
	})

	t.Run("romance content with high sentiment", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:   "romance",
			SentimentScore: 0.5,
			Tone: analyzer.ToneProfile{
				Creative: 0.7,
			},
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		assert.NotEmpty(t, pairs)
	})

	t.Run("multi-genre content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:   "fiction",
			SecondaryGenre: "mystery",
			Tone: analyzer.ToneProfile{
				Formal: 0.5,
			},
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		
		// Deve ter mais opções quando há gênero secundário
		assert.GreaterOrEqual(t, len(pairs), 3)
	})

	t.Run("nil analysis", func(t *testing.T) {
		_, err := suggester.Suggest(nil)
		assert.Error(t, err)
	})

	t.Run("unknown genre defaults to fiction", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre: "unknown_genre",
			Tone: analyzer.ToneProfile{
				Formal: 0.5,
			},
		}

		pairs, err := suggester.Suggest(analysis)
		require.NoError(t, err)
		assert.NotEmpty(t, pairs)
	})
}

func TestFontSuggester_GetFontMetadata(t *testing.T) {
	suggester, err := NewFontSuggester()
	require.NoError(t, err)

	t.Run("existing font", func(t *testing.T) {
		meta, ok := suggester.GetFontMetadata("Crimson Text")
		assert.True(t, ok)
		assert.NotNil(t, meta)
		assert.Equal(t, "Crimson Text", meta.Name)
		assert.True(t, meta.SupportsLatin)
	})

	t.Run("non-existing font", func(t *testing.T) {
		meta, ok := suggester.GetFontMetadata("NonExistentFont")
		assert.False(t, ok)
		assert.Nil(t, meta)
	})
}

func TestFontSuggester_ListAvailableFonts(t *testing.T) {
	suggester, err := NewFontSuggester()
	require.NoError(t, err)

	t.Run("serif_body category", func(t *testing.T) {
		fonts := suggester.ListAvailableFonts("serif_body")
		assert.NotEmpty(t, fonts)
		assert.Contains(t, fonts, "Crimson Text")
	})

	t.Run("sans_body category", func(t *testing.T) {
		fonts := suggester.ListAvailableFonts("sans_body")
		assert.NotEmpty(t, fonts)
	})

	t.Run("monospace category", func(t *testing.T) {
		fonts := suggester.ListAvailableFonts("monospace")
		assert.NotEmpty(t, fonts)
		assert.Contains(t, fonts, "Fira Code")
	})

	t.Run("non-existing category", func(t *testing.T) {
		fonts := suggester.ListAvailableFonts("non_existing")
		assert.Empty(t, fonts)
	})
}

func TestFontSuggester_CalculateScore(t *testing.T) {
	suggester, err := NewFontSuggester()
	require.NoError(t, err)

	tests := []struct {
		name        string
		pair        FontPair
		analysis    *analyzer.ContentAnalysis
		expectBoost bool
	}{
		{
			name: "classic mood with high formality",
			pair: FontPair{Mood: "classic"},
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{Formal: 0.7},
			},
			expectBoost: true,
		},
		{
			name: "modern mood with low formality",
			pair: FontPair{Mood: "modern"},
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{Formal: 0.3},
			},
			expectBoost: true,
		},
		{
			name: "technical mood with high tech density",
			pair: FontPair{Mood: "technical"},
			analysis: &analyzer.ContentAnalysis{
				TechnicalDensity: 0.6,
			},
			expectBoost: true,
		},
		{
			name: "academic mood with high academic tone",
			pair: FontPair{Mood: "academic"},
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{Academic: 0.7},
			},
			expectBoost: true,
		},
		{
			name: "romantic mood with positive sentiment",
			pair: FontPair{Mood: "romantic"},
			analysis: &analyzer.ContentAnalysis{
				SentimentScore: 0.4,
			},
			expectBoost: true,
		},
		{
			name: "dark mood with negative sentiment",
			pair: FontPair{Mood: "dark"},
			analysis: &analyzer.ContentAnalysis{
				SentimentScore: -0.3,
			},
			expectBoost: true,
		},
		{
			name: "monospace with technical content",
			pair: FontPair{Monospace: "Fira Code"},
			analysis: &analyzer.ContentAnalysis{
				TechnicalDensity: 0.5,
			},
			expectBoost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := suggester.calculateScore(&tt.pair, tt.analysis)
			
			if tt.expectBoost {
				assert.Greater(t, score, 1.0, "Expected score boost")
			}
			
			assert.Greater(t, score, 0.0)
		})
	}
}

func BenchmarkFontSuggester_Suggest(b *testing.B) {
	suggester, err := NewFontSuggester()
	if err != nil {
		b.Fatal(err)
	}

	analysis := &analyzer.ContentAnalysis{
		PrimaryGenre:   "fiction",
		SecondaryGenre: "mystery",
		Complexity:     0.6,
		SentimentScore: -0.2,
		Tone: analyzer.ToneProfile{
			Formal:    0.5,
			Academic:  0.3,
			Technical: 0.2,
		},
		TechnicalDensity: 0.3,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := suggester.Suggest(analysis)
		if err != nil {
			b.Fatal(err)
		}
	}
}
