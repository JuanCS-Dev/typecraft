package design

import (
	"fmt"
	"math"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
)

// Color representa uma cor em RGB
type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	Hex string `json:"hex"`
}

// ColorPalette representa uma paleta de cores completa
type ColorPalette struct {
	Primary    Color           `json:"primary"`
	Secondary  Color           `json:"secondary"`
	Accent     Color           `json:"accent"`
	Background Color           `json:"background"`
	Text       Color           `json:"text"`
	Metadata   PaletteMetadata `json:"metadata"`
}

// PaletteMetadata contém informações sobre a paleta
type PaletteMetadata struct {
	Mood        string  `json:"mood"`
	Harmony     string  `json:"harmony"`
	Contrast    float64 `json:"contrast"`
	Accessibility string `json:"accessibility"`
	Rationale   string  `json:"rationale"`
}

// HSL representa uma cor em HSL (Hue, Saturation, Lightness)
type HSL struct {
	H float64 // 0-360
	S float64 // 0-1
	L float64 // 0-1
}

// ColorGenerator gera paletas de cores baseadas em análise de conteúdo
type ColorGenerator struct {
	// Mapping de sentimentos para hues
	sentimentHueMap map[string]float64
}

// NewColorGenerator cria um novo gerador de paletas
func NewColorGenerator() *ColorGenerator {
	return &ColorGenerator{
		sentimentHueMap: initSentimentHueMap(),
	}
}

// Generate gera uma paleta de cores contextual
func (cg *ColorGenerator) Generate(analysis *analyzer.ContentAnalysis) (*ColorPalette, error) {
	if analysis == nil {
		return nil, fmt.Errorf("analysis cannot be nil")
	}
	
	// 1. Determinar hue primária baseada em gênero e sentimento
	primaryHue := cg.selectPrimaryHue(analysis)
	
	// 2. Determinar tipo de harmonia
	harmonyType := cg.selectHarmony(analysis)
	
	// 3. Gerar cores
	primary := cg.generatePrimaryColor(primaryHue, analysis)
	secondary := cg.generateSecondaryColor(primaryHue, harmonyType, analysis)
	accent := cg.generateAccentColor(primaryHue, harmonyType, analysis)
	
	// 4. Cores neutras
	background := Color{R: 253, G: 254, B: 254, Hex: "#FDFEFE"}
	text := Color{R: 28, G: 28, B: 28, Hex: "#1C1C1C"}
	
	// 5. Ajustar para contraste se necessário
	if analysis.Tone.Academic > 0.7 {
		text = Color{R: 0, G: 0, B: 0, Hex: "#000000"}
		background = Color{R: 255, G: 255, B: 255, Hex: "#FFFFFF"}
	}
	
	// 6. Calcular contraste
	contrast := cg.calculateContrast(text, background)
	
	// 7. Criar paleta
	palette := &ColorPalette{
		Primary:    primary,
		Secondary:  secondary,
		Accent:     accent,
		Background: background,
		Text:       text,
		Metadata: PaletteMetadata{
			Mood:          cg.determineMood(analysis),
			Harmony:       harmonyType,
			Contrast:      contrast,
			Accessibility: cg.getAccessibilityLevel(contrast),
			Rationale:     cg.generateRationale(analysis, harmonyType),
		},
	}
	
	return palette, nil
}

