// Package ai provides AI-powered content analysis and design generation
// Following DETER-AGENT framework from Constituição Vértice v3.0
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Genre represents the detected genre of a manuscript
type Genre string

const (
	GenreFiction          Genre = "fiction"
	GenreNonFiction       Genre = "nonfiction"
	GenreTechnical        Genre = "technical"
	GenreAcademic         Genre = "academic"
	GenrePoetry           Genre = "poetry"
	GenreChildrens        Genre = "childrens"
	GenreSelfHelp         Genre = "self_help"
	GenreBiography        Genre = "biography"
	GenreHistorical       Genre = "historical"
	GenreSciFi            Genre = "sci_fi"
	GenreFantasy          Genre = "fantasy"
	GenreRomance          Genre = "romance"
	GenreMystery          Genre = "mystery"
	GenreCookbook         Genre = "cookbook"
	GenreTravel           Genre = "travel"
	GenreArt              Genre = "art"
)

// ToneAnalysis represents the emotional tone of the content
type ToneAnalysis struct {
	Primary    string  `json:"primary"`     // e.g., "formal", "casual", "poetic", "technical"
	Formality  float64 `json:"formality"`   // 0.0 (very casual) to 1.0 (very formal)
	Emotion    string  `json:"emotion"`     // e.g., "serene", "intense", "playful"
	Confidence float64 `json:"confidence"`  // 0.0 to 1.0
}

// ComplexityMetrics measures text complexity
type ComplexityMetrics struct {
	AvgSentenceLength    float64 `json:"avg_sentence_length"`
	VocabularyRichness   float64 `json:"vocabulary_richness"`   // 0.0 to 1.0
	SyntaxComplexity     float64 `json:"syntax_complexity"`     // 0.0 to 1.0
	TechnicalDensity     float64 `json:"technical_density"`     // 0.0 to 1.0
	ReadingLevel         string  `json:"reading_level"`         // e.g., "high_school", "college", "graduate"
}

// EmotionalKeywords extracted for color palette generation
type EmotionalKeywords struct {
	Keywords   []string `json:"keywords"`
	Sentiments []string `json:"sentiments"` // e.g., "joyful", "melancholic", "energetic"
}

// ContentAnalysis is the complete analysis of a manuscript
type ContentAnalysis struct {
	Genre             Genre              `json:"genre"`
	GenreConfidence   float64            `json:"genre_confidence"`
	SubGenres         []string           `json:"sub_genres,omitempty"`
	Tone              ToneAnalysis       `json:"tone"`
	Complexity        ComplexityMetrics  `json:"complexity"`
	EmotionalKeywords EmotionalKeywords  `json:"emotional_keywords"`
	HasMath           bool               `json:"has_math"`
	HasCode           bool               `json:"has_code"`
	HasImages         bool               `json:"has_images"`
	WordCount         int                `json:"word_count"`
	EstimatedPages    int                `json:"estimated_pages"`
	RecommendedPipeline string           `json:"recommended_pipeline"` // "latex" or "html"
}

// Analyzer performs AI-powered content analysis
type Analyzer struct {
	client *openai.Client
	model  string
}

