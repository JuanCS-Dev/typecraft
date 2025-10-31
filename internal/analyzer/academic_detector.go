package analyzer

import (
	"regexp"
	"strings"
)

// AcademicDetector detecta se o conteúdo é acadêmico/científico
type AcademicDetector struct {
	academicKeywords []string
	citationPatterns []*regexp.Regexp
	structureMarkers []string
}

// NewAcademicDetector cria um novo detector
func NewAcademicDetector() *AcademicDetector {
	return &AcademicDetector{
		academicKeywords: initAcademicKeywords(),
		citationPatterns: initCitationPatterns(),
		structureMarkers: initStructureMarkers(),
	}
}

// AcademicScore representa o score de conteúdo acadêmico
type AcademicScore struct {
	IsAcademic       bool    `json:"is_academic"`
	Confidence       float64 `json:"confidence"`        // 0-1
	KeywordDensity   float64 `json:"keyword_density"`   // Academic keywords per 1000 words
	CitationCount    int     `json:"citation_count"`    // Number of citations found
	HasBibliography  bool    `json:"has_bibliography"`  // References section detected
	HasAbstract      bool    `json:"has_abstract"`      // Abstract section detected
	StructureScore   float64 `json:"structure_score"`   // 0-1 based on academic structure
	EquationDensity  float64 `json:"equation_density"`  // Equations per page
}

// Detect analisa se o conteúdo é acadêmico
func (ad *AcademicDetector) Detect(content string, analysis *ContentAnalysis) *AcademicScore {
	score := &AcademicScore{}
	
	contentLower := strings.ToLower(content)
	wordCount := float64(analysis.WordCount)
	if wordCount == 0 {
		wordCount = 1
	}
	
	// 1. Contar keywords acadêmicas
	keywordCount := ad.countAcademicKeywords(contentLower)
	score.KeywordDensity = (float64(keywordCount) / wordCount) * 1000
	
	// 2. Detectar citações
	score.CitationCount = ad.countCitations(content)
	
	// 3. Detectar estrutura acadêmica
	score.HasAbstract = ad.hasAbstract(contentLower)
	score.HasBibliography = ad.hasBibliography(contentLower)
	score.StructureScore = ad.calculateStructureScore(contentLower)
	
	// 4. Densidade de equações
	if analysis.WordCount > 0 {
		estimatedPages := float64(analysis.WordCount) / 250.0
		if estimatedPages > 0 {
			score.EquationDensity = float64(analysis.EquationCount) / estimatedPages
		}
	}
	
	// 5. Calcular confidence e decisão
	score.Confidence = ad.calculateConfidence(score, analysis)
	score.IsAcademic = score.Confidence > 0.6
	
	return score
}

// countAcademicKeywords conta palavras acadêmicas
func (ad *AcademicDetector) countAcademicKeywords(content string) int {
	count := 0
	for _, keyword := range ad.academicKeywords {
		count += strings.Count(content, keyword)
	}
	return count
}

// countCitations conta citações no texto
func (ad *AcademicDetector) countCitations(content string) int {
	count := 0
	for _, pattern := range ad.citationPatterns {
		matches := pattern.FindAllString(content, -1)
		count += len(matches)
	}
	return count
}

// hasAbstract detecta se tem seção de abstract
func (ad *AcademicDetector) hasAbstract(content string) bool {
	abstractMarkers := []string{
		"abstract",
		"resumo",
		"summary",
		"síntese",
	}
	
	for _, marker := range abstractMarkers {
		if strings.Contains(content, marker) {
			return true
		}
	}
	return false
}

// hasBibliography detecta se tem seção de referências
func (ad *AcademicDetector) hasBibliography(content string) bool {
	bibMarkers := []string{
		"references",
		"bibliography",
		"referências",
		"bibliografia",
		"works cited",
		"obras citadas",
	}
	
	for _, marker := range bibMarkers {
		if strings.Contains(content, marker) {
			return true
		}
	}
	return false
}

// calculateStructureScore calcula score baseado em estrutura acadêmica
func (ad *AcademicDetector) calculateStructureScore(content string) float64 {
	score := 0.0
	maxScore := float64(len(ad.structureMarkers))
	
	for _, marker := range ad.structureMarkers {
		if strings.Contains(content, marker) {
			score += 1.0
		}
	}
	
	if maxScore > 0 {
		return score / maxScore
	}
	return 0.0
}

