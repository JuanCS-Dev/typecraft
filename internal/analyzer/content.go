package analyzer

import (
	"math"
	"regexp"
	"strings"
)

// ContentAnalysis representa a análise completa de um manuscrito
type ContentAnalysis struct {
	// Genre signals: palavras-chave e seus pesos
	GenreSignals map[string]float64 `json:"genre_signals"`
	
	// Métricas de gênero detectadas
	PrimaryGenre   string   `json:"primary_genre"`
	SecondaryGenre string   `json:"secondary_genre,omitempty"`
	GenreScores    map[string]float64 `json:"genre_scores"`
	
	// Métricas de complexidade
	Complexity       float64 `json:"complexity"`        // 0-1 (Flesch Reading Ease normalizado)
	FleschScore      float64 `json:"flesch_score"`      // Score bruto
	AvgSentenceLen   float64 `json:"avg_sentence_len"`
	AvgWordLen       float64 `json:"avg_word_len"`
	
	// Métricas de tom
	Tone             ToneProfile `json:"tone"`
	Formality        float64     `json:"formality"`      // 0-1
	TechnicalDensity float64     `json:"technical_density"` // 0-1
	
	// Contadores
	WordCount     int `json:"word_count"`
	SentenceCount int `json:"sentence_count"`
	ParagraphCount int `json:"paragraph_count"`
	ImageCount    int `json:"image_count"`
	TableCount    int `json:"table_count"`
	EquationCount int `json:"equation_count"`
	
	// Métricas de sentimento
	SentimentScore float64 `json:"sentiment_score"` // -1 a +1
	
	// Ratios úteis
	ImageRatio float64 `json:"image_ratio"` // images/pages estimado
}

// ToneProfile descreve o tom do conteúdo
type ToneProfile struct {
	Formal     float64 `json:"formal"`
	Casual     float64 `json:"casual"`
	Technical  float64 `json:"technical"`
	Creative   float64 `json:"creative"`
	Academic   float64 `json:"academic"`
}

// ContentAnalyzer analisa manuscritos
type ContentAnalyzer struct {
	genreKeywords map[string][]string
	techTerms     map[string]bool
	formalWords   []string
	casualWords   []string
}

// NewContentAnalyzer cria um novo analisador
func NewContentAnalyzer() *ContentAnalyzer {
	return &ContentAnalyzer{
		genreKeywords: initGenreKeywords(),
		techTerms:     initTechnicalTerms(),
		formalWords:   initFormalWords(),
		casualWords:   initCasualWords(),
	}
}

// Analyze realiza análise completa do conteúdo
func (ca *ContentAnalyzer) Analyze(content string) (*ContentAnalysis, error) {
	analysis := &ContentAnalysis{
		GenreSignals: make(map[string]float64),
		GenreScores:  make(map[string]float64),
	}
	
	ca.countBasicMetrics(content, analysis)
	ca.analyzeGenre(content, analysis)
	ca.analyzeComplexity(content, analysis)
	ca.analyzeTone(content, analysis)
	ca.analyzeSpecialElements(content, analysis)
	ca.analyzeSentiment(content, analysis)
	ca.calculateRatios(analysis)
	
	return analysis, nil
}

func (ca *ContentAnalyzer) countBasicMetrics(content string, analysis *ContentAnalysis) {
	cleaned := strings.TrimSpace(content)
	
	paragraphs := strings.Split(cleaned, "\n\n")
	analysis.ParagraphCount = len(paragraphs)
	
	sentenceRegex := regexp.MustCompile(`[.!?]+`)
	sentences := sentenceRegex.Split(cleaned, -1)
	analysis.SentenceCount = 0
	for _, s := range sentences {
		if strings.TrimSpace(s) != "" {
			analysis.SentenceCount++
		}
	}
	
	wordRegex := regexp.MustCompile(`\b\w+\b`)
	words := wordRegex.FindAllString(cleaned, -1)
	analysis.WordCount = len(words)
	
	if analysis.SentenceCount > 0 {
		analysis.AvgSentenceLen = float64(analysis.WordCount) / float64(analysis.SentenceCount)
	}
	
	totalChars := 0
	for _, word := range words {
		totalChars += len(word)
	}
	if len(words) > 0 {
		analysis.AvgWordLen = float64(totalChars) / float64(len(words))
	}
}