// NewAnalyzer creates a new content analyzer
func NewAnalyzer(apiKey string, model string) *Analyzer {
	if model == "" {
		model = "gpt-4o"
	}
	return &Analyzer{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

// AnalyzeManuscript performs comprehensive analysis of manuscript content
// Following Artigo VI: Camada Constitucional - Princípios P2 (Validação Preventiva) e P4 (Rastreabilidade)
func (a *Analyzer) AnalyzeManuscript(ctx context.Context, textSample string, fullWordCount int) (*ContentAnalysis, error) {
	// P6: Eficiência de Token - limitar sample se muito grande
	if len(textSample) > 5000 {
		textSample = textSample[:5000]
	}

	// Tree of Thoughts: Phase 1 - Generate multiple analysis approaches
	prompt := a.buildAnalysisPrompt(textSample, fullWordCount)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: a.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: `You are an expert literary analyst and typographer operating under the Constituição Vértice v3.0.
Your task is to analyze manuscript content and provide structured, accurate analysis for automated book design.

CRITICAL PRINCIPLES (P1-P6):
- P1: Completude Obrigatória - Provide complete analysis, no placeholders
- P2: Validação Preventiva - Validate all assumptions before concluding
- P3: Ceticismo Crítico - Question obvious classifications, dig deeper
- P4: Rastreabilidade Total - Base all conclusions on evidence from the text
- P5: Consciência Sistêmica - Consider how each field impacts downstream design decisions
- P6: Eficiência de Token - Be concise yet complete

Output MUST be valid JSON matching the ContentAnalysis schema exactly.`,
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Temperature: 0.3, // Lower temperature for analytical task
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: "json_object",
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var analysis ContentAnalysis
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// P2: Validação Preventiva - validate analysis completeness
	if err := a.validateAnalysis(&analysis); err != nil {
		return nil, fmt.Errorf("invalid analysis from AI: %w", err)
	}

	// Set full word count
	analysis.WordCount = fullWordCount
	
	// Calculate estimated pages (250 words/page average)
	analysis.EstimatedPages = (fullWordCount / 250) + 1

	// Auto-critica obrigatória (Artigo VII, Seção 2)
	a.selfCritique(&analysis, textSample)

	return &analysis, nil
}

// buildAnalysisPrompt constructs the prompt for manuscript analysis
func (a *Analyzer) buildAnalysisPrompt(textSample string, wordCount int) string {
	return fmt.Sprintf(`Analyze this manuscript sample and provide a comprehensive ContentAnalysis.

MANUSCRIPT SAMPLE (first ~5000 chars):
"""
%s
"""

METADATA:
- Total word count: %d words

ANALYSIS REQUIREMENTS:

1. GENRE DETECTION:
   - Primary genre (from: fiction, nonfiction, technical, academic, poetry, childrens, self_help, biography, historical, sci_fi, fantasy, romance, mystery, cookbook, travel, art)
   - Confidence level (0.0-1.0)
   - Sub-genres if applicable

2. TONE ANALYSIS:
   - Primary tone descriptor (formal, casual, poetic, technical, conversational, authoritative, playful)
   - Formality score (0.0=very casual, 1.0=very formal)
   - Emotional tone (serene, intense, melancholic, joyful, energetic, contemplative)
   - Confidence in tone assessment

3. COMPLEXITY METRICS:
   - Average sentence length
   - Vocabulary richness (0.0=simple, 1.0=very sophisticated)
   - Syntax complexity (0.0=simple, 1.0=very complex)
   - Technical density (0.0=no jargon, 1.0=highly technical)
   - Reading level (elementary, middle_school, high_school, college, graduate, expert)

4. EMOTIONAL KEYWORDS:
   - Extract 5-10 words with strong emotional/aesthetic charge
   - Identify dominant sentiments (for color palette generation)

5. CONTENT FEATURES:
   - has_math: true if mathematical notation detected
   - has_code: true if programming code detected
   - has_images: true if references to images/figures detected

6. PIPELINE RECOMMENDATION:
   - "latex" for: technical, academic, heavy math, complex footnotes
   - "html" for: visual-rich, cookbooks, children's books, art books

Return ONLY valid JSON matching this exact structure (no markdown, no extra text):
{
  "genre": "string",
  "genre_confidence": 0.0,
  "sub_genres": ["string"],
  "tone": {
    "primary": "string",
    "formality": 0.0,
    "emotion": "string",
    "confidence": 0.0
  },
  "complexity": {
    "avg_sentence_length": 0.0,
    "vocabulary_richness": 0.0,
    "syntax_complexity": 0.0,
    "technical_density": 0.0,
    "reading_level": "string"
  },
  "emotional_keywords": {
    "keywords": ["string"],
    "sentiments": ["string"]
  },
  "has_math": false,
  "has_code": false,
  "has_images": false,
  "recommended_pipeline": "latex"
}`, textSample, wordCount)
}

// validateAnalysis ensures the AI output is complete and valid (P1: Completude Obrigatória)
func (a *Analyzer) validateAnalysis(analysis *ContentAnalysis) error {
	if analysis.Genre == "" {
		return fmt.Errorf("genre is required")
	}
	
	if analysis.GenreConfidence < 0 || analysis.GenreConfidence > 1 {
		return fmt.Errorf("genre_confidence must be between 0 and 1")
	}

	if analysis.Tone.Primary == "" {
		return fmt.Errorf("tone.primary is required")
	}

	if analysis.Complexity.ReadingLevel == "" {
		return fmt.Errorf("complexity.reading_level is required")
	}

	if analysis.RecommendedPipeline != "latex" && analysis.RecommendedPipeline != "html" {
		return fmt.Errorf("recommended_pipeline must be 'latex' or 'html'")
	}

	return nil
}

// selfCritique applies self-criticism to the analysis (Artigo VII: Camada de Deliberação)
func (a *Analyzer) selfCritique(analysis *ContentAnalysis, textSample string) {
	// Questões de auto-crítica obrigatórias:
	// - O gênero está correto ou foi uma classificação superficial?
	// - O pipeline recomendado faz sentido para o conteúdo?
	// - As métricas de complexidade são coerentes entre si?

	// Ajuste 1: Se há muita matemática mas pipeline é HTML, corrigir
	if analysis.HasMath && analysis.RecommendedPipeline == "html" {
		analysis.RecommendedPipeline = "latex"
	}

	// Ajuste 2: Se vocabulário rico mas reading_level é básico, questionar
	if analysis.Complexity.VocabularyRichness > 0.7 && 
	   (analysis.Complexity.ReadingLevel == "elementary" || analysis.Complexity.ReadingLevel == "middle_school") {
		// Elevar reading level
		if analysis.Complexity.VocabularyRichness > 0.85 {
			analysis.Complexity.ReadingLevel = "college"
		} else {
			analysis.Complexity.ReadingLevel = "high_school"
		}
	}

	// Ajuste 3: Se detectou código ou imagens mas não refletiu no pipeline
	if (analysis.HasCode || analysis.HasImages) && analysis.Genre == "technical" {
		// LaTeX é superior para documentação técnica
		analysis.RecommendedPipeline = "latex"
	}

	// Ajuste 4: Livros visuais (cookbook, art, travel, childrens) devem usar HTML
	visualGenres := []Genre{GenreCookbook, GenreArt, GenreTravel, GenreChildrens}
	for _, vg := range visualGenres {
		if analysis.Genre == vg {
			analysis.RecommendedPipeline = "html"
			break
		}
	}
}

// ExtractEmotionalKeywordsForColors extracts keywords specifically for color palette generation
// This is used by the Design Generator module
func (a *Analyzer) ExtractEmotionalKeywordsForColors(analysis *ContentAnalysis) []string {
	// Combine keywords and sentiments
	allKeywords := append(analysis.EmotionalKeywords.Keywords, analysis.EmotionalKeywords.Sentiments...)
	
	// Filter to most relevant for color (top 5)
	if len(allKeywords) > 5 {
		return allKeywords[:5]
	}
	
	return allKeywords
}

// ShouldUseLaTeX determines if LaTeX pipeline is recommended
func (a *Analyzer) ShouldUseLaTeX(analysis *ContentAnalysis) bool {
	return analysis.RecommendedPipeline == "latex"
}

// GetTypographicProfile generates a summary for typography decisions
func (a *Analyzer) GetTypographicProfile(analysis *ContentAnalysis) string {
	var profile strings.Builder
	
	profile.WriteString(fmt.Sprintf("Genre: %s (%.0f%% confidence)\n", analysis.Genre, analysis.GenreConfidence*100))
	profile.WriteString(fmt.Sprintf("Tone: %s, %s\n", analysis.Tone.Primary, analysis.Tone.Emotion))
	profile.WriteString(fmt.Sprintf("Formality: %.0f%%\n", analysis.Tone.Formality*100))
	profile.WriteString(fmt.Sprintf("Reading Level: %s\n", analysis.Complexity.ReadingLevel))
	profile.WriteString(fmt.Sprintf("Complexity: Vocabulary=%.0f%%, Syntax=%.0f%%, Technical=%.0f%%\n",
		analysis.Complexity.VocabularyRichness*100,
		analysis.Complexity.SyntaxComplexity*100,
		analysis.Complexity.TechnicalDensity*100))
	profile.WriteString(fmt.Sprintf("Special Features: Math=%v, Code=%v, Images=%v\n",
		analysis.HasMath, analysis.HasCode, analysis.HasImages))
	profile.WriteString(fmt.Sprintf("Recommended Pipeline: %s\n", analysis.RecommendedPipeline))
	
	return profile.String()
}