// calculateConfidence calcula confidence final
func (ad *AcademicDetector) calculateConfidence(score *AcademicScore, analysis *ContentAnalysis) float64 {
	confidence := 0.0
	
	// Keyword density (peso: 0.25)
	if score.KeywordDensity > 10 {
		confidence += 0.25
	} else if score.KeywordDensity > 5 {
		confidence += 0.15
	} else if score.KeywordDensity > 2 {
		confidence += 0.05
	}
	
	// Citations (peso: 0.25)
	if score.CitationCount > 20 {
		confidence += 0.25
	} else if score.CitationCount > 10 {
		confidence += 0.15
	} else if score.CitationCount > 5 {
		confidence += 0.10
	} else if score.CitationCount > 0 {
		confidence += 0.05
	}
	
	// Structure (peso: 0.20)
	if score.HasAbstract {
		confidence += 0.10
	}
	if score.HasBibliography {
		confidence += 0.10
	}
	
	// Structure markers (peso: 0.15)
	confidence += score.StructureScore * 0.15
	
	// Tone analysis (peso: 0.10)
	if analysis.Tone.Academic > 0.7 {
		confidence += 0.10
	} else if analysis.Tone.Academic > 0.5 {
		confidence += 0.05
	}
	
	// Formality (peso: 0.05)
	if analysis.Formality > 0.7 {
		confidence += 0.05
	}
	
	return confidence
}

// initAcademicKeywords inicializa keywords acadêmicas
func initAcademicKeywords() []string {
	return []string{
		// Metodologia
		"hypothesis", "hipótese", "methodology", "metodologia",
		"experiment", "experimento", "analysis", "análise",
		"research", "pesquisa", "study", "estudo",
		"investigation", "investigação", "findings", "achados",
		
		// Estrutura acadêmica
		"abstract", "resumo", "introduction", "introdução",
		"conclusion", "conclusão", "discussion", "discussão",
		"results", "resultados", "methods", "métodos",
		
		// Linguagem científica
		"significant", "significativo", "correlation", "correlação",
		"statistical", "estatístico", "data", "dados",
		"variable", "variável", "sample", "amostra",
		"theory", "teoria", "model", "modelo",
		
		// Citações e referências
		"according to", "de acordo com", "et al", "et al.",
		"cited", "citado", "references", "referências",
		"published", "publicado", "journal", "revista",
		
		// Verbos acadêmicos
		"demonstrate", "demonstrar", "indicate", "indicar",
		"suggest", "sugerir", "conclude", "concluir",
		"observe", "observar", "examine", "examinar",
		"analyze", "analisar", "evaluate", "avaliar",
		
		// Qualificadores
		"therefore", "portanto", "however", "entretanto",
		"furthermore", "além disso", "consequently", "consequentemente",
		"nevertheless", "não obstante", "moreover", "ademais",
	}
}

// initCitationPatterns inicializa padrões de citação
func initCitationPatterns() []*regexp.Regexp {
	patterns := []string{
		// APA style: (Author, Year)
		`\([A-Z][a-z]+,\s*\d{4}\)`,
		`\([A-Z][a-z]+\s+et\s+al\.,\s*\d{4}\)`,
		
		// Vancouver style: [1], [1-3], [1,2,3]
		`\[\d+\]`,
		`\[\d+-\d+\]`,
		`\[\d+,\s*\d+\]`,
		
		// Harvard style: Author (Year)
		`[A-Z][a-z]+\s+\(\d{4}\)`,
		
		// Et al.
		`et\s+al\.`,
		
		// DOI
		`doi:\s*10\.\d+`,
		`https?://doi\.org/`,
		
		// ISSN/ISBN
		`ISSN\s+\d{4}-\d{4}`,
		`ISBN\s+\d{3}-\d+-\d+-\d+-\d`,
	}
	
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	
	return compiled
}

// initStructureMarkers inicializa marcadores de estrutura
func initStructureMarkers() []string {
	return []string{
		"abstract",
		"introduction",
		"methodology",
		"results",
		"discussion",
		"conclusion",
		"references",
		"appendix",
		"acknowledgments",
		"figure",
		"table",
		"equation",
	}
}
