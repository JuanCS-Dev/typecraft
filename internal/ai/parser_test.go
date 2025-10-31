package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/JuanCS-Dev/typecraft/internal/domain"
)

func TestParseAnalysisResponse(t *testing.T) {
	tests := []struct {
		name        string
		jsonResp    string
		manuscriptID string
		wantErr     bool
		checkFunc   func(*testing.T, interface{})
	}{
		{
			name: "valid response - fiction",
			jsonResp: `{
				"genre": "fiction",
				"sub_genre": "romance",
				"confidence": 0.92,
				"tone": "conversational",
				"tone_score": 0.85,
				"special_elements": {
					"equation_percentage": 0.0,
					"code_percentage": 0.0,
					"table_count": 0,
					"image_count": 2
				},
				"recommendations": {
					"pipeline": "standard",
					"pipeline_reason": "Standard fiction with dialogue and chapters",
					"confidence": 0.88,
					"body_font": "Garamond",
					"title_font": "Cinzel",
					"mono_font": "Courier",
					"font_rationale": "Classic serif for readability"
				}
			}`,
			manuscriptID: "test-manuscript-123",
			wantErr: false,
			checkFunc: func(t *testing.T, result interface{}) {
				analysis := result.(*domain.AIAnalysis)
				assert.Equal(t, "fiction", analysis.Genre)
				assert.Equal(t, "romance", analysis.SubGenre)
				assert.Equal(t, 0.92, analysis.GenreConfidence)
				assert.Equal(t, "conversational", analysis.Tone)
				assert.Equal(t, 0.85, analysis.ToneScore)
				assert.Equal(t, 0.0, analysis.EquationPercentage)
				assert.Equal(t, 2, analysis.ImageCount)
				assert.Equal(t, "standard", analysis.RecommendedPipeline)
				assert.Equal(t, "Garamond", analysis.RecommendedBodyFont)
			},
		},
		{
			name: "valid response - technical",
			jsonResp: `{
				"genre": "technical",
				"sub_genre": "computer_science",
				"confidence": 0.95,
				"tone": "formal",
				"tone_score": 0.90,
				"special_elements": {
					"equation_percentage": 0.15,
					"code_percentage": 0.25,
					"table_count": 8,
					"image_count": 12
				},
				"recommendations": {
					"pipeline": "latex",
					"pipeline_reason": "High code and equation content requires LaTeX",
					"confidence": 0.93,
					"body_font": "Computer Modern",
					"title_font": "Latin Modern Sans",
					"mono_font": "Inconsolata",
					"font_rationale": "Technical fonts for code and math"
				}
			}`,
			manuscriptID: "test-tech-456",
			wantErr: false,
			checkFunc: func(t *testing.T, result interface{}) {
				analysis := result.(*domain.AIAnalysis)
				assert.Equal(t, "technical", analysis.Genre)
				assert.Equal(t, 0.15, analysis.EquationPercentage)
				assert.Equal(t, 0.25, analysis.CodePercentage)
				assert.Equal(t, "latex", analysis.RecommendedPipeline)
				assert.True(t, analysis.ShouldUseLaTeX())
			},
		},
		{
			name: "invalid JSON",
			jsonResp: `{invalid json}`,
			manuscriptID: "test-789",
			wantErr: true,
		},
		{
			name: "missing required fields",
			jsonResp: `{
				"genre": "fiction",
				"tone": ""
			}`,
			manuscriptID: "test-999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAnalysisResponse(tt.jsonResp, tt.manuscriptID)
			
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			require.NotNil(t, result)
			
			assert.Equal(t, tt.manuscriptID, result.ManuscriptID)
			assert.True(t, result.IsValid())
			
			if tt.checkFunc != nil {
				tt.checkFunc(t, result)
			}
		})
	}
}

func TestParseAnalysisResponse_EdgeCases(t *testing.T) {
	t.Run("empty manuscript ID", func(t *testing.T) {
		jsonResp := `{
			"genre": "fiction",
			"confidence": 0.9,
			"tone": "formal",
			"tone_score": 0.8,
			"special_elements": {"equation_percentage": 0, "code_percentage": 0, "table_count": 0, "image_count": 0},
			"recommendations": {"pipeline": "simple", "pipeline_reason": "test", "confidence": 0.9, "body_font": "Arial", "title_font": "Arial", "mono_font": "Courier", "font_rationale": "test"}
		}`
		
		result, err := ParseAnalysisResponse(jsonResp, "")
		require.Error(t, err)
		assert.Nil(t, result)
	})
	
	t.Run("zero confidence should still work", func(t *testing.T) {
		jsonResp := `{
			"genre": "unknown",
			"confidence": 0.0,
			"tone": "neutral",
			"tone_score": 0.0,
			"special_elements": {"equation_percentage": 0, "code_percentage": 0, "table_count": 0, "image_count": 0},
			"recommendations": {"pipeline": "simple", "pipeline_reason": "test", "confidence": 0.0, "body_font": "Arial", "title_font": "Arial", "mono_font": "Courier", "font_rationale": "test"}
		}`
		
		result, err := ParseAnalysisResponse(jsonResp, "test-id")
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0.0, result.GenreConfidence)
	})
}
