package design

import (
	"testing"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewColorGenerator(t *testing.T) {
	generator := NewColorGenerator()
	require.NotNil(t, generator)
	require.NotNil(t, generator.sentimentHueMap)
}

func TestColorGenerator_Generate(t *testing.T) {
	generator := NewColorGenerator()

	t.Run("fiction with positive sentiment", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:   "fiction",
			SentimentScore: 0.5,
			Tone: analyzer.ToneProfile{
				Formal:   0.5,
				Creative: 0.6,
			},
		}

		palette, err := generator.Generate(analysis)
		require.NoError(t, err)
		require.NotNil(t, palette)

		// Verificar que todas as cores foram geradas
		assert.NotEmpty(t, palette.Primary.Hex)
		assert.NotEmpty(t, palette.Secondary.Hex)
		assert.NotEmpty(t, palette.Accent.Hex)
		assert.NotEmpty(t, palette.Background.Hex)
		assert.NotEmpty(t, palette.Text.Hex)

		// Verificar metadata
		assert.NotEmpty(t, palette.Metadata.Mood)
		assert.NotEmpty(t, palette.Metadata.Harmony)
		assert.Greater(t, palette.Metadata.Contrast, 0.0)
		assert.NotEmpty(t, palette.Metadata.Accessibility)
	})

	t.Run("mystery with negative sentiment", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:   "mystery",
			SentimentScore: -0.4,
			Tone: analyzer.ToneProfile{
				Formal: 0.7,
			},
		}

		palette, err := generator.Generate(analysis)
		require.NoError(t, err)

		// Deve gerar paleta escura/misteriosa
		assert.Equal(t, "dark", palette.Metadata.Mood)
	})

	t.Run("romance with high positive sentiment", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:   "romance",
			SentimentScore: 0.6,
			Tone: analyzer.ToneProfile{
				Creative: 0.8,
			},
		}

		palette, err := generator.Generate(analysis)
		require.NoError(t, err)

		// Deve ter mood bright ou vibrant
		assert.Contains(t, []string{"bright", "vibrant"}, palette.Metadata.Mood)
	})

	t.Run("technical content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre:     "technical",
			TechnicalDensity: 0.8,
			Tone: analyzer.ToneProfile{
				Formal:    0.8,
				Technical: 0.9,
			},
		}

		palette, err := generator.Generate(analysis)
		require.NoError(t, err)

		// Deve usar harmonia complementar para contraste
		assert.Equal(t, "complementary", palette.Metadata.Harmony)
		assert.Equal(t, "professional", palette.Metadata.Mood)
	})

	t.Run("academic content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre: "academic",
			Tone: analyzer.ToneProfile{
				Academic: 0.9,
				Formal:   0.9,
			},
			Formality: 0.85,
		}

		palette, err := generator.Generate(analysis)
		require.NoError(t, err)

		// Deve ter cores conservadoras e alto contraste
		assert.GreaterOrEqual(t, palette.Metadata.Contrast, 4.5)
		assert.Contains(t, []string{"AAA", "AA"}, palette.Metadata.Accessibility)
		
		// Background deve ser branco puro para academic
		assert.Equal(t, "#FFFFFF", palette.Background.Hex)
		assert.Equal(t, "#000000", palette.Text.Hex)
	})

	t.Run("creative content", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			PrimaryGenre: "fiction",
			Tone: analyzer.ToneProfile{
				Creative: 0.9,
			},
		}

		palette, err := generator.Generate(analysis)
		require.NoError(t, err)

		// Deve usar triadic para cores vibrantes
		assert.Equal(t, "triadic", palette.Metadata.Harmony)
	})

	t.Run("nil analysis", func(t *testing.T) {
		_, err := generator.Generate(nil)
		assert.Error(t, err)
	})
}

func TestColorGenerator_SelectPrimaryHue(t *testing.T) {
	generator := NewColorGenerator()

	tests := []struct {
		name      string
		analysis  *analyzer.ContentAnalysis
		expectMin float64
		expectMax float64
	}{
		{
			name: "fiction neutral",
			analysis: &analyzer.ContentAnalysis{
				PrimaryGenre:   "fiction",
				SentimentScore: 0.0,
			},
			expectMin: 200.0,
			expectMax: 250.0,
		},
		{
			name: "positive sentiment shifts to warm",
			analysis: &analyzer.ContentAnalysis{
				PrimaryGenre:   "fiction",
				SentimentScore: 0.5,
			},
			expectMin: 30.0,
			expectMax: 60.0,
		},
		{
			name: "negative sentiment stays cool",
			analysis: &analyzer.ContentAnalysis{
				PrimaryGenre:   "fiction",
				SentimentScore: -0.5,
			},
			expectMin: 230.0,
			expectMax: 250.0,
		},
		{
			name: "romance with positive",
			analysis: &analyzer.ContentAnalysis{
				PrimaryGenre:   "romance",
				SentimentScore: 0.5,
			},
			expectMin: 340.0,
			expectMax: 360.0,
		},
		{
			name: "mystery with negative",
			analysis: &analyzer.ContentAnalysis{
				PrimaryGenre:   "mystery",
				SentimentScore: -0.5,
			},
			expectMin: 250.0,
			expectMax: 270.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue := generator.selectPrimaryHue(tt.analysis)
			assert.GreaterOrEqual(t, hue, tt.expectMin)
			assert.LessOrEqual(t, hue, tt.expectMax)
		})
	}
}

