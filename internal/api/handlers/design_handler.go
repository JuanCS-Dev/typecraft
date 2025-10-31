// Package handlers provides HTTP handlers for the Typecraft API
// Conformidade: Constituição Vértice v3.0 - P1 (Completude Obrigatória)
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DesignGenerateRequest represents the request to generate design
type DesignGenerateRequest struct {
	Genre       string   `json:"genre"`        // Book genre (e.g., "fiction", "technical")
	Keywords    []string `json:"keywords"`     // Keywords for analysis
	Tone        string   `json:"tone"`         // Desired tone (e.g., "professional", "playful")
	ColorScheme string   `json:"color_scheme"` // Preferred color scheme (optional)
}

// DesignGenerateResponse represents the generated design
type DesignGenerateResponse struct {
	ProjectID   int               `json:"project_id"`
	ColorPalette ColorPalette     `json:"color_palette"`
	FontPairing  FontPairing      `json:"font_pairing"`
	GeneratedAt  string           `json:"generated_at"`
}

// ColorPalette represents the color palette
type ColorPalette struct {
	Primary     string   `json:"primary"`
	Secondary   string   `json:"secondary"`
	Accent      string   `json:"accent"`
	Background  string   `json:"background"`
	Text        string   `json:"text"`
	AllColors   []string `json:"all_colors"`
	Sentiment   string   `json:"sentiment"`
}

// FontPairing represents the font pairing suggestion
type FontPairing struct {
	Heading string `json:"heading"`
	Body    string `json:"body"`
	Code    string `json:"code,omitempty"`
	Rationale string `json:"rationale"`
}

// DesignHandler handles design-related requests
type DesignHandler struct {
	// Dependencies will be injected here
}

// NewDesignHandler creates a new design handler
func NewDesignHandler() *DesignHandler {
	return &DesignHandler{}
}

// GenerateDesign handles POST /api/v1/projects/:id/design/generate
func (h *DesignHandler) GenerateDesign(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req DesignGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validation
	if req.Genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Genre is required"})
		return
	}

	// TODO: Integrate with AI design modules
	// For now, return mock response
	response := DesignGenerateResponse{
		ProjectID: projectID,
		ColorPalette: ColorPalette{
			Primary:    "#2C3E50",
			Secondary:  "#34495E",
			Accent:     "#E74C3C",
			Background: "#ECF0F1",
			Text:       "#2C3E50",
			AllColors:  []string{"#2C3E50", "#34495E", "#E74C3C", "#ECF0F1"},
			Sentiment:  "professional",
		},
		FontPairing: FontPairing{
			Heading:   "Playfair Display",
			Body:      "Source Serif Pro",
			Code:      "JetBrains Mono",
			Rationale: "Classic serif pairing for professional technical content",
		},
		GeneratedAt: "2024-10-31T20:00:00Z",
	}

	c.JSON(http.StatusOK, response)
}

// ListFonts handles GET /api/v1/fonts
func (h *DesignHandler) ListFonts(c *gin.Context) {
	category := c.Query("category") // serif, sans-serif, monospace

	fonts := getFontDatabase()
	
	if category != "" {
		filtered := make([]Font, 0)
		for _, f := range fonts {
			if f.Category == category {
				filtered = append(filtered, f)
			}
		}
		fonts = filtered
	}

	c.JSON(http.StatusOK, gin.H{
		"fonts": fonts,
		"count": len(fonts),
	})
}

// Font represents a font in the database
type Font struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Source      string   `json:"source"`
	Variants    []string `json:"variants"`
	UseCase     []string `json:"use_case"`
	GoogleFonts bool     `json:"google_fonts"`
}

// getFontDatabase returns the available fonts
func getFontDatabase() []Font {
	return []Font{
		{
			Name:        "Playfair Display",
			Category:    "serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "bold", "bold-italic"},
			UseCase:     []string{"heading", "display"},
			GoogleFonts: true,
		},
		{
			Name:        "Source Serif Pro",
			Category:    "serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "semibold", "bold"},
			UseCase:     []string{"body", "text"},
			GoogleFonts: true,
		},
		{
			Name:        "Libre Baskerville",
			Category:    "serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "bold"},
			UseCase:     []string{"body", "text"},
			GoogleFonts: true,
		},
		{
			Name:        "Crimson Text",
			Category:    "serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "semibold", "bold"},
			UseCase:     []string{"body", "text"},
			GoogleFonts: true,
		},
		{
			Name:        "Inter",
			Category:    "sans-serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "medium", "semibold", "bold"},
			UseCase:     []string{"heading", "body", "ui"},
			GoogleFonts: true,
		},
		{
			Name:        "Roboto",
			Category:    "sans-serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "medium", "bold"},
			UseCase:     []string{"body", "ui"},
			GoogleFonts: true,
		},
		{
			Name:        "Open Sans",
			Category:    "sans-serif",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "semibold", "bold"},
			UseCase:     []string{"body", "ui"},
			GoogleFonts: true,
		},
		{
			Name:        "JetBrains Mono",
			Category:    "monospace",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "medium", "bold"},
			UseCase:     []string{"code"},
			GoogleFonts: true,
		},
		{
			Name:        "Fira Code",
			Category:    "monospace",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "medium", "semibold", "bold"},
			UseCase:     []string{"code"},
			GoogleFonts: true,
		},
		{
			Name:        "Source Code Pro",
			Category:    "monospace",
			Source:      "Google Fonts",
			Variants:    []string{"regular", "italic", "semibold", "bold"},
			UseCase:     []string{"code"},
			GoogleFonts: true,
		},
	}
}
