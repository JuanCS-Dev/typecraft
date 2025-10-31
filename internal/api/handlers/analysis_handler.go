// Package handlers provides HTTP handlers for the API
// Following Constituição Vértice v3.0 - Artigo I
package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JuanCS-Dev/typecraft/internal/ai"
	"github.com/JuanCS-Dev/typecraft/internal/repository"
	"github.com/JuanCS-Dev/typecraft/internal/service"
	"github.com/gin-gonic/gin"
)

// AnalysisHandler handles manuscript analysis requests
type AnalysisHandler struct {
	analysisService *service.AnalysisService
	projectService  *service.ProjectService
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisService *service.AnalysisService, projectService *service.ProjectService) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
		projectService:  projectService,
	}
}

// NewAnalysisHandlerWithDeps creates a new analysis handler with all dependencies initialized
func NewAnalysisHandlerWithDeps() (*AnalysisHandler, error) {
	// Get OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	
	// Create analyzer with API key and model
	model := getEnvOrDefault("OPENAI_MODEL", "gpt-4o")
	analyzer := ai.NewAnalyzer(apiKey, model)
	
	// Create repositories (they get DB from internal package)
	projectRepo := repository.NewProjectRepository()
	
	// Create services
	analysisService := service.NewAnalysisService(analyzer, projectRepo)
	projectService := service.NewProjectService()
	
	return NewAnalysisHandler(analysisService, projectService), nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		if intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}

func getFloat32Env(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		var floatValue float32
		fmt.Sscanf(value, "%f", &floatValue)
		if floatValue > 0 {
			return floatValue
		}
	}
	return defaultValue
}

// AnalyzeProjectRequest represents the request body for project analysis
type AnalyzeProjectRequest struct {
	ForceReanalysis         bool `json:"force_reanalysis"`
	IncludeRecommendations  bool `json:"include_recommendations"`
}

// AnalyzeProjectResponse represents the response for project analysis
type AnalyzeProjectResponse struct {
	Analysis       interface{} `json:"analysis"`
	Recommendations interface{} `json:"recommendations,omitempty"`
}

// AnalyzeProject analyzes a project's manuscript content
// POST /api/v1/projects/:id/analyze
func (h *AnalysisHandler) AnalyzeProject(c *gin.Context) {
	projectID := c.Param("id")
	
	var req AnalyzeProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Use defaults if no body provided
		req.IncludeRecommendations = true
	}
	
	// Get project
	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	
	// Get manuscript content (placeholder - needs actual implementation)
	manuscriptText := project.Description // For now, analyze description
	if manuscriptText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project has no content to analyze"})
		return
	}
	
	// If force_reanalysis is true, we could invalidate cache here
	// For now, the service will handle caching automatically
	
	// Perform analysis
	analysis, err := h.analysisService.AnalyzeProject(c.Request.Context(), projectID, manuscriptText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "analysis failed", "details": err.Error()})
		return
	}
	
	response := AnalyzeProjectResponse{
		Analysis: analysis,
	}
	
	// Include typography recommendations if requested
	if req.IncludeRecommendations {
		recommendations := h.analysisService.GetTypographicRecommendations(analysis)
		response.Recommendations = recommendations
	}
	
	c.JSON(http.StatusOK, response)
}

// ListGenres lists all supported genres
// GET /api/v1/genres
func (h *AnalysisHandler) ListGenres(c *gin.Context) {
	genres := []map[string]interface{}{
		{
			"id":   "fiction",
			"name": "Fiction",
			"sub_genres": []string{"literary", "commercial", "speculative"},
		},
		{
			"id":   "nonfiction",
			"name": "Non-Fiction",
			"sub_genres": []string{"memoir", "biography", "self-help", "history"},
		},
		{
			"id":   "technical",
			"name": "Technical",
			"sub_genres": []string{"programming", "engineering", "science"},
		},
		{
			"id":   "academic",
			"name": "Academic",
			"sub_genres": []string{"textbook", "research", "monograph"},
		},
		{
			"id":   "poetry",
			"name": "Poetry",
			"sub_genres": []string{"free_verse", "sonnet", "haiku"},
		},
		{
			"id":   "childrens",
			"name": "Children's",
			"sub_genres": []string{"picture_book", "middle_grade", "young_adult"},
		},
		{
			"id":   "self_help",
			"name": "Self-Help",
			"sub_genres": []string{"motivational", "wellness", "business"},
		},
		{
			"id":   "biography",
			"name": "Biography",
			"sub_genres": []string{"autobiography", "memoir", "historical"},
		},
		{
			"id":   "historical",
			"name": "Historical",
			"sub_genres": []string{"ancient", "medieval", "modern"},
		},
		{
			"id":   "sci_fi",
			"name": "Science Fiction",
			"sub_genres": []string{"hard_sf", "space_opera", "cyberpunk"},
		},
		{
			"id":   "fantasy",
			"name": "Fantasy",
			"sub_genres": []string{"epic", "urban", "dark"},
		},
		{
			"id":   "romance",
			"name": "Romance",
			"sub_genres": []string{"contemporary", "historical", "paranormal"},
		},
		{
			"id":   "mystery",
			"name": "Mystery",
			"sub_genres": []string{"cozy", "hardboiled", "police_procedural"},
		},
		{
			"id":   "cookbook",
			"name": "Cookbook",
			"sub_genres": []string{"regional", "diet", "baking"},
		},
		{
			"id":   "travel",
			"name": "Travel",
			"sub_genres": []string{"guide", "memoir", "adventure"},
		},
		{
			"id":   "art",
			"name": "Art & Photography",
			"sub_genres": []string{"photography", "design", "architecture"},
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"genres": genres,
		"total":  len(genres),
	})
}