func TestColorGenerator_CalculateContrast(t *testing.T) {
	generator := NewColorGenerator()

	tests := []struct {
		name     string
		c1       Color
		c2       Color
		minRatio float64
	}{
		{
			name:     "black on white",
			c1:       Color{R: 0, G: 0, B: 0},
			c2:       Color{R: 255, G: 255, B: 255},
			minRatio: 21.0,
		},
		{
			name:     "white on black",
			c1:       Color{R: 255, G: 255, B: 255},
			c2:       Color{R: 0, G: 0, B: 0},
			minRatio: 21.0,
		},
		{
			name:     "gray on white",
			c1:       Color{R: 118, G: 118, B: 118},
			c2:       Color{R: 255, G: 255, B: 255},
			minRatio: 3.5,
		},
		{
			name:     "dark gray on light gray",
			c1:       Color{R: 68, G: 68, B: 68},
			c2:       Color{R: 238, G: 238, B: 238},
			minRatio: 7.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := generator.calculateContrast(tt.c1, tt.c2)
			assert.GreaterOrEqual(t, ratio, tt.minRatio)
		})
	}
}

func TestColorGenerator_GetAccessibilityLevel(t *testing.T) {
	generator := NewColorGenerator()

	tests := []struct {
		contrast float64
		expected string
	}{
		{21.0, "AAA"},
		{7.5, "AAA"},
		{7.0, "AAA"},
		{6.0, "AA"},
		{4.5, "AA"},
		{4.0, "AA Large"},
		{3.0, "AA Large"},
		{2.5, "Fail"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			level := generator.getAccessibilityLevel(tt.contrast)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestHSLToRGB(t *testing.T) {
	tests := []struct {
		name string
		hsl  HSL
		want Color
	}{
		{
			name: "pure red",
			hsl:  HSL{H: 0, S: 1, L: 0.5},
			want: Color{R: 255, G: 0, B: 0},
		},
		{
			name: "pure green",
			hsl:  HSL{H: 120, S: 1, L: 0.5},
			want: Color{R: 0, G: 255, B: 0},
		},
		{
			name: "pure blue",
			hsl:  HSL{H: 240, S: 1, L: 0.5},
			want: Color{R: 0, G: 0, B: 255},
		},
		{
			name: "white",
			hsl:  HSL{H: 0, S: 0, L: 1},
			want: Color{R: 255, G: 255, B: 255},
		},
		{
			name: "black",
			hsl:  HSL{H: 0, S: 0, L: 0},
			want: Color{R: 0, G: 0, B: 0},
		},
		{
			name: "gray",
			hsl:  HSL{H: 0, S: 0, L: 0.5},
			want: Color{R: 128, G: 128, B: 128},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hslToRGB(tt.hsl)
			
			// Allow small rounding differences
			assert.InDelta(t, tt.want.R, got.R, 1)
			assert.InDelta(t, tt.want.G, got.G, 1)
			assert.InDelta(t, tt.want.B, got.B, 1)
			assert.NotEmpty(t, got.Hex)
		})
	}
}

func TestColorGenerator_DetermineMood(t *testing.T) {
	generator := NewColorGenerator()

	tests := []struct {
		name     string
		analysis *analyzer.ContentAnalysis
		expected string
	}{
		{
			name: "positive sentiment",
			analysis: &analyzer.ContentAnalysis{
				SentimentScore: 0.5,
			},
			expected: "bright",
		},
		{
			name: "negative sentiment",
			analysis: &analyzer.ContentAnalysis{
				SentimentScore: -0.5,
			},
			expected: "dark",
		},
		{
			name: "formal tone",
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{
					Formal: 0.8,
				},
			},
			expected: "professional",
		},
		{
			name: "creative tone",
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{
					Creative: 0.8,
				},
			},
			expected: "vibrant",
		},
		{
			name: "neutral",
			analysis: &analyzer.ContentAnalysis{
				SentimentScore: 0.0,
				Tone: analyzer.ToneProfile{
					Formal: 0.5,
				},
			},
			expected: "balanced",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mood := generator.determineMood(tt.analysis)
			assert.Equal(t, tt.expected, mood)
		})
	}
}

func TestColorGenerator_SelectHarmony(t *testing.T) {
	generator := NewColorGenerator()

	tests := []struct {
		name     string
		analysis *analyzer.ContentAnalysis
		expected string
	}{
		{
			name: "academic prefers analogous",
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{
					Academic: 0.8,
				},
			},
			expected: "analogous",
		},
		{
			name: "creative prefers triadic",
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{
					Creative: 0.8,
				},
			},
			expected: "triadic",
		},
		{
			name: "technical prefers complementary",
			analysis: &analyzer.ContentAnalysis{
				TechnicalDensity: 0.6,
			},
			expected: "complementary",
		},
		{
			name: "default analogous",
			analysis: &analyzer.ContentAnalysis{
				Tone: analyzer.ToneProfile{
					Formal: 0.5,
				},
			},
			expected: "analogous",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			harmony := generator.selectHarmony(tt.analysis)
			assert.Equal(t, tt.expected, harmony)
		})
	}
}

func BenchmarkColorGenerator_Generate(b *testing.B) {
	generator := NewColorGenerator()
	
	analysis := &analyzer.ContentAnalysis{
		PrimaryGenre:   "fiction",
		SentimentScore: 0.3,
		Tone: analyzer.ToneProfile{
			Formal:    0.5,
			Creative:  0.6,
			Academic:  0.3,
			Technical: 0.2,
		},
		Formality:        0.5,
		TechnicalDensity: 0.3,
		Complexity:       0.5,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(analysis)
		if err != nil {
			b.Fatal(err)
		}
	}
}
