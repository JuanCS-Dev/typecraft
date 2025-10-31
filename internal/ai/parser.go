package ai

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
)

type AnalysisResponse struct {
	Genre      string  `json:"genre"`
	SubGenre   string  `json:"sub_genre"`
	Confidence float64 `json:"confidence"`
	
	Tone      string  `json:"tone"`
	ToneScore float64 `json:"tone_score"`
	
	SpecialElements struct {
		EquationPercentage float64 `json:"equation_percentage"`
		CodePercentage     float64 `json:"code_percentage"`
		TableCount         int     `json:"table_count"`
		ImageCount         int     `json:"image_count"`
	} `json:"special_elements"`
	
	Recommendations struct {
		Pipeline         string  `json:"pipeline"`
		PipelineReason   string  `json:"pipeline_reason"`
		Confidence       float64 `json:"confidence"`
		BodyFont         string  `json:"body_font"`
		TitleFont        string  `json:"title_font"`
		MonoFont         string  `json:"mono_font"`
		FontRationale    string  `json:"font_rationale"`
	} `json:"recommendations"`
}

func ParseAnalysisResponse(jsonResponse, manuscriptID string) (*domain.AIAnalysis, error) {
	var response AnalysisResponse
	
	if err := json.Unmarshal([]byte(jsonResponse), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	analysis := &domain.AIAnalysis{
		ManuscriptID:         manuscriptID,
		Genre:                response.Genre,
		SubGenre:             response.SubGenre,
		GenreConfidence:      response.Confidence,
		Tone:                 response.Tone,
		ToneScore:            response.ToneScore,
		EquationPercentage:   response.SpecialElements.EquationPercentage,
		CodePercentage:       response.SpecialElements.CodePercentage,
		TableCount:           response.SpecialElements.TableCount,
		ImageCount:           response.SpecialElements.ImageCount,
		RecommendedPipeline:  response.Recommendations.Pipeline,
		PipelineConfidence:   response.Recommendations.Confidence,
		PipelineReason:       response.Recommendations.PipelineReason,
		RecommendedBodyFont:  response.Recommendations.BodyFont,
		RecommendedTitleFont: response.Recommendations.TitleFont,
		RecommendedMonoFont:  response.Recommendations.MonoFont,
		FontRationale:        response.Recommendations.FontRationale,
		AnalyzedAt:           time.Now(),
	}

	if !analysis.IsValid() {
		return nil, fmt.Errorf("invalid analysis data")
	}

	return analysis, nil
}