// selectPrimaryHue seleciona a hue primária baseada no conteúdo
func (cg *ColorGenerator) selectPrimaryHue(analysis *analyzer.ContentAnalysis) float64 {
	genre := analysis.PrimaryGenre
	sentiment := analysis.SentimentScore
	
	// Base hue por gênero
	genreHues := map[string]float64{
		"fiction":   220.0, // Blue
		"mystery":   270.0, // Purple
		"romance":   340.0, // Pink/Red
		"scifi":     200.0, // Cyan
		"fantasy":   280.0, // Purple/Magenta
		"technical": 210.0, // Blue
		"academic":  220.0, // Blue
		"business":  210.0, // Blue
	}
	
	baseHue := genreHues[genre]
	if baseHue == 0 {
		baseHue = 220.0 // Default blue
	}
	
	// Ajustar baseado em sentimento
	if sentiment > 0.3 {
		// Positive: shift para warm
		baseHue = 45.0 // Orange/Yellow
	} else if sentiment < -0.3 {
		// Negative: shift para cool/dark
		baseHue = 240.0 // Deep blue
	}
	
	// Ajustar para gêneros específicos com sentimento
	if genre == "romance" && sentiment > 0 {
		baseHue = 350.0 // Romantic pink
	}
	if genre == "mystery" && sentiment < 0 {
		baseHue = 260.0 // Dark purple
	}
	
	return baseHue
}

// selectHarmony seleciona o tipo de harmonia de cores
func (cg *ColorGenerator) selectHarmony(analysis *analyzer.ContentAnalysis) string {
	if analysis.Tone.Academic > 0.7 {
		return "analogous"
	}
	if analysis.Tone.Creative > 0.7 {
		return "triadic"
	}
	if analysis.TechnicalDensity > 0.5 {
		return "complementary"
	}
	return "analogous"
}

// generatePrimaryColor gera a cor primária
func (cg *ColorGenerator) generatePrimaryColor(hue float64, analysis *analyzer.ContentAnalysis) Color {
	saturation := 0.6
	lightness := 0.4
	
	// Ajustar saturação baseado em formalidade
	if analysis.Formality > 0.7 {
		saturation = 0.4
	} else if analysis.Formality < 0.3 {
		saturation = 0.8
	}
	
	// Ajustar lightness baseado em complexidade
	if analysis.Complexity > 0.7 {
		lightness = 0.35
	}
	
	return hslToRGB(HSL{H: hue, S: saturation, L: lightness})
}

// generateSecondaryColor gera a cor secundária
func (cg *ColorGenerator) generateSecondaryColor(baseHue float64, harmony string, analysis *analyzer.ContentAnalysis) Color {
	var hue float64
	
	switch harmony {
	case "complementary":
		hue = math.Mod(baseHue+180.0, 360.0)
	case "analogous":
		hue = math.Mod(baseHue+30.0, 360.0)
	case "triadic":
		hue = math.Mod(baseHue+120.0, 360.0)
	default:
		hue = math.Mod(baseHue+30.0, 360.0)
	}
	
	saturation := 0.5
	lightness := 0.5
	
	return hslToRGB(HSL{H: hue, S: saturation, L: lightness})
}

// generateAccentColor gera a cor de acento
func (cg *ColorGenerator) generateAccentColor(baseHue float64, harmony string, analysis *analyzer.ContentAnalysis) Color {
	var hue float64
	
	switch harmony {
	case "complementary":
		hue = math.Mod(baseHue+200.0, 360.0)
	case "analogous":
		hue = math.Mod(baseHue-30.0, 360.0)
	case "triadic":
		hue = math.Mod(baseHue+240.0, 360.0)
	default:
		hue = math.Mod(baseHue+60.0, 360.0)
	}
	
	saturation := 0.7
	lightness := 0.55
	
	// Accent mais vibrante para conteúdo criativo
	if analysis.Tone.Creative > 0.6 {
		saturation = 0.85
	}
	
	return hslToRGB(HSL{H: hue, S: saturation, L: lightness})
}

// calculateContrast calcula o contrast ratio entre duas cores
func (cg *ColorGenerator) calculateContrast(c1, c2 Color) float64 {
	l1 := cg.relativeLuminance(c1)
	l2 := cg.relativeLuminance(c2)
	
	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)
	
	return (lighter + 0.05) / (darker + 0.05)
}

