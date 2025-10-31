package ai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyzer_AnalyzeManuscript_Fiction tests analysis of fiction content
func TestAnalyzer_AnalyzeManuscript_Fiction(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	analyzer := NewAnalyzer(apiKey, "gpt-4o")
	
	fictionSample := `
		It was a dark and stormy night when Elizabeth finally understood the truth. 
		The old mansion on Ravencroft Hill had held secrets for generations, but tonight, 
		under the flickering candlelight, she discovered the journal that would change everything.
		
		Her hands trembled as she opened the leather-bound book. The handwriting was elegant,
		almost too perfect. "My dearest daughter," it began, "if you are reading this,
		then I am no longer among the living. What I must tell you now will challenge
		everything you believe about our family."
		
		Outside, thunder rolled across the moors. Elizabeth pulled her shawl tighter
		and continued reading. Page after page revealed a tapestry of lies, love, and
		betrayal that spanned three continents and two world wars. Her grandmother hadn't
		been a simple governess after all. She had been a spy.
	`

	analysis, err := analyzer.AnalyzeManuscript(context.Background(), fictionSample, 50000)
	
	require.NoError(t, err, "Analysis should not error")
	require.NotNil(t, analysis, "Analysis should not be nil")

	// Assertions on genre
	assert.Contains(t, []Genre{GenreFiction, GenreMystery, GenreHistorical}, analysis.Genre,
		"Should detect fiction, mystery, or historical genre")
	assert.GreaterOrEqual(t, analysis.GenreConfidence, 0.6,
		"Should have reasonable confidence")

	// Assertions on tone
	assert.NotEmpty(t, analysis.Tone.Primary, "Should have primary tone")
	assert.GreaterOrEqual(t, analysis.Tone.Formality, 0.5,
		"Fiction prose should be moderately to highly formal")
	assert.NotEmpty(t, analysis.Tone.Emotion, "Should detect emotional tone")

	// Assertions on complexity
	assert.Greater(t, analysis.Complexity.AvgSentenceLength, 0.0,
		"Should calculate sentence length")
	assert.Greater(t, analysis.Complexity.VocabularyRichness, 0.0,
		"Should assess vocabulary")
	
	// Assertions on word count
	assert.Equal(t, 50000, analysis.WordCount)
	assert.Greater(t, analysis.EstimatedPages, 0)

	// Log results for manual inspection
	t.Logf("Fiction Analysis Results:\n%s", analyzer.GetTypographicProfile(analysis))
}

// TestAnalyzer_AnalyzeManuscript_Technical tests analysis of technical content
func TestAnalyzer_AnalyzeManuscript_Technical(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	analyzer := NewAnalyzer(apiKey, "gpt-4o")
	
	technicalSample := `
		# Introduction to Quantum Computing

		## 1.1 Quantum Bits and Superposition

		In classical computing, information is stored in bits that exist in either
		a 0 or 1 state. Quantum computers, however, utilize quantum bits or "qubits"
		that can exist in a superposition of both states simultaneously.

		Mathematically, a qubit state is represented as:
		|ψ⟩ = α|0⟩ + β|1⟩

		where α and β are complex probability amplitudes satisfying |α|² + |β|² = 1.

		## 1.2 Quantum Gates

		Quantum gates manipulate qubits through unitary transformations. The Hadamard
		gate (H) creates superposition:

		H|0⟩ = (|0⟩ + |1⟩)/√2

		Consider the following Python implementation:

		[CODE]
		from qiskit import QuantumCircuit, execute, Aer

		def create_bell_state():
			qc = QuantumCircuit(2, 2)
			qc.h(0)  # Apply Hadamard to qubit 0
			qc.cx(0, 1)  # CNOT gate
			return qc
		[/CODE]

		This circuit creates a maximally entangled Bell state, demonstrating
		quantum entanglement—a phenomenon where qubits become correlated such
		that measuring one instantly affects the other, regardless of distance.
	`

	analysis, err := analyzer.AnalyzeManuscript(context.Background(), technicalSample, 30000)
	
	require.NoError(t, err)
	require.NotNil(t, analysis)

	// Technical genre should be detected
	assert.Equal(t, GenreTechnical, analysis.Genre,
		"Should detect technical genre")

	// Should detect math
	assert.True(t, analysis.HasMath,
		"Should detect mathematical notation")

	// Should detect code
	assert.True(t, analysis.HasCode,
		"Should detect code blocks")

	// Should recommend LaTeX for technical content with math
	assert.Equal(t, "latex", analysis.RecommendedPipeline,
		"Should recommend LaTeX pipeline for technical/math content")

	// Technical content should have high technical density
	assert.GreaterOrEqual(t, analysis.Complexity.TechnicalDensity, 0.5,
		"Should detect high technical density")

	t.Logf("Technical Analysis Results:\n%s", analyzer.GetTypographicProfile(analysis))
}