func (ca *ContentAnalyzer) analyzeGenre(content string, analysis *ContentAnalysis) {
	contentLower := strings.ToLower(content)
	
	for genre, keywords := range ca.genreKeywords {
		score := 0.0
		for _, keyword := range keywords {
			count := strings.Count(contentLower, strings.ToLower(keyword))
			if count > 0 {
				score += float64(count)
				analysis.GenreSignals[keyword] = float64(count)
			}
		}
		analysis.GenreScores[genre] = score
	}
	
	maxScore := 0.0
	for _, score := range analysis.GenreScores {
		if score > maxScore {
			maxScore = score
		}
	}
	
	if maxScore > 0 {
		for genre := range analysis.GenreScores {
			analysis.GenreScores[genre] /= maxScore
		}
	}
	
	primary := ""
	secondary := ""
	primaryScore := 0.0
	secondaryScore := 0.0
	
	for genre, score := range analysis.GenreScores {
		if score > primaryScore {
			secondary = primary
			secondaryScore = primaryScore
			primary = genre
			primaryScore = score
		} else if score > secondaryScore {
			secondary = genre
			secondaryScore = score
		}
	}
	
	analysis.PrimaryGenre = primary
	if secondaryScore > 0.3 {
		analysis.SecondaryGenre = secondary
	}
}

func (ca *ContentAnalyzer) analyzeComplexity(content string, analysis *ContentAnalysis) {
	if analysis.SentenceCount == 0 || analysis.WordCount == 0 {
		analysis.Complexity = 0.5
		analysis.FleschScore = 50.0
		return
	}
	
	syllables := ca.countSyllables(content)
	
	avgWordsPerSentence := float64(analysis.WordCount) / float64(analysis.SentenceCount)
	avgSyllablesPerWord := float64(syllables) / float64(analysis.WordCount)
	
	flesch := 206.835 - (1.015 * avgWordsPerSentence) - (84.6 * avgSyllablesPerWord)
	analysis.FleschScore = flesch
	
	if flesch < 0 {
		flesch = 0
	}
	if flesch > 100 {
		flesch = 100
	}
	
	analysis.Complexity = 1.0 - (flesch / 100.0)
}

func (ca *ContentAnalyzer) countSyllables(text string) int {
	wordRegex := regexp.MustCompile(`\b\w+\b`)
	words := wordRegex.FindAllString(strings.ToLower(text), -1)
	
	totalSyllables := 0
	for _, word := range words {
		syllables := ca.syllablesInWord(word)
		totalSyllables += syllables
	}
	
	return totalSyllables
}

func (ca *ContentAnalyzer) syllablesInWord(word string) int {
	if len(word) == 0 {
		return 0
	}
	
	vowels := "aeiouáéíóúâêôãõy"
	syllables := 0
	prevWasVowel := false
	
	for _, char := range strings.ToLower(word) {
		isVowel := strings.ContainsRune(vowels, char)
		if isVowel && !prevWasVowel {
			syllables++
		}
		prevWasVowel = isVowel
	}
	
	if strings.HasSuffix(word, "e") && syllables > 1 {
		syllables--
	}
	
	if syllables == 0 {
		syllables = 1
	}
	
	return syllables
}

func (ca *ContentAnalyzer) analyzeTone(content string, analysis *ContentAnalysis) {
	contentLower := strings.ToLower(content)
	words := strings.Fields(contentLower)
	
	if len(words) == 0 {
		return
	}
	
	formalCount := 0
	casualCount := 0
	technicalCount := 0
	
	for _, formalWord := range ca.formalWords {
		formalCount += strings.Count(contentLower, formalWord)
	}
	
	for _, casualWord := range ca.casualWords {
		casualCount += strings.Count(contentLower, casualWord)
	}
	
	for _, word := range words {
		if ca.techTerms[word] {
			technicalCount++
		}
	}
	
	totalWords := float64(len(words))
	analysis.Tone.Formal = float64(formalCount) / totalWords
	analysis.Tone.Casual = float64(casualCount) / totalWords
	analysis.Tone.Technical = float64(technicalCount) / totalWords
	
	if formalCount+casualCount > 0 {
		analysis.Formality = float64(formalCount) / float64(formalCount+casualCount)
	} else {
		analysis.Formality = 0.5
	}
	
	analysis.TechnicalDensity = math.Min(float64(technicalCount)/totalWords*10, 1.0)
	
	analysis.Tone.Academic = (analysis.Formality + analysis.Complexity + analysis.TechnicalDensity) / 3.0
	
	analysis.Tone.Creative = 1.0 - analysis.Tone.Academic
}

func (ca *ContentAnalyzer) analyzeSpecialElements(content string, analysis *ContentAnalysis) {
	imageRegex := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	images := imageRegex.FindAllString(content, -1)
	analysis.ImageCount = len(images)
	
	tableRegex := regexp.MustCompile(`\|.*\|`)
	tables := tableRegex.FindAllString(content, -1)
	analysis.TableCount = len(tables) / 3
	
	equationRegex := regexp.MustCompile(`\$\$.*?\$\$|\$.*?\$`)
	equations := equationRegex.FindAllString(content, -1)
	analysis.EquationCount = len(equations)
}

