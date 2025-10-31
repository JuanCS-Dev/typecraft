package epub

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MetadataEnhancer melhora metadata de ePub
type MetadataEnhancer struct {
	languageDetector *LanguageDetector
}

// NewMetadataEnhancer cria um novo enhancer
func NewMetadataEnhancer() *MetadataEnhancer {
	return &MetadataEnhancer{
		languageDetector: NewLanguageDetector(),
	}
}

// Enhance melhora metadata automaticamente
func (m *MetadataEnhancer) Enhance(metadata *Metadata, content string) {
	// 1. Generate UUID if missing
	if metadata.Identifier == "" {
		metadata.Identifier = m.GenerateUUID()
	}
	
	// 2. Detect language if missing
	if metadata.Language == "" && content != "" {
		metadata.Language = m.languageDetector.Detect(content)
	}
	
	// 3. Set date if missing
	if metadata.Date.IsZero() {
		metadata.Date = time.Now()
	}
	
	// 4. Enhance subjects with taxonomy
	if len(metadata.Subject) > 0 {
		metadata.Subject = m.EnhanceSubjects(metadata.Subject)
	}
	
	// 5. Normalize fields
	metadata.Title = m.NormalizeTitle(metadata.Title)
	metadata.Author = m.NormalizeAuthor(metadata.Author)
}

// GenerateUUID gera um UUID para o livro
func (m *MetadataEnhancer) GenerateUUID() string {
	return fmt.Sprintf("urn:uuid:%s", uuid.New().String())
}

// NormalizeTitle normaliza o título
func (m *MetadataEnhancer) NormalizeTitle(title string) string {
	title = strings.TrimSpace(title)
	
	// Remove multiple spaces
	for strings.Contains(title, "  ") {
		title = strings.ReplaceAll(title, "  ", " ")
	}
	
	return title
}

// NormalizeAuthor normaliza o nome do autor
func (m *MetadataEnhancer) NormalizeAuthor(author string) string {
	author = strings.TrimSpace(author)
	
	// Remove multiple spaces
	for strings.Contains(author, "  ") {
		author = strings.ReplaceAll(author, "  ", " ")
	}
	
	return author
}

// EnhanceSubjects melhora subjects com taxonomia
func (m *MetadataEnhancer) EnhanceSubjects(subjects []string) []string {
	enhanced := make([]string, 0, len(subjects))
	seen := make(map[string]bool)
	
	for _, subject := range subjects {
		subject = strings.TrimSpace(subject)
		if subject == "" {
			continue
		}
		
		// Avoid duplicates
		lower := strings.ToLower(subject)
		if seen[lower] {
			continue
		}
		seen[lower] = true
		
		enhanced = append(enhanced, subject)
	}
	
	return enhanced
}

// ValidateMetadata valida metadata
func (m *MetadataEnhancer) ValidateMetadata(metadata *Metadata) []string {
	issues := make([]string, 0)
	
	// Required fields
	if metadata.Title == "" {
		issues = append(issues, "Title is required")
	}
	
	if metadata.Author == "" {
		issues = append(issues, "Author is required")
	}
	
	if metadata.Language == "" {
		issues = append(issues, "Language is required")
	}
	
	if metadata.Identifier == "" {
		issues = append(issues, "Identifier is required")
	}
	
	// Recommended fields
	if metadata.Publisher == "" {
		issues = append(issues, "Publisher is recommended")
	}
	
	if metadata.Description == "" {
		issues = append(issues, "Description is recommended")
	}
	
	if len(metadata.Subject) == 0 {
		issues = append(issues, "At least one subject is recommended")
	}
	
	return issues
}

// LanguageDetector detecta idioma do conteúdo
type LanguageDetector struct {
	patterns map[string][]string
}

// NewLanguageDetector cria um novo detector
func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{
		patterns: map[string][]string{
			"en": {"the", "and", "is", "in", "to", "of", "a"},
			"pt": {"o", "a", "de", "e", "em", "os", "as", "da", "do"},
			"es": {"el", "la", "de", "y", "en", "los", "las", "del"},
			"fr": {"le", "la", "de", "et", "les", "des", "un", "une"},
			"de": {"der", "die", "das", "und", "in", "den", "dem"},
			"it": {"il", "la", "di", "e", "in", "i", "le", "del"},
		},
	}
}

// Detect detecta o idioma do conteúdo
func (ld *LanguageDetector) Detect(content string) string {
	if content == "" {
		return "en" // default
	}
	
	// Normalize content
	content = strings.ToLower(content)
	words := strings.Fields(content)
	
	if len(words) < 10 {
		return "en" // too short to detect
	}
	
	// Count matches for each language
	scores := make(map[string]int)
	
	for lang, patterns := range ld.patterns {
		score := 0
		for _, word := range words {
			for _, pattern := range patterns {
				if word == pattern {
					score++
				}
			}
		}
		scores[lang] = score
	}
	
	// Find highest score
	maxScore := 0
	detected := "en"
	
	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			detected = lang
		}
	}
	
	// If no clear winner, default to English
	if maxScore < 3 {
		return "en"
	}
	
	return detected
}

// SubjectTaxonomy taxonomia de subjects
type SubjectTaxonomy struct {
	categories map[string][]string
}

// NewSubjectTaxonomy cria uma nova taxonomia
func NewSubjectTaxonomy() *SubjectTaxonomy {
	return &SubjectTaxonomy{
		categories: map[string][]string{
			"Fiction": {
				"Literary Fiction",
				"Science Fiction",
				"Fantasy",
				"Mystery",
				"Thriller",
				"Romance",
				"Horror",
				"Historical Fiction",
			},
			"Non-Fiction": {
				"Biography",
				"History",
				"Science",
				"Technology",
				"Business",
				"Self-Help",
				"Travel",
				"Cooking",
			},
			"Academic": {
				"Mathematics",
				"Physics",
				"Chemistry",
				"Biology",
				"Computer Science",
				"Engineering",
				"Medicine",
				"Philosophy",
			},
			"Reference": {
				"Dictionary",
				"Encyclopedia",
				"Manual",
				"Handbook",
			},
		},
	}
}

// GetCategory retorna a categoria de um subject
func (st *SubjectTaxonomy) GetCategory(subject string) string {
	subject = strings.ToLower(subject)
	
	for category, subjects := range st.categories {
		for _, s := range subjects {
			if strings.ToLower(s) == subject {
				return category
			}
		}
	}
	
	return "General"
}

// GetRelated retorna subjects relacionados
func (st *SubjectTaxonomy) GetRelated(subject string) []string {
	category := st.GetCategory(subject)
	if category == "General" {
		return []string{}
	}
	
	related := make([]string, 0)
	if subjects, ok := st.categories[category]; ok {
		for _, s := range subjects {
			if strings.ToLower(s) != strings.ToLower(subject) {
				related = append(related, s)
			}
		}
	}
	
	return related
}
