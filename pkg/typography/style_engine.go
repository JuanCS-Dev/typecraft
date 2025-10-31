package typography

import (
	"regexp"
	"strings"
)

// StyleEngine aplica regras tipográficas ao texto
type StyleEngine struct {
	rules []Rule
}

// Rule representa uma regra tipográfica
type Rule struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
	Enabled     bool
}

// NewStyleEngine cria uma nova engine de estilo com regras padrão
func NewStyleEngine() *StyleEngine {
	return &StyleEngine{
		rules: defaultRules(),
	}
}

// ApplyRules aplica todas as regras tipográficas ao texto
func (s *StyleEngine) ApplyRules(text string) string {
	result := text

	for _, rule := range s.rules {
		if !rule.Enabled {
			continue
		}

		if rule.Pattern != nil {
			result = rule.Pattern.ReplaceAllString(result, rule.Replacement)
		}
	}

	return result
}

// defaultRules retorna as regras tipográficas padrão em português
func defaultRules() []Rule {
	return []Rule{
		{
			Name:        "Aspas duplas tipográficas",
			Pattern:     regexp.MustCompile(`"([^"]+)"`),
			Replacement: "\u201c$1\u201d",
			Enabled:     true,
		},
		{
			Name:        "Aspas simples tipográficas",
			Pattern:     regexp.MustCompile(`'([^']+)'`),
			Replacement: "\u2018$1\u2019",
			Enabled:     true,
		},
		{
			Name:        "Reticências",
			Pattern:     regexp.MustCompile(`\.\.\.`),
			Replacement: "…",
			Enabled:     true,
		},
		{
			Name:        "Travessão em diálogos",
			Pattern:     regexp.MustCompile(`(?m)^-\s+`),
			Replacement: "— ",
			Enabled:     true,
		},
		{
			Name:        "Espaço não-quebrável antes de pontuação",
			Pattern:     regexp.MustCompile(`\s+([!?;:])`),
			Replacement: " $1",
			Enabled:     true,
		},
		{
			Name:        "Múltiplos espaços",
			Pattern:     regexp.MustCompile(`\s{2,}`),
			Replacement: " ",
			Enabled:     true,
		},
		{
			Name:        "Espaço após pontuação",
			Pattern:     regexp.MustCompile(`([.!?;:,])\s*([^\s])`),
			Replacement: "$1 $2",
			Enabled:     true,
		},
	}
}

// AddRule adiciona uma nova regra customizada
func (s *StyleEngine) AddRule(rule Rule) {
	s.rules = append(s.rules, rule)
}

// DisableRule desabilita uma regra por nome
func (s *StyleEngine) DisableRule(name string) {
	for i := range s.rules {
		if s.rules[i].Name == name {
			s.rules[i].Enabled = false
			return
		}
	}
}

// EnableRule habilita uma regra por nome
func (s *StyleEngine) EnableRule(name string) {
	for i := range s.rules {
		if s.rules[i].Name == name {
			s.rules[i].Enabled = true
			return
		}
	}
}

// ApplySmartQuotes aplica aspas inteligentes
func ApplySmartQuotes(text string) string {
	// Aspas duplas
	text = regexp.MustCompile(`"([^"]+)"`).ReplaceAllString(text, "\u201c$1\u201d")
	// Aspas simples  
	text = regexp.MustCompile(`'([^']+)'`).ReplaceAllString(text, "\u2018$1\u2019")
	return text
}

// ApplyDashes converte hífens em travessões quando apropriado
func ApplyDashes(text string) string {
	// Travessão em diálogos
	text = regexp.MustCompile(`(?m)^-\s+`).ReplaceAllString(text, "— ")
	// Travessão entre frases
	text = regexp.MustCompile(`\s+-\s+`).ReplaceAllString(text, " — ")
	return text
}

// ApplyEllipsis converte três pontos em reticências
func ApplyEllipsis(text string) string {
	return strings.ReplaceAll(text, "...", "…")
}

// CleanSpacing limpa espaçamento excessivo
func CleanSpacing(text string) string {
	// Remove múltiplos espaços
	text = regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " ")
	// Remove espaços no início de linhas
	text = regexp.MustCompile(`(?m)^\s+`).ReplaceAllString(text, "")
	// Remove espaços no final de linhas
	text = regexp.MustCompile(`(?m)\s+$`).ReplaceAllString(text, "")
	return text
}

// FormatParagraphs adiciona formatação adequada a parágrafos
func FormatParagraphs(text string) string {
	// Garante que parágrafos tenham quebra dupla
	paragraphs := strings.Split(text, "\n\n")
	var formatted []string

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			// Remove quebras simples dentro do parágrafo
			p = regexp.MustCompile(`\n`).ReplaceAllString(p, " ")
			// Limpa espaços múltiplos
			p = regexp.MustCompile(`\s{2,}`).ReplaceAllString(p, " ")
			formatted = append(formatted, p)
		}
	}

	return strings.Join(formatted, "\n\n")
}
