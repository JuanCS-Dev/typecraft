package design

import (
	"context"
)

// Service orchestrates design generation
type Service struct {
	fontSuggester  *FontSuggester
	colorGenerator *ColorGenerator
}

// NewService creates a new design service
func NewService() *Service {
	fontSuggester, _ := NewFontSuggester()
	return &Service{
		fontSuggester:  fontSuggester,
		colorGenerator: NewColorGenerator(),
	}
}

// DesignRequest encapsulates design generation parameters
type DesignRequest struct {
	Genre         string
	Tone          string
	Complexity    float64
	PageFormat    string
	BodyFont      string
	HeadingFont   string
	ColorScheme   []string
	MarginPreset  string
	CustomMargins *Margins
}

// DesignResult contains the generated design
type DesignResult struct {
	Fonts   Fonts
	Colors  []string
	Margins Margins
}

// Fonts contains font selections
type Fonts struct {
	Body    string
	Heading string
}

// Margins contains page margins in millimeters
type Margins struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

// GenerateDesign generates a complete design based on content analysis
func (s *Service) GenerateDesign(ctx context.Context, req *DesignRequest) (*DesignResult, error) {
	result := &DesignResult{}

	// Generate fonts
	if req.BodyFont != "" && req.HeadingFont != "" {
		result.Fonts.Body = req.BodyFont
		result.Fonts.Heading = req.HeadingFont
	} else {
		// Use default fonts based on genre
		result.Fonts.Body = getDefaultBodyFont(req.Genre)
		result.Fonts.Heading = getDefaultHeadingFont(req.Genre)
	}

	// Generate colors
	if len(req.ColorScheme) > 0 {
		result.Colors = req.ColorScheme
	} else {
		result.Colors = getDefaultColors(req.Genre)
	}

	// Generate margins (Van de Graaf Canon)
	if req.CustomMargins != nil {
		result.Margins = *req.CustomMargins
	} else {
		result.Margins = generateVanDeGraafMargins(req.PageFormat)
	}

	return result, nil
}

// getDefaultBodyFont returns a default body font based on genre
func getDefaultBodyFont(genre string) string {
	fonts := map[string]string{
		"Fiction":   "Garamond",
		"Academic":  "Palatino",
		"Technical": "Computer Modern",
		"Poetry":    "Baskerville",
	}
	if font, ok := fonts[genre]; ok {
		return font
	}
	return "Garamond" // default
}

// getDefaultHeadingFont returns a default heading font based on genre
func getDefaultHeadingFont(genre string) string {
	fonts := map[string]string{
		"Fiction":   "Futura",
		"Academic":  "Helvetica",
		"Technical": "Computer Modern Sans",
		"Poetry":    "Optima",
	}
	if font, ok := fonts[genre]; ok {
		return font
	}
	return "Futura" // default
}

// getDefaultColors returns default colors based on genre
func getDefaultColors(genre string) []string {
	colors := map[string][]string{
		"Fiction":   {"#2C3E50", "#ECF0F1"},
		"Academic":  {"#34495E", "#FFFFFF"},
		"Technical": {"#1A1A1A", "#F5F5F5"},
		"Poetry":    {"#8E44AD", "#F4ECF7"},
	}
	if palette, ok := colors[genre]; ok {
		return palette
	}
	return []string{"#000000", "#FFFFFF"} // default
}

// generateVanDeGraafMargins generates classical book margins
func generateVanDeGraafMargins(pageFormat string) Margins {
	// Simplified Van de Graaf Canon
	// Inner margin : Top margin : Outer margin : Bottom margin = 2:3:4:6
	
	baseMargin := 20.0 // mm
	
	return Margins{
		Top:    baseMargin * 3 / 2,  // 30mm
		Bottom: baseMargin * 6 / 2,  // 60mm
		Left:   baseMargin,          // 20mm (inner)
		Right:  baseMargin * 4 / 2,  // 40mm (outer)
	}
}
