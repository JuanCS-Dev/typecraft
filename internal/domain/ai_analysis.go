package domain

import (
	"time"

	"github.com/google/uuid"
)

// AIAnalysis representa a análise de IA de um manuscrito
type AIAnalysis struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	ManuscriptID string    `gorm:"type:uuid;not null;index" json:"manuscript_id"`
	
	// Classificação de gênero
	Genre           string  `gorm:"type:varchar(100);not null" json:"genre"`
	SubGenre        string  `gorm:"type:varchar(100)" json:"sub_genre"`
	GenreConfidence float64 `gorm:"type:decimal(5,4)" json:"genre_confidence"`
	
	// Análise de tom
	Tone      string  `gorm:"type:varchar(100);not null" json:"tone"`
	ToneScore float64 `gorm:"type:decimal(5,4)" json:"tone_score"`
	
	// Elementos especiais detectados
	EquationPercentage float64 `gorm:"type:decimal(5,4)" json:"equation_percentage"`
	CodePercentage     float64 `gorm:"type:decimal(5,4)" json:"code_percentage"`
	TableCount         int     `gorm:"type:integer;default:0" json:"table_count"`
	ImageCount         int     `gorm:"type:integer;default:0" json:"image_count"`
	
	// Recomendações de pipeline
	RecommendedPipeline string  `gorm:"type:varchar(50);not null" json:"recommended_pipeline"`
	PipelineConfidence  float64 `gorm:"type:decimal(5,4)" json:"pipeline_confidence"`
	PipelineReason      string  `gorm:"type:text" json:"pipeline_reason"`
	
	// Recomendações de fontes
	RecommendedBodyFont  string `gorm:"type:varchar(100)" json:"recommended_body_font"`
	RecommendedTitleFont string `gorm:"type:varchar(100)" json:"recommended_title_font"`
	RecommendedMonoFont  string `gorm:"type:varchar(100)" json:"recommended_mono_font"`
	FontRationale        string `gorm:"type:text" json:"font_rationale"`
	
	// Metadados de análise
	AnalyzedAt time.Time `gorm:"not null" json:"analyzed_at"`
	TokensUsed int       `gorm:"type:integer;default:0" json:"tokens_used"`
	
	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
}

// TableName especifica o nome da tabela no banco de dados
func (AIAnalysis) TableName() string {
	return "ai_analyses"
}

// BeforeCreate hook do GORM executado antes de criar o registro
func (a *AIAnalysis) BeforeCreate() error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.AnalyzedAt.IsZero() {
		a.AnalyzedAt = time.Now()
	}
	return nil
}

// IsValid valida se a análise tem dados mínimos necessários
func (a *AIAnalysis) IsValid() bool {
	return a.ManuscriptID != "" &&
		a.Genre != "" &&
		a.Tone != "" &&
		a.RecommendedPipeline != ""
}

// SpecialElements retorna um resumo dos elementos especiais detectados
type SpecialElements struct {
	HasEquations bool    `json:"has_equations"`
	HasCode      bool    `json:"has_code"`
	HasTables    bool    `json:"has_tables"`
	HasImages    bool    `json:"has_images"`
	Complexity   string  `json:"complexity"`
	Score        float64 `json:"score"`
}

// GetSpecialElements retorna análise agregada dos elementos especiais
func (a *AIAnalysis) GetSpecialElements() SpecialElements {
	se := SpecialElements{
		HasEquations: a.EquationPercentage > 0.01,
		HasCode:      a.CodePercentage > 0.01,
		HasTables:    a.TableCount > 0,
		HasImages:    a.ImageCount > 0,
	}
	
	// Calcular score de complexidade (0-1)
	score := (a.EquationPercentage + a.CodePercentage) / 2
	if se.HasTables {
		score += 0.1
	}
	if se.HasImages {
		score += 0.05
	}
	
	se.Score = score
	
	// Classificar complexidade
	switch {
	case score < 0.05:
		se.Complexity = "simple"
	case score < 0.15:
		se.Complexity = "moderate"
	default:
		se.Complexity = "complex"
	}
	
	return se
}

// ShouldUseLaTeX determina se deve usar LaTeX baseado na análise
func (a *AIAnalysis) ShouldUseLaTeX() bool {
	return a.EquationPercentage > 0.05 ||
		a.CodePercentage > 0.05 ||
		a.TableCount > 10 ||
		a.RecommendedPipeline == "latex"
}
