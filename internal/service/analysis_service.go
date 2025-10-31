// Package service provides business logic for content analysis
// Operating under Constituição Vértice v3.0 - Artigo I, Seção 3
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/JuanCS-Dev/typecraft/internal/ai"
	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"github.com/JuanCS-Dev/typecraft/internal/repository"
)

// AnalysisService handles manuscript content analysis
type AnalysisService struct {
	analyzer  *ai.Analyzer
	projectRepo *repository.ProjectRepository
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(analyzer *ai.Analyzer, projectRepo *repository.ProjectRepository) *AnalysisService {
	return &AnalysisService{
		analyzer:  analyzer,
		projectRepo: projectRepo,
	}
}

// AnalyzeProject performs AI analysis on a project's manuscript
// Following Artigo VI-VII: Camadas Constitucional e de Deliberação
func (s *AnalysisService) AnalyzeProject(ctx context.Context, projectID string, manuscriptText string) (*domain.AIAnalysis, error) {
	// P5: Consciência Sistêmica - Validate project exists
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Calculate word count
	wordCount := countWords(manuscriptText)
	
	// Extract sample for analysis (first 5000 chars)
	sample := manuscriptText
	if len(sample) > 5000 {
		sample = sample[:5000]
	}

	// Perform AI analysis
	analysis, err := s.analyzer.AnalyzeManuscript(ctx, sample, wordCount)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	// Convert to domain model
	domainAnalysis := s.convertToDomainAnalysis(analysis)
	domainAnalysis.ProjectID = projectID
	
	// Update project status
	project.Status = "analyzing" // Use string instead of constant
	
	if err := s.projectRepo.Update(project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return domainAnalysis, nil
}

// convertToDomainAnalysis converts AI analysis to domain model
func (s *AnalysisService) convertToDomainAnalysis(analysis *ai.ContentAnalysis) *domain.AIAnalysis {
	return &domain.AIAnalysis{
		Genre:             string(analysis.Genre),
		GenreConfidence:   analysis.GenreConfidence,
		SubGenres:         analysis.SubGenres,
		Tone:              convertTone(analysis.Tone),
		Complexity:        convertComplexity(analysis.Complexity),
		EmotionalKeywords: analysis.EmotionalKeywords.Keywords,
		Sentiments:        analysis.EmotionalKeywords.Sentiments,
		HasMath:           analysis.HasMath,
		HasCode:           analysis.HasCode,
		HasImages:         analysis.HasImages,
		WordCount:         analysis.WordCount,
		EstimatedPages:    analysis.EstimatedPages,
		RecommendedPipeline: analysis.RecommendedPipeline,
	}
}

// convertTone converts AI tone to domain tone
func convertTone(tone ai.ToneAnalysis) domain.ToneAnalysis {
	return domain.ToneAnalysis{
		Primary:    tone.Primary,
		Formality:  tone.Formality,
		Emotion:    tone.Emotion,
		Confidence: tone.Confidence,
	}
}

// convertComplexity converts AI complexity to domain complexity
func convertComplexity(complexity ai.ComplexityMetrics) domain.ComplexityMetrics {
	return domain.ComplexityMetrics{
		AvgSentenceLength:  complexity.AvgSentenceLength,
		VocabularyRichness: complexity.VocabularyRichness,
		SyntaxComplexity:   complexity.SyntaxComplexity,
		TechnicalDensity:   complexity.TechnicalDensity,
		ReadingLevel:       complexity.ReadingLevel,
	}
}

// countWords counts words in a text
func countWords(text string) int {
	// Simple word counter (can be improved)
	words := strings.Fields(text)
	return len(words)
}

// GetTypographicRecommendations provides typography recommendations based on analysis
func (s *AnalysisService) GetTypographicRecommendations(analysis *domain.AIAnalysis) *domain.TypographicRecommendations {
	recommendations := &domain.TypographicRecommendations{}

	// Font recommendations based on genre and tone
	recommendations.FontPair = s.recommendFontPair(analysis)
	
	// Layout parameters based on complexity
	recommendations.LayoutParams = s.recommendLayoutParams(analysis)
	
	// Color palette keywords
	recommendations.ColorKeywords = analysis.EmotionalKeywords

	return recommendations
}

// recommendFontPair suggests font pairing based on content analysis
func (s *AnalysisService) recommendFontPair(analysis *domain.AIAnalysis) domain.FontPair {
	// Default: classic serif + sans-serif pair
	pair := domain.FontPair{
		Body:    "Source Serif Pro",
		Heading: "Source Sans Pro",
		Code:    "Source Code Pro",
	}

	// Adjust based on genre
	switch analysis.Genre {
	case "fiction", "poetry", "historical":
		// Classic literature fonts
		pair.Body = "Garamond"
		pair.Heading = "Trajan Pro"
	case "technical", "academic":
		// Clear, readable fonts
		pair.Body = "Source Serif Pro"
		pair.Heading = "Source Sans Pro"
		pair.Code = "Fira Code"
	case "childrens":
		// Friendly, rounded fonts
		pair.Body = "Comic Neue"
		pair.Heading = "Fredoka One"
	case "art", "photography":
		// Modern, elegant fonts
		pair.Body = "Lora"
		pair.Heading = "Montserrat"
	}

	// Adjust based on formality
	if analysis.Tone.Formality > 0.8 {
		// Very formal: use classic serif
		pair.Body = "Baskerville"
	} else if analysis.Tone.Formality < 0.3 {
		// Very casual: use sans-serif for body too
		pair.Body = "Open Sans"
	}

	return pair
}

// recommendLayoutParams suggests layout parameters
func (s *AnalysisService) recommendLayoutParams(analysis *domain.AIAnalysis) domain.LayoutParams {
	params := domain.LayoutParams{
		PageSize:     "6x9",
		BaselineGrid: 12.0,
		Margins: domain.Margins{
			Inner:  0.75,
			Outer:  1.0,
			Top:    1.0,
			Bottom: 1.5,
		},
		BodyFontSize: 11.0,
		Leading:      13.2, // 120% of body size
		GridColumns:  1,
	}

	// Adjust based on complexity
	if analysis.Complexity.TechnicalDensity > 0.7 {
		// Technical content: two-column for code/examples
		params.GridColumns = 2
		params.PageSize = "8.5x11"
	}

	// Adjust based on reading level
	switch analysis.Complexity.ReadingLevel {
	case "elementary", "middle_school":
		// Larger text for younger readers
		params.BodyFontSize = 12.0
		params.Leading = 15.6
	case "graduate", "expert":
		// Denser text for advanced readers
		params.BodyFontSize = 10.0
		params.Leading = 12.0
	}

	return params
}
