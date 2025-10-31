// Package handlers implements HTTP handlers for the book generation API.
// This exposes the Orchestrator through REST endpoints (Dia 7, Sprint 7-8).
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/JuanCS-Dev/typecraft/internal/service"
)

// BookGenerationHandler handles book generation endpoints
type BookGenerationHandler struct {
	orchestrator *service.BookOrchestrator
}

// NewBookGenerationHandler creates a new handler
func NewBookGenerationHandler(orchestrator *service.BookOrchestrator) *BookGenerationHandler {
	return &BookGenerationHandler{
		orchestrator: orchestrator,
	}
}

// GenerateBookRequest is the HTTP request body
type GenerateBookRequest struct {
	ContentPath      string                   `json:"content_path" binding:"required"`
	OutputFormats    []string                 `json:"output_formats" binding:"required,dive,oneof=pdf epub"`
	OverridePipeline string                   `json:"override_pipeline,omitempty" binding:"omitempty,oneof=latex html"`
	CustomDesign     *CustomDesignRequest     `json:"custom_design,omitempty"`
}

// CustomDesignRequest allows custom design parameters
type CustomDesignRequest struct {
	BodyFont      string   `json:"body_font,omitempty"`
	HeadingFont   string   `json:"heading_font,omitempty"`
	ColorScheme   []string `json:"color_scheme,omitempty"`
	MarginPreset  string   `json:"margin_preset,omitempty"`
}

// GenerateBookResponse is the HTTP response
type GenerateBookResponse struct {
	Success        bool                  `json:"success"`
	Message        string                `json:"message,omitempty"`
	ProjectID      uint                  `json:"project_id"`
	Pipeline       string                `json:"pipeline"`
	OutputFiles    map[string]string     `json:"output_files"`
	DesignMetadata *DesignMetadataResponse `json:"design_metadata,omitempty"`
	Metrics        *MetricsResponse      `json:"metrics,omitempty"`
	Error          string                `json:"error,omitempty"`
}

// DesignMetadataResponse contains design information
type DesignMetadataResponse struct {
	Fonts   FontsResponse   `json:"fonts"`
	Colors  []string        `json:"colors"`
	Margins MarginsResponse `json:"margins"`
}

// FontsResponse contains font information
type FontsResponse struct {
	Body    string `json:"body"`
	Heading string `json:"heading"`
}

// MarginsResponse contains margin information
type MarginsResponse struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

// MetricsResponse contains generation metrics
type MetricsResponse struct {
	DurationMs         int64   `json:"duration_ms"`
	ContentAnalysisMs  int64   `json:"content_analysis_ms"`
	DesignGenerationMs int64   `json:"design_generation_ms"`
	RenderingMs        int64   `json:"rendering_ms"`
	ValidationMs       int64   `json:"validation_ms"`
	TotalPages         int     `json:"total_pages"`
	FileSize           int64   `json:"file_size"`
	QualityScore       float64 `json:"quality_score"`
}

// Generate handles POST /api/v1/projects/:id/generate
// @Summary Generate book in multiple formats
// @Description Orchestrates the complete book generation pipeline
// @Tags generation
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body GenerateBookRequest true "Generation parameters"
// @Success 200 {object} GenerateBookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/projects/{id}/generate [post]
func (h *BookGenerationHandler) Generate(c *gin.Context) {
	// Parse project ID
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid project ID",
			Message: err.Error(),
		})
		return
	}

	// Parse request body
	var req GenerateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Convert to service request
	serviceReq := &service.GenerationRequest{
		ProjectID:        uint(projectID),
		ContentPath:      req.ContentPath,
		OutputFormats:    req.OutputFormats,
		OverridePipeline: req.OverridePipeline,
	}

	// Apply custom design if provided
	if req.CustomDesign != nil {
		serviceReq.CustomDesign = &service.DesignOptions{
			BodyFont:     req.CustomDesign.BodyFont,
			HeadingFont:  req.CustomDesign.HeadingFont,
			ColorScheme:  req.CustomDesign.ColorScheme,
			MarginPreset: req.CustomDesign.MarginPreset,
		}
	}

	// Execute generation
	result, err := h.orchestrator.Generate(c.Request.Context(), serviceReq)
	
	// Build response
	response := &GenerateBookResponse{
		Success:     result.Success,
		ProjectID:   result.ProjectID,
		Pipeline:    result.Pipeline,
		OutputFiles: result.OutputFiles,
	}

	// Add design metadata if available
	if result.DesignMetadata != nil {
		response.DesignMetadata = &DesignMetadataResponse{
			Fonts: FontsResponse{
				Body:    result.DesignMetadata.Fonts.Body,
				Heading: result.DesignMetadata.Fonts.Heading,
			},
			Colors: result.DesignMetadata.Colors,
			Margins: MarginsResponse{
				Top:    result.DesignMetadata.Margins.Top,
				Bottom: result.DesignMetadata.Margins.Bottom,
				Left:   result.DesignMetadata.Margins.Left,
				Right:  result.DesignMetadata.Margins.Right,
			},
		}
	}

	// Add metrics if available
	if result.Metrics != nil {
		response.Metrics = &MetricsResponse{
			DurationMs:         result.Metrics.Duration.Milliseconds(),
			ContentAnalysisMs:  result.Metrics.ContentAnalysisMs,
			DesignGenerationMs: result.Metrics.DesignGenerationMs,
			RenderingMs:        result.Metrics.RenderingMs,
			ValidationMs:       result.Metrics.ValidationMs,
			TotalPages:         result.Metrics.TotalPages,
			FileSize:           result.Metrics.FileSize,
			QualityScore:       result.Metrics.QualityScore,
		}
	}

	// Handle errors
	if err != nil {
		response.Error = err.Error()
		response.Message = "Book generation failed"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Message = "Book generated successfully"
	c.JSON(http.StatusOK, response)
}

// GetProgress handles GET /api/v1/projects/:id/generation/progress
// @Summary Get generation progress
// @Description Returns the current progress of a book generation job
// @Tags generation
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/projects/{id}/generation/progress [get]
func (h *BookGenerationHandler) GetProgress(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid project ID",
			Message: err.Error(),
		})
		return
	}

	progress, err := h.orchestrator.GetProgress(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Progress not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ProgressResponse{
		ProjectID:    progress.ProjectID,
		Status:       progress.Status,
		CurrentStage: progress.CurrentStage,
		Progress:     progress.Progress,
		Message:      progress.Message,
	})
}

// CancelGeneration handles DELETE /api/v1/projects/:id/generation
// @Summary Cancel generation
// @Description Cancels an in-progress book generation job
// @Tags generation
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/projects/{id}/generation [delete]
func (h *BookGenerationHandler) CancelGeneration(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid project ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.orchestrator.CancelGeneration(c.Request.Context(), uint(projectID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to cancel generation",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: fmt.Sprintf("Generation cancelled for project %d", projectID),
	})
}

// ProgressResponse contains progress information
type ProgressResponse struct {
	ProjectID    uint   `json:"project_id"`
	Status       string `json:"status"`
	CurrentStage string `json:"current_stage"`
	Progress     int    `json:"progress"`
	Message      string `json:"message,omitempty"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// MessageResponse is a standard message response
type MessageResponse struct {
	Message string `json:"message"`
}

// RegisterRoutes registers all generation routes
func (h *BookGenerationHandler) RegisterRoutes(router *gin.RouterGroup) {
	generation := router.Group("/projects/:id/generation")
	{
		generation.POST("", h.Generate)
		generation.GET("/progress", h.GetProgress)
		generation.DELETE("", h.CancelGeneration)
	}
}