// TestAnalyzer_AnalyzeManuscript_Poetry tests analysis of poetic content
func TestAnalyzer_AnalyzeManuscript_Poetry(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	analyzer := NewAnalyzer(apiKey, "gpt-4o")
	
	poetrySample := `
		Whispers of Dawn

		In the hushed moment before sunrise,
		when shadows dance with fading stars,
		the world holds its breath—
		a silent prayer suspended
		between night's velvet embrace
		and morning's golden promise.

		Dewdrops cling to spider silk,
		each one a universe of light,
		reflecting the infinite
		in the ephemeral,
		beauty in the transient.

		Listen:
		the first bird's song
		shatters the crystalline silence,
		cascading notes of pure joy
		awakening the sleeping earth.

		And we, dreamers of daybreak,
		stand witness to this ancient ritual—
		the daily miracle
		of darkness becoming light,
		of endings becoming beginnings,
		of sorrow transforming into hope.
	`

	analysis, err := analyzer.AnalyzeManuscript(context.Background(), poetrySample, 15000)
	
	require.NoError(t, err)
	require.NotNil(t, analysis)

	// Should detect poetry
	assert.Equal(t, GenrePoetry, analysis.Genre,
		"Should detect poetry genre")

	// Poetry should have rich vocabulary
	assert.GreaterOrEqual(t, analysis.Complexity.VocabularyRichness, 0.6,
		"Poetry should have rich vocabulary")

	// Should extract emotional keywords
	assert.NotEmpty(t, analysis.EmotionalKeywords.Keywords,
		"Should extract emotional keywords")
	assert.NotEmpty(t, analysis.EmotionalKeywords.Sentiments,
		"Should identify sentiments")

	// Emotional keywords should capture poetic imagery
	keywords := analyzer.ExtractEmotionalKeywordsForColors(analysis)
	assert.NotEmpty(t, keywords, "Should have keywords for color palette")

	t.Logf("Poetry Analysis Results:\n%s", analyzer.GetTypographicProfile(analysis))
	t.Logf("Emotional Keywords for Colors: %v", keywords)
}

// TestAnalyzer_Validation tests that validation catches incomplete analysis
func TestAnalyzer_Validation(t *testing.T) {
	analyzer := NewAnalyzer("dummy_key", "gpt-4o")

	tests := []struct {
		name      string
		analysis  *ContentAnalysis
		wantError bool
	}{
		{
			name: "valid analysis",
			analysis: &ContentAnalysis{
				Genre:           GenreFiction,
				GenreConfidence: 0.85,
				Tone: ToneAnalysis{
					Primary:   "formal",
					Formality: 0.7,
					Emotion:   "serene",
				},
				Complexity: ComplexityMetrics{
					ReadingLevel: "college",
				},
				RecommendedPipeline: "latex",
			},
			wantError: false,
		},
		{
			name: "missing genre",
			analysis: &ContentAnalysis{
				Tone: ToneAnalysis{Primary: "formal"},
				Complexity: ComplexityMetrics{ReadingLevel: "college"},
				RecommendedPipeline: "latex",
			},
			wantError: true,
		},
		{
			name: "invalid genre_confidence",
			analysis: &ContentAnalysis{
				Genre:           GenreFiction,
				GenreConfidence: 1.5, // invalid
				Tone: ToneAnalysis{Primary: "formal"},
				Complexity: ComplexityMetrics{ReadingLevel: "college"},
				RecommendedPipeline: "latex",
			},
			wantError: true,
		},
		{
			name: "missing tone.primary",
			analysis: &ContentAnalysis{
				Genre:           GenreFiction,
				GenreConfidence: 0.8,
				Tone: ToneAnalysis{}, // empty
				Complexity: ComplexityMetrics{ReadingLevel: "college"},
				RecommendedPipeline: "latex",
			},
			wantError: true,
		},
		{
			name: "invalid pipeline",
			analysis: &ContentAnalysis{
				Genre:           GenreFiction,
				GenreConfidence: 0.8,
				Tone: ToneAnalysis{Primary: "formal"},
				Complexity: ComplexityMetrics{ReadingLevel: "college"},
				RecommendedPipeline: "invalid", // invalid
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := analyzer.validateAnalysis(tt.analysis)
			if tt.wantError {
				assert.Error(t, err, "Should return error for invalid analysis")
			} else {
				assert.NoError(t, err, "Should not error for valid analysis")
			}
		})
	}
}

// TestAnalyzer_SelfCritique tests the self-critique logic
func TestAnalyzer_SelfCritique(t *testing.T) {
	analyzer := NewAnalyzer("dummy_key", "gpt-4o")

	t.Run("corrects pipeline for math content", func(t *testing.T) {
		analysis := &ContentAnalysis{
			Genre:               GenreTechnical,
			HasMath:             true,
			RecommendedPipeline: "html", // wrong for math
		}

		analyzer.selfCritique(analysis, "")

		assert.Equal(t, "latex", analysis.RecommendedPipeline,
			"Should correct to LaTeX for math content")
	})

	t.Run("corrects reading_level for rich vocabulary", func(t *testing.T) {
		analysis := &ContentAnalysis{
			Complexity: ComplexityMetrics{
				VocabularyRichness: 0.9,
				ReadingLevel:       "elementary", // inconsistent
			},
		}

		analyzer.selfCritique(analysis, "")

		assert.NotEqual(t, "elementary", analysis.Complexity.ReadingLevel,
			"Should elevate reading level for rich vocabulary")
	})

	t.Run("sets HTML for visual genres", func(t *testing.T) {
		visualGenres := []Genre{GenreCookbook, GenreArt, GenreTravel, GenreChildrens}

		for _, genre := range visualGenres {
			analysis := &ContentAnalysis{
				Genre:               genre,
				RecommendedPipeline: "latex", // wrong for visual content
			}

			analyzer.selfCritique(analysis, "")

			assert.Equal(t, "html", analysis.RecommendedPipeline,
				"Should use HTML pipeline for genre: %s", genre)
		}
	})
}

// TestAnalyzer_ShouldUseLaTeX tests pipeline decision helper
func TestAnalyzer_ShouldUseLaTeX(t *testing.T) {
	analyzer := NewAnalyzer("dummy_key", "gpt-4o")

	tests := []struct {
		name     string
		pipeline string
		want     bool
	}{
		{"latex pipeline", "latex", true},
		{"html pipeline", "html", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &ContentAnalysis{
				RecommendedPipeline: tt.pipeline,
			}
			got := analyzer.ShouldUseLaTeX(analysis)
			assert.Equal(t, tt.want, got)
		})
	}
}
