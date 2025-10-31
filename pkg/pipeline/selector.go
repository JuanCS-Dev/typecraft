package pipeline

import (
	"fmt"
	"math"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
)

// PipelineType representa o tipo de pipeline
type PipelineType string

const (
	PipelineLaTeX PipelineType = "latex"
	PipelineHTML  PipelineType = "html"
)

// PipelineDecision representa a decisão de pipeline
type PipelineDecision struct {
	Pipeline   PipelineType `json:"pipeline"`
	Confidence float64      `json:"confidence"` // 0-1
	Reasons    []string     `json:"reasons"`
	Scores     DecisionScores `json:"scores"`
}

// DecisionScores mostra os scores de cada pipeline
type DecisionScores struct {
	LaTeXScore float64 `json:"latex_score"`
	HTMLScore  float64 `json:"html_score"`
}

// PipelineThresholds define thresholds para decisão
type PipelineThresholds struct {
	// LaTeX favorito se:
	MathEquationsMin int     // >= N equações
	AcademicConfMin  float64 // >= confiança acadêmica
	ComplexTablesMin int     // >= N tabelas complexas
	
	// HTML favorito se:
	ImageRatioMin    float64 // >= ratio de imagens
	InteractiveMin   int     // >= elementos interativos
	
	// Neutral zone
	NeutralZone      float64 // Diferença mínima para decisão clara
}

// DefaultThresholds retorna thresholds padrão
func DefaultThresholds() PipelineThresholds {
	return PipelineThresholds{
		MathEquationsMin: 10,
		AcademicConfMin:  0.7,
		ComplexTablesMin: 5,
		ImageRatioMin:    0.1,
		InteractiveMin:   3,
		NeutralZone:      0.2,
	}
}

// PipelineSelector seleciona o pipeline apropriado
type PipelineSelector struct {
	thresholds       PipelineThresholds
	academicDetector *analyzer.AcademicDetector
}

// NewPipelineSelector cria um novo seletor
func NewPipelineSelector() *PipelineSelector {
	return &PipelineSelector{
		thresholds:       DefaultThresholds(),
		academicDetector: analyzer.NewAcademicDetector(),
	}
}

// NewPipelineSelectorWithThresholds cria seletor com thresholds customizados
func NewPipelineSelectorWithThresholds(thresholds PipelineThresholds) *PipelineSelector {
	return &PipelineSelector{
		thresholds:       thresholds,
		academicDetector: analyzer.NewAcademicDetector(),
	}
}

// Select seleciona o pipeline baseado na análise
func (ps *PipelineSelector) Select(analysis *analyzer.ContentAnalysis) (*PipelineDecision, error) {
	if analysis == nil {
		return nil, fmt.Errorf("analysis cannot be nil")
	}
	
	// Criar academic score básico a partir da análise existente
	academicScore := &analyzer.AcademicScore{
		IsAcademic:      analysis.IsAcademic(),
		Confidence:      analysis.Tone.Academic,
		EquationDensity: float64(analysis.EquationCount) / math.Max(float64(analysis.WordCount)/250.0, 1.0),
		HasAbstract:     false,
		HasBibliography: false,
	}
	
	// Calcular scores
	latexScore := ps.calculateLaTeXScore(analysis, academicScore)
	htmlScore := ps.calculateHTMLScore(analysis, academicScore)
	
	// Determinar pipeline
	var pipeline PipelineType
	var confidence float64
	reasons := make([]string, 0)
	
	scoreDiff := latexScore - htmlScore
	
	if scoreDiff > ps.thresholds.NeutralZone {
		// LaTeX wins
		pipeline = PipelineLaTeX
		confidence = latexScore / (latexScore + htmlScore)
		reasons = ps.getLaTeXReasons(analysis, academicScore)
	} else if scoreDiff < -ps.thresholds.NeutralZone {
		// HTML wins
		pipeline = PipelineHTML
		confidence = htmlScore / (latexScore + htmlScore)
		reasons = ps.getHTMLReasons(analysis, academicScore)
	} else {
		// Neutral zone - default to HTML (mais flexível)
		pipeline = PipelineHTML
		confidence = 0.5
		reasons = []string{
			"Scores are balanced",
			"Defaulting to HTML for flexibility",
		}
	}
	
	return &PipelineDecision{
		Pipeline:   pipeline,
		Confidence: confidence,
		Reasons:    reasons,
		Scores: DecisionScores{
			LaTeXScore: latexScore,
			HTMLScore:  htmlScore,
		},
	}, nil
}

