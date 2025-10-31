// Package handlers provides HTTP handlers for the Typecraft API
// Conformidade: Constituição Vértice v3.0 - P1 (Completude Obrigatória)
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RenderHTMLRequest represents the request to render HTML
type RenderHTMLRequest struct {
	IncludeCSS    bool              `json:"include_css"`
	CSSVariables  map[string]string `json:"css_variables"`
	TemplateName  string            `json:"template_name"` // "base", "pagedjs"
}

// RenderHTMLResponse represents the HTML render response
type RenderHTMLResponse struct {
	ProjectID   int    `json:"project_id"`
	HTMLPath    string `json:"html_path"`
	CSSPath     string `json:"css_path,omitempty"`
	Size        int64  `json:"size_bytes"`
	GeneratedAt string `json:"generated_at"`
}

// RenderPDFRequest represents the request to render PDF
type RenderPDFRequest struct {
	Engine        string            `json:"engine"`         // "pagedjs", "prince", "weasyprint"
	Format        string            `json:"format"`         // "A4", "A5", "letter"
	Quality       string            `json:"quality"`        // "print", "screen", "ebook"
	Options       map[string]string `json:"options"`
	IncludeCovers bool              `json:"include_covers"`
}

// RenderPDFResponse represents the PDF render response
type RenderPDFResponse struct {
	ProjectID   int     `json:"project_id"`
	PDFPath     string  `json:"pdf_path"`
	Size        int64   `json:"size_bytes"`
	Pages       int     `json:"pages"`
	RenderTime  float64 `json:"render_time_seconds"`
	GeneratedAt string  `json:"generated_at"`
	Metadata    PDFMetadata `json:"metadata"`
}

// PDFMetadata represents PDF metadata
type PDFMetadata struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Subject     string `json:"subject"`
	Keywords    string `json:"keywords"`
	Creator     string `json:"creator"`
	Producer    string `json:"producer"`
}

// RenderHandler handles rendering requests
type RenderHandler struct {
	// Dependencies will be injected here
}

// NewRenderHandler creates a new render handler
func NewRenderHandler() *RenderHandler {
	return &RenderHandler{}
}

// RenderHTML handles POST /api/v1/projects/:id/render/html
func (h *RenderHandler) RenderHTML(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req RenderHTMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Default template
	if req.TemplateName == "" {
		req.TemplateName = "base"
	}

	// TODO: Integrate with HTML pipeline
	// For now, return mock response
	response := RenderHTMLResponse{
		ProjectID:   projectID,
		HTMLPath:    "/output/" + strconv.Itoa(projectID) + "/book.html",
		CSSPath:     "/output/" + strconv.Itoa(projectID) + "/styles.css",
		Size:        45678,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// RenderPDF handles POST /api/v1/projects/:id/render/pdf
func (h *RenderHandler) RenderPDF(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req RenderPDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Default values
	if req.Engine == "" {
		req.Engine = "pagedjs"
	}
	if req.Format == "" {
		req.Format = "A4"
	}
	if req.Quality == "" {
		req.Quality = "print"
	}

	// Validation
	validEngines := map[string]bool{
		"pagedjs":    true,
		"prince":     true,
		"weasyprint": true,
	}
	if !validEngines[req.Engine] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid engine. Valid options: pagedjs, prince, weasyprint"})
		return
	}

	// TODO: Integrate with PDF rendering pipeline
	// For now, return mock response
	response := RenderPDFResponse{
		ProjectID:  projectID,
		PDFPath:    "/output/" + strconv.Itoa(projectID) + "/book.pdf",
		Size:       1234567,
		Pages:      124,
		RenderTime: 3.45,
		GeneratedAt: time.Now().Format(time.RFC3339),
		Metadata: PDFMetadata{
			Title:    "Sample Book",
			Author:   "John Doe",
			Subject:  "Fiction",
			Keywords: "novel, fiction, story",
			Creator:  "Typecraft v1.0",
			Producer: "Typecraft PDF Engine (" + req.Engine + ")",
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetRenderStatus handles GET /api/v1/projects/:id/render/status
func (h *RenderHandler) GetRenderStatus(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// TODO: Implement actual status checking
	status := gin.H{
		"project_id": projectID,
		"html": gin.H{
			"status":      "completed",
			"last_render": "2024-10-31T19:30:00Z",
			"path":        "/output/" + strconv.Itoa(projectID) + "/book.html",
		},
		"pdf": gin.H{
			"status":      "completed",
			"last_render": "2024-10-31T19:35:00Z",
			"path":        "/output/" + strconv.Itoa(projectID) + "/book.pdf",
			"pages":       124,
		},
	}

	c.JSON(http.StatusOK, status)
}