func (ca *ContentAnalyzer) analyzeSentiment(content string, analysis *ContentAnalysis) {
	contentLower := strings.ToLower(content)
	
	positiveWords := []string{
		"happy", "joy", "love", "excellent", "wonderful", "beautiful",
		"feliz", "alegria", "amor", "excelente", "maravilhoso", "belo",
		"great", "fantastic", "amazing", "perfect", "brilliant",
	}
	
	negativeWords := []string{
		"sad", "fear", "hate", "terrible", "horrible", "awful",
		"triste", "medo", "ódio", "terrível", "horrível", "péssimo",
		"bad", "worst", "angry", "pain", "death", "dark",
	}
	
	positiveCount := 0
	negativeCount := 0
	
	for _, word := range positiveWords {
		positiveCount += strings.Count(contentLower, word)
	}
	
	for _, word := range negativeWords {
		negativeCount += strings.Count(contentLower, word)
	}
	
	total := positiveCount + negativeCount
	if total > 0 {
		analysis.SentimentScore = float64(positiveCount-negativeCount) / float64(total)
	} else {
		analysis.SentimentScore = 0.0
	}
}

func (ca *ContentAnalyzer) calculateRatios(analysis *ContentAnalysis) {
	estimatedPages := float64(analysis.WordCount) / 250.0
	if estimatedPages > 0 {
		analysis.ImageRatio = float64(analysis.ImageCount) / estimatedPages
	}
}

func (analysis *ContentAnalysis) IsAcademic() bool {
	return analysis.Tone.Academic > 0.6 || 
	       (analysis.Formality > 0.7 && analysis.TechnicalDensity > 0.5)
}

func (analysis *ContentAnalysis) HasComplexMath() bool {
	return analysis.EquationCount > 10
}

func (analysis *ContentAnalysis) HasRichMedia() bool {
	return analysis.ImageRatio > 0.1
}

func initGenreKeywords() map[string][]string {
	return map[string][]string{
		"fiction": {
			"chapter", "character", "story", "plot", "scene",
			"capítulo", "personagem", "história", "enredo", "cena",
			"dialogue", "narrator", "protagonist",
		},
		"mystery": {
			"detective", "murder", "crime", "investigation", "clue",
			"detetive", "assassinato", "crime", "investigação", "pista",
			"suspect", "evidence", "victim",
		},
		"romance": {
			"love", "heart", "kiss", "relationship", "passion",
			"amor", "coração", "beijo", "relacionamento", "paixão",
			"romance", "desire", "embrace",
		},
		"scifi": {
			"space", "alien", "technology", "future", "robot",
			"espaço", "alienígena", "tecnologia", "futuro", "robô",
			"planet", "galaxy", "cybernetic",
		},
		"fantasy": {
			"magic", "wizard", "dragon", "quest", "kingdom",
			"magia", "mago", "dragão", "missão", "reino",
			"spell", "enchanted", "mythical",
		},
		"technical": {
			"algorithm", "function", "system", "analysis", "method",
			"algoritmo", "função", "sistema", "análise", "método",
			"implementation", "architecture", "framework",
		},
		"academic": {
			"research", "study", "hypothesis", "methodology", "conclusion",
			"pesquisa", "estudo", "hipótese", "metodologia", "conclusão",
			"thesis", "analysis", "experiment",
		},
		"business": {
			"market", "strategy", "revenue", "customer", "growth",
			"mercado", "estratégia", "receita", "cliente", "crescimento",
			"investment", "profit", "management",
		},
	}
}

func initTechnicalTerms() map[string]bool {
	terms := []string{
		"algorithm", "data", "system", "process", "function",
		"method", "analysis", "implementation", "architecture",
		"framework", "protocol", "interface", "database",
		"api", "server", "client", "network", "security",
		"algoritmo", "dados", "sistema", "processo", "função",
		"método", "análise", "implementação", "arquitetura",
	}
	
	termMap := make(map[string]bool)
	for _, term := range terms {
		termMap[term] = true
	}
	return termMap
}

func initFormalWords() []string {
	return []string{
		"furthermore", "moreover", "consequently", "therefore",
		"however", "nevertheless", "additionally", "subsequently",
		"portanto", "entretanto", "todavia", "contudo",
		"ademais", "outrossim", "destarte",
	}
}

func initCasualWords() []string {
	return []string{
		"gonna", "wanna", "yeah", "okay", "cool",
		"stuff", "things", "guy", "hey",
		"vai", "tá", "né", "cara", "coisa",
	}
}