// relativeLuminance calcula a luminância relativa
func (cg *ColorGenerator) relativeLuminance(c Color) float64 {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0
	
	// Convert to linear
	if r <= 0.03928 {
		r = r / 12.92
	} else {
		r = math.Pow((r+0.055)/1.055, 2.4)
	}
	
	if g <= 0.03928 {
		g = g / 12.92
	} else {
		g = math.Pow((g+0.055)/1.055, 2.4)
	}
	
	if b <= 0.03928 {
		b = b / 12.92
	} else {
		b = math.Pow((b+0.055)/1.055, 2.4)
	}
	
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// getAccessibilityLevel retorna o nível de acessibilidade WCAG
func (cg *ColorGenerator) getAccessibilityLevel(contrast float64) string {
	if contrast >= 7.0 {
		return "AAA"
	}
	if contrast >= 4.5 {
		return "AA"
	}
	if contrast >= 3.0 {
		return "AA Large"
	}
	return "Fail"
}

// determineMood determina o mood da paleta
func (cg *ColorGenerator) determineMood(analysis *analyzer.ContentAnalysis) string {
	if analysis.SentimentScore > 0.3 {
		return "bright"
	}
	if analysis.SentimentScore < -0.3 {
		return "dark"
	}
	if analysis.Tone.Formal > 0.7 {
		return "professional"
	}
	if analysis.Tone.Creative > 0.7 {
		return "vibrant"
	}
	return "balanced"
}

// generateRationale gera explicação da paleta
func (cg *ColorGenerator) generateRationale(analysis *analyzer.ContentAnalysis, harmony string) string {
	mood := cg.determineMood(analysis)
	genre := analysis.PrimaryGenre
	
	rationales := map[string]map[string]string{
		"fiction": {
			"bright":       "Warm, inviting colors for engaging storytelling",
			"dark":         "Cool, subdued tones for contemplative narrative",
			"professional": "Balanced palette for serious literary work",
			"vibrant":      "Dynamic colors reflecting creative expression",
			"balanced":     "Harmonious colors for traditional fiction",
		},
		"mystery": {
			"dark":         "Deep, mysterious tones for suspenseful atmosphere",
			"professional": "Sophisticated colors for detective stories",
			"balanced":     "Intriguing palette for mystery novels",
		},
		"romance": {
			"bright":  "Soft, romantic colors for love stories",
			"vibrant": "Passionate tones for contemporary romance",
		},
		"technical": {
			"professional": "Clean, professional palette for technical content",
			"balanced":     "Clear colors for technical documentation",
		},
		"academic": {
			"professional": "Scholarly palette for academic work",
			"balanced":     "Traditional colors for research publication",
		},
	}
	
	if genreMap, ok := rationales[genre]; ok {
		if rationale, ok := genreMap[mood]; ok {
			return rationale
		}
	}
	
	return fmt.Sprintf("%s %s palette for %s content", harmony, mood, genre)
}

// hslToRGB converte HSL para RGB
func hslToRGB(hsl HSL) Color {
	h := hsl.H / 360.0
	s := hsl.S
	l := hsl.L
	
	var r, g, b float64
	
	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}
	
	rByte := uint8(math.Round(r * 255))
	gByte := uint8(math.Round(g * 255))
	bByte := uint8(math.Round(b * 255))
	
	hex := fmt.Sprintf("#%02X%02X%02X", rByte, gByte, bByte)
	
	return Color{
		R:   rByte,
		G:   gByte,
		B:   bByte,
		Hex: hex,
	}
}

// hueToRGB helper para conversão HSL->RGB
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// initSentimentHueMap inicializa mapping de sentimentos
func initSentimentHueMap() map[string]float64 {
	return map[string]float64{
		"joy":     50.0,  // Yellow
		"love":    350.0, // Pink
		"calm":    200.0, // Blue-green
		"sadness": 240.0, // Blue
		"fear":    270.0, // Purple
		"anger":   0.0,   // Red
	}
}