// calculateLaTeXScore calcula score favorecendo LaTeX
func (ps *PipelineSelector) calculateLaTeXScore(analysis *analyzer.ContentAnalysis, academicScore *analyzer.AcademicScore) float64 {
	score := 0.0
	
	// Equações matemáticas (peso: 30%)
	if analysis.EquationCount >= ps.thresholds.MathEquationsMin {
		score += 0.30
	} else if analysis.EquationCount > 0 {
		// Proporcional
		ratio := float64(analysis.EquationCount) / float64(ps.thresholds.MathEquationsMin)
		score += 0.30 * ratio
	}
	
	// Conteúdo acadêmico (peso: 25%)
	if academicScore.IsAcademic {
		score += 0.25 * academicScore.Confidence
	}
	
	// Tabelas complexas (peso: 15%)
	if analysis.TableCount >= ps.thresholds.ComplexTablesMin {
		score += 0.15
	} else if analysis.TableCount > 0 {
		ratio := float64(analysis.TableCount) / float64(ps.thresholds.ComplexTablesMin)
		score += 0.15 * ratio
	}
	
	// Alta complexidade textual (peso: 15%)
	if analysis.Complexity > 0.7 {
		score += 0.15
	} else if analysis.Complexity > 0.5 {
		score += 0.15 * ((analysis.Complexity - 0.5) / 0.2)
	}
	
	// Densidade de equações alta (peso: 15%)
	if academicScore.EquationDensity > 0.5 {
		score += 0.15
	} else if academicScore.EquationDensity > 0 {
		score += 0.15 * (academicScore.EquationDensity / 0.5)
	}
	
	return score
}

// calculateHTMLScore calcula score favorecendo HTML
func (ps *PipelineSelector) calculateHTMLScore(analysis *analyzer.ContentAnalysis, academicScore *analyzer.AcademicScore) float64 {
	score := 0.0
	
	// Alto ratio de imagens (peso: 30%)
	if analysis.ImageRatio >= ps.thresholds.ImageRatioMin {
		score += 0.30
	} else if analysis.ImageRatio > 0 {
		ratio := analysis.ImageRatio / ps.thresholds.ImageRatioMin
		score += 0.30 * ratio
	}
	
	// Muitas imagens absolutas (peso: 20%)
	if analysis.ImageCount > 20 {
		score += 0.20
	} else if analysis.ImageCount > 5 {
		score += 0.20 * (float64(analysis.ImageCount-5) / 15.0)
	}
	
	// Conteúdo criativo/não-acadêmico (peso: 20%)
	if !academicScore.IsAcademic {
		score += 0.20
	}
	if analysis.Tone.Creative > 0.6 {
		score += 0.10
	}
	
	// Baixa complexidade (mais acessível em HTML) (peso: 15%)
	if analysis.Complexity < 0.4 {
		score += 0.15
	} else if analysis.Complexity < 0.6 {
		score += 0.15 * ((0.6 - analysis.Complexity) / 0.2)
	}
	
	// Tom casual (peso: 15%)
	if analysis.Tone.Casual > 0.5 {
		score += 0.15 * analysis.Tone.Casual
	}
	
	return score
}

// getLaTeXReasons retorna razões para escolher LaTeX
func (ps *PipelineSelector) getLaTeXReasons(analysis *analyzer.ContentAnalysis, academicScore *analyzer.AcademicScore) []string {
	reasons := make([]string, 0)
	
	if analysis.EquationCount >= ps.thresholds.MathEquationsMin {
		reasons = append(reasons, fmt.Sprintf("High equation count (%d)", analysis.EquationCount))
	}
	
	if academicScore.IsAcademic {
		reasons = append(reasons, fmt.Sprintf("Academic content (%.0f%% confidence)", academicScore.Confidence*100))
	}
	
	if analysis.TableCount >= ps.thresholds.ComplexTablesMin {
		reasons = append(reasons, fmt.Sprintf("Complex tables (%d)", analysis.TableCount))
	}
	
	if analysis.Complexity > 0.7 {
		reasons = append(reasons, "High text complexity")
	}
	
	if academicScore.EquationDensity > 0.5 {
		reasons = append(reasons, fmt.Sprintf("High equation density (%.2f per page)", academicScore.EquationDensity))
	}
	
	if len(reasons) == 0 {
		reasons = append(reasons, "LaTeX features detected")
	}
	
	return reasons
}

// getHTMLReasons retorna razões para escolher HTML
func (ps *PipelineSelector) getHTMLReasons(analysis *analyzer.ContentAnalysis, academicScore *analyzer.AcademicScore) []string {
	reasons := make([]string, 0)
	
	if analysis.ImageRatio >= ps.thresholds.ImageRatioMin {
		reasons = append(reasons, fmt.Sprintf("High image ratio (%.2f)", analysis.ImageRatio))
	}
	
	if analysis.ImageCount > 20 {
		reasons = append(reasons, fmt.Sprintf("Many images (%d)", analysis.ImageCount))
	}
	
	if !academicScore.IsAcademic {
		reasons = append(reasons, "Non-academic content")
	}
	
	if analysis.Tone.Creative > 0.6 {
		reasons = append(reasons, "Creative/narrative tone")
	}
	
	if analysis.Complexity < 0.4 {
		reasons = append(reasons, "Accessible reading level")
	}
	
	if analysis.Tone.Casual > 0.5 {
		reasons = append(reasons, "Casual tone")
	}
	
	if len(reasons) == 0 {
		reasons = append(reasons, "Web-first format preferred")
	}
	
	return reasons
}
