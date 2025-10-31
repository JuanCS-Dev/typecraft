package design

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
)

// FontPair representa um par de fontes (corpo + título)
type FontPair struct {
	Body      string  `json:"body"`
	Heading   string  `json:"heading"`
	Monospace string  `json:"monospace,omitempty"`
	Mood      string  `json:"mood"`
	Rationale string  `json:"rationale"`
	Score     float64 `json:"score"`
}

// FontDatabase representa o banco de dados de fontes
type FontDatabase struct {
	Version       string                        `json:"version"`
	LastUpdated   string                        `json:"last_updated"`
	CuratedFonts  map[string][]string          `json:"curated_fonts"`
	GenrePairings map[string][]FontPair        `json:"genre_pairings"`
	FontMetadata  map[string]FontMetadata      `json:"font_metadata"`
}

// FontMetadata representa metadata de uma fonte
type FontMetadata struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	SupportsLatin bool   `json:"supports_latin"`
	Variable      bool   `json:"variable"`
	GoogleFonts   bool   `json:"google_fonts"`
}

// FontSuggester sugere fontes baseado em análise de conteúdo
type FontSuggester struct {
	db *FontDatabase
}

// NewFontSuggester cria um novo sugestor de fontes
func NewFontSuggester() (*FontSuggester, error) {
	db, err := loadFontDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to load font database: %w", err)
	}
	
	return &FontSuggester{db: db}, nil
}

// Suggest sugere pares de fontes baseado na análise de conteúdo
func (fs *FontSuggester) Suggest(analysis *analyzer.ContentAnalysis) ([]FontPair, error) {
	if analysis == nil {
		return nil, fmt.Errorf("analysis cannot be nil")
	}
	
	// 1. Obter pairings do gênero primário
	primaryPairs := fs.getPairsForGenre(analysis.PrimaryGenre)
	
	// 2. Obter pairings do gênero secundário (se houver)
	var secondaryPairs []FontPair
	if analysis.SecondaryGenre != "" {
		secondaryPairs = fs.getPairsForGenre(analysis.SecondaryGenre)
	}
	
	// 3. Calcular scores baseados nas métricas de análise
	allPairs := append(primaryPairs, secondaryPairs...)
	for i := range allPairs {
		allPairs[i].Score = fs.calculateScore(&allPairs[i], analysis)
	}
	
	// 4. Ordenar por score
	sort.Slice(allPairs, func(i, j int) bool {
		return allPairs[i].Score > allPairs[j].Score
	})
	
	// 5. Retornar top 3
	maxResults := 3
	if len(allPairs) < maxResults {
		maxResults = len(allPairs)
	}
	
	return allPairs[:maxResults], nil
}

// getPairsForGenre retorna pares de fontes para um gênero
func (fs *FontSuggester) getPairsForGenre(genre string) []FontPair {
	if pairs, ok := fs.db.GenrePairings[genre]; ok {
		result := make([]FontPair, len(pairs))
		copy(result, pairs)
		return result
	}
	
	// Default: fiction pairings
	if pairs, ok := fs.db.GenrePairings["fiction"]; ok {
		result := make([]FontPair, len(pairs))
		copy(result, pairs)
		return result
	}
	
	return []FontPair{}
}

// calculateScore calcula score de um par baseado na análise
func (fs *FontSuggester) calculateScore(pair *FontPair, analysis *analyzer.ContentAnalysis) float64 {
	score := 1.0
	
	// Boost para matches de mood
	if pair.Mood == "classic" && analysis.Tone.Formal > 0.6 {
		score += 0.3
	}
	if pair.Mood == "modern" && analysis.Tone.Formal < 0.4 {
		score += 0.3
	}
	if pair.Mood == "technical" && analysis.TechnicalDensity > 0.5 {
		score += 0.4
	}
	if pair.Mood == "academic" && analysis.Tone.Academic > 0.6 {
		score += 0.4
	}
	if pair.Mood == "romantic" && analysis.SentimentScore > 0.3 {
		score += 0.3
	}
	if pair.Mood == "dark" && analysis.SentimentScore < -0.2 {
		score += 0.3
	}
	
	// Boost para complexidade
	if analysis.Complexity > 0.7 && pair.Mood == "scholarly" {
		score += 0.2
	}
	if analysis.Complexity < 0.4 && pair.Mood == "accessible" {
		score += 0.2
	}
	
	// Boost se tem código e o pair tem monospace
	if analysis.TechnicalDensity > 0.3 && pair.Monospace != "" {
		score += 0.2
	}
	
	// Penalidade para pares muito formais em conteúdo casual
	if pair.Mood == "classic" && analysis.Tone.Casual > 0.5 {
		score -= 0.3
	}
	
	return score
}

// GetFontMetadata retorna metadata de uma fonte
func (fs *FontSuggester) GetFontMetadata(fontName string) (*FontMetadata, bool) {
	if meta, ok := fs.db.FontMetadata[fontName]; ok {
		return &meta, true
	}
	return nil, false
}

// ListAvailableFonts lista todas as fontes disponíveis por categoria
func (fs *FontSuggester) ListAvailableFonts(category string) []string {
	if fonts, ok := fs.db.CuratedFonts[category]; ok {
		result := make([]string, len(fonts))
		copy(result, fonts)
		return result
	}
	return []string{}
}

// loadFontDatabase carrega o banco de dados de fontes
func loadFontDatabase() (*FontDatabase, error) {
	// Try multiple paths
	paths := []string{
		"pkg/design/google_fonts_db.json",
		"google_fonts_db.json",
		"../../pkg/design/google_fonts_db.json",
	}
	
	var lastErr error
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}
		
		var db FontDatabase
		if err := json.Unmarshal(data, &db); err != nil {
			return nil, fmt.Errorf("failed to parse font database: %w", err)
		}
		
		return &db, nil
	}
	
	return nil, fmt.Errorf("failed to read font database from any path: %w", lastErr)
}
