package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AIAnalysis representa a análise de IA de um manuscrito
type AIAnalysis struct {
	ID              string             `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID       string             `gorm:"type:uuid;not null;index" json:"project_id"`
	
	// Classificação de gênero
	Genre           string             `gorm:"type:varchar(100);not null" json:"genre"`
	GenreConfidence float64            `gorm:"type:decimal(5,4)" json:"genre_confidence"`
	SubGenres       []string           `gorm:"-" json:"sub_genres,omitempty"` // Stored as JSON
	SubGenresJSON   string             `gorm:"type:text;column:sub_genres" json:"-"`
	
	// Análise de tom
	Tone            ToneAnalysis       `gorm:"embedded;embeddedPrefix:tone_" json:"tone"`
	
	// Métricas de complexidade
	Complexity      ComplexityMetrics  `gorm:"embedded;embeddedPrefix:complexity_" json:"complexity"`
	
	// Palavras-chave emocionais (para paleta de cores)
	EmotionalKeywords []string         `gorm:"-" json:"emotional_keywords,omitempty"`
	EmotionalKeywordsJSON string       `gorm:"type:text;column:emotional_keywords" json:"-"`
	Sentiments       []string          `gorm:"-" json:"sentiments,omitempty"`
	SentimentsJSON   string            `gorm:"type:text;column:sentiments" json:"-"`
	
	// Elementos especiais detectados
	HasMath         bool               `gorm:"type:boolean;default:false" json:"has_math"`
	HasCode         bool               `gorm:"type:boolean;default:false" json:"has_code"`
	HasImages       bool               `gorm:"type:boolean;default:false" json:"has_images"`
	
	// Estatísticas do conteúdo
	WordCount       int                `gorm:"type:integer;not null" json:"word_count"`
	EstimatedPages  int                `gorm:"type:integer;not null" json:"estimated_pages"`
	
	// Recomendações de pipeline
	RecommendedPipeline string         `gorm:"type:varchar(50);not null" json:"recommended_pipeline"`
	
	// Metadados de análise
	AnalyzedAt time.Time `gorm:"not null" json:"analyzed_at"`
	TokensUsed int       `gorm:"type:integer;default:0" json:"tokens_used"`
	
	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// ToneAnalysis embedded struct
type ToneAnalysis struct {
	Primary    string  `gorm:"type:varchar(100)" json:"primary"`
	Formality  float64 `gorm:"type:decimal(5,4)" json:"formality"`
	Emotion    string  `gorm:"type:varchar(100)" json:"emotion"`
	Confidence float64 `gorm:"type:decimal(5,4)" json:"confidence"`
}

// ComplexityMetrics embedded struct
type ComplexityMetrics struct {
	AvgSentenceLength  float64 `gorm:"type:decimal(7,2)" json:"avg_sentence_length"`
	VocabularyRichness float64 `gorm:"type:decimal(5,4)" json:"vocabulary_richness"`
	SyntaxComplexity   float64 `gorm:"type:decimal(5,4)" json:"syntax_complexity"`
	TechnicalDensity   float64 `gorm:"type:decimal(5,4)" json:"technical_density"`
	ReadingLevel       string  `gorm:"type:varchar(50)" json:"reading_level"`
}

// Typography recommendations
type TypographicRecommendations struct {
	FontPair      FontPair      `json:"font_pair"`
	LayoutParams  LayoutParams  `json:"layout_params"`
	ColorKeywords []string      `json:"color_keywords"`
}

// FontPair recommended fonts
type FontPair struct {
	Body    string `json:"body"`
	Heading string `json:"heading"`
	Code    string `json:"code,omitempty"`
}

// LayoutParams layout parameters
type LayoutParams struct {
	PageSize     string  `json:"page_size"`
	BaselineGrid float64 `json:"baseline_grid"`
	Margins      Margins `json:"margins"`
	BodyFontSize float64 `json:"body_font_size"`
	Leading      float64 `json:"leading"`
	GridColumns  int     `json:"grid_columns"`
}

// Margins page margins
type Margins struct {
	Inner  float64 `json:"inner"`
	Outer  float64 `json:"outer"`
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
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
	
	// Serialize slices to JSON
	if len(a.SubGenres) > 0 {
		data, _ := json.Marshal(a.SubGenres)
		a.SubGenresJSON = string(data)
	}
	if len(a.EmotionalKeywords) > 0 {
		data, _ := json.Marshal(a.EmotionalKeywords)
		a.EmotionalKeywordsJSON = string(data)
	}
	if len(a.Sentiments) > 0 {
		data, _ := json.Marshal(a.Sentiments)
		a.SentimentsJSON = string(data)
	}
	
	return nil
}

// AfterFind hook do GORM executado após buscar registro
func (a *AIAnalysis) AfterFind() error {
	// Deserialize JSON to slices
	if a.SubGenresJSON != "" {
		json.Unmarshal([]byte(a.SubGenresJSON), &a.SubGenres)
	}
	if a.EmotionalKeywordsJSON != "" {
		json.Unmarshal([]byte(a.EmotionalKeywordsJSON), &a.EmotionalKeywords)
	}
	if a.SentimentsJSON != "" {
		json.Unmarshal([]byte(a.SentimentsJSON), &a.Sentiments)
	}
	return nil
}

// IsValid valida se a análise tem dados mínimos necessários
func (a *AIAnalysis) IsValid() bool {
	return a.ProjectID != "" &&
		a.Genre != "" &&
		a.Tone.Primary != "" &&
		a.RecommendedPipeline != ""
}

// ShouldUseLaTeX determina se deve usar LaTeX baseado na análise
func (a *AIAnalysis) ShouldUseLaTeX() bool {
	return a.HasMath ||
		a.HasCode ||
		a.Complexity.TechnicalDensity > 0.5 ||
		a.RecommendedPipeline == "latex"
}
