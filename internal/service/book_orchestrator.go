// Package service implements the Orchestrator - the central conductor of the book generation pipeline.
// This is the implementation of the INTEGRATION (Dia 7, Sprint 7-8) as per Constituição VÉRTICE v3.0.
package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"github.com/JuanCS-Dev/typecraft/pkg/design"
)

// BookOrchestrator coordinates the complete book generation workflow.
// Following VÉRTICE P5 (System Awareness): designed for evolution, not rewrites.
type BookOrchestrator struct {
	projectRepo    domain.ProjectRepository
	analysisClient AnalysisClient
	designService  *design.Service
	
	// Output configuration
	outputDir string
}

// AnalysisClient interface for AI content analysis
type AnalysisClient interface {
	AnalyzeContent(ctx context.Context, content string) (*domain.Analysis, error)
}

// NewBookOrchestrator creates a new orchestrator instance
func NewBookOrchestrator(
	projectRepo domain.ProjectRepository,
	analysisClient AnalysisClient,
	outputDir string,
) *BookOrchestrator {
	return &BookOrchestrator{
		projectRepo:    projectRepo,
		analysisClient: analysisClient,
		designService:  design.NewService(),
		outputDir:      outputDir,
	}
}

// GenerationRequest encapsulates all parameters for book generation
type GenerationRequest struct {
	ProjectID       uint
	ContentPath     string
	OutputFormats   []string // ["pdf", "epub"]
	OverridePipeline string   // "latex" or "html" (optional)
	CustomDesign    *DesignOptions
}

// DesignOptions allows custom design parameters
type DesignOptions struct {
	BodyFont        string
	HeadingFont     string
	ColorScheme     []string
	MarginPreset    string
	CustomMargins   *design.Margins
}

// GenerationResult contains all outputs and metadata
type GenerationResult struct {
	ProjectID      uint
	Pipeline       string
	OutputFiles    map[string]string // format -> filepath
	DesignMetadata *design.DesignResult
	Analysis       *domain.Analysis
	Metrics        *GenerationMetrics
	Success        bool
	Error          error
}

// GenerationMetrics tracks performance and quality
type GenerationMetrics struct {
	StartTime         time.Time
	EndTime           time.Time
	Duration          time.Duration
	ContentAnalysisMs int64
	DesignGenerationMs int64
	PipelineSelectionMs int64
	RenderingMs       int64
	ValidationMs      int64
	TotalPages        int
	FileSize          int64
	QualityScore      float64
}

// Generate orchestrates the complete book generation pipeline.
// This is the MAIN INTEGRATION POINT following VÉRTICE architecture.
func (o *BookOrchestrator) Generate(ctx context.Context, req *GenerationRequest) (*GenerationResult, error) {
	metrics := &GenerationMetrics{StartTime: time.Now()}
	result := &GenerationResult{
		ProjectID:   req.ProjectID,
		OutputFiles: make(map[string]string),
		Metrics:     metrics,
	}

	// STEP 1: Load project
	project, err := o.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		result.Error = fmt.Errorf("failed to load project: %w", err)
		return result, result.Error
	}

	// STEP 2: Read and validate content
	content, err := o.readContent(req.ContentPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to read content: %w", err)
		return result, result.Error
	}

	// STEP 3: AI Content Analysis
	analysisStart := time.Now()
	analysis, err := o.analysisClient.AnalyzeContent(ctx, content)
	if err != nil {
		result.Error = fmt.Errorf("content analysis failed: %w", err)
		return result, result.Error
	}
	result.Analysis = analysis
	metrics.ContentAnalysisMs = time.Since(analysisStart).Milliseconds()

	// STEP 4: Design Generation
	designStart := time.Now()
	designReq := o.buildDesignRequest(project, analysis, req.CustomDesign)
	designResult, err := o.designService.GenerateDesign(ctx, designReq)
	if err != nil {
		result.Error = fmt.Errorf("design generation failed: %w", err)
		return result, result.Error
	}
	result.DesignMetadata = designResult
	metrics.DesignGenerationMs = time.Since(designStart).Milliseconds()

	// STEP 5: Pipeline Selection
	selectionStart := time.Now()
	selectedPipeline := o.selectPipeline(analysis, req.OverridePipeline)
	result.Pipeline = selectedPipeline
	metrics.PipelineSelectionMs = time.Since(selectionStart).Milliseconds()

	// STEP 6: Rendering
	renderStart := time.Now()
	if err := o.renderOutputs(ctx, req, content, designResult, selectedPipeline, result); err != nil {
		result.Error = fmt.Errorf("rendering failed: %w", err)
		return result, result.Error
	}
	metrics.RenderingMs = time.Since(renderStart).Milliseconds()

	// STEP 7: Validation
	validationStart := time.Now()
	if err := o.validateOutputs(result); err != nil {
		result.Error = fmt.Errorf("validation failed: %w", err)
		return result, result.Error
	}
	metrics.ValidationMs = time.Since(validationStart).Milliseconds()

	// Finalize metrics
	metrics.EndTime = time.Now()
	metrics.Duration = metrics.EndTime.Sub(metrics.StartTime)
	result.Success = true

	return result, nil
}

// readContent reads and validates the input content
func (o *BookOrchestrator) readContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	content := string(data)
	if len(content) == 0 {
		return "", fmt.Errorf("content file is empty")
	}
	
	return content, nil
}

// buildDesignRequest creates design generation request from project data
func (o *BookOrchestrator) buildDesignRequest(
	project *domain.Project,
	analysis *domain.Analysis,
	customDesign *DesignOptions,
) *design.DesignRequest {
	req := &design.DesignRequest{
		Genre:      project.Genre,
		Tone:       analysis.Tone,
		Complexity: analysis.Complexity,
		PageFormat: project.PageFormat,
	}

	// Apply custom design if provided
	if customDesign != nil {
		if customDesign.BodyFont != "" {
			req.BodyFont = customDesign.BodyFont
		}
		if customDesign.HeadingFont != "" {
			req.HeadingFont = customDesign.HeadingFont
		}
		if len(customDesign.ColorScheme) > 0 {
			req.ColorScheme = customDesign.ColorScheme
		}
		if customDesign.CustomMargins != nil {
			req.CustomMargins = customDesign.CustomMargins
		}
	}

	return req
}

// selectPipeline determines which rendering pipeline to use
func (o *BookOrchestrator) selectPipeline(analysis *domain.Analysis, override string) string {
	if override != "" {
		return override
	}
	
	// Simple decision logic
	if analysis.HasMath || analysis.CodeBlocks > 5 {
		return "latex"
	}
	
	if analysis.ImageCount > 10 {
		return "html"
	}
	
	// Default to LaTeX for text-heavy content
	return "latex"
}

// renderOutputs generates all requested output formats
func (o *BookOrchestrator) renderOutputs(
	ctx context.Context,
	req *GenerationRequest,
	content string,
	design *design.DesignResult,
	selectedPipeline string,
	result *GenerationResult,
) error {
	for _, format := range req.OutputFormats {
		switch format {
		case "pdf":
			pdfPath, err := o.renderPDF(ctx, content, design, selectedPipeline, req.ProjectID)
			if err != nil {
				return fmt.Errorf("PDF rendering failed: %w", err)
			}
			result.OutputFiles["pdf"] = pdfPath

		case "epub":
			epubPath, err := o.renderEPUB(ctx, content, design, req.ProjectID)
			if err != nil {
				return fmt.Errorf("ePub rendering failed: %w", err)
			}
			result.OutputFiles["epub"] = epubPath

		default:
			return fmt.Errorf("unsupported output format: %s", format)
		}
	}

	return nil
}

// renderPDF generates PDF using the selected pipeline
func (o *BookOrchestrator) renderPDF(
	ctx context.Context,
	content string,
	design *design.DesignResult,
	pipelineType string,
	projectID uint,
) (string, error) {
	outputPath := filepath.Join(o.outputDir, fmt.Sprintf("project_%d.pdf", projectID))

	switch pipelineType {
	case "latex":
		return o.renderPDFLaTeX(content, design, outputPath)
	case "html":
		return o.renderPDFHTML(content, design, outputPath)
	default:
		return "", fmt.Errorf("unknown pipeline type: %s", pipelineType)
	}
}

// renderPDFLaTeX generates PDF via LaTeX pipeline
func (o *BookOrchestrator) renderPDFLaTeX(content string, design *design.DesignResult, outputPath string) (string, error) {
	// Simplified stub - demonstrates integration point
	// In real implementation, would call latex.Compiler
	
	// Create a placeholder PDF file
	placeholder := []byte("%PDF-1.4\n%Placeholder PDF generated by orchestrator\n%%EOF")
	if err := os.WriteFile(outputPath, placeholder, 0644); err != nil {
		return "", fmt.Errorf("failed to write PDF: %w", err)
	}

	return outputPath, nil
}

// renderPDFHTML generates PDF via HTML/CSS + Paged.js pipeline
func (o *BookOrchestrator) renderPDFHTML(content string, design *design.DesignResult, outputPath string) (string, error) {
	// Simplified stub - demonstrates integration point
	// In real implementation, would call pipeline.HTMLGenerator and PDFGenerator
	
	// Create a placeholder PDF file
	placeholder := []byte("%PDF-1.4\n%Placeholder PDF generated via HTML pipeline\n%%EOF")
	if err := os.WriteFile(outputPath, placeholder, 0644); err != nil {
		return "", fmt.Errorf("failed to write PDF: %w", err)
	}

	return outputPath, nil
}

// renderEPUB generates ePub output
func (o *BookOrchestrator) renderEPUB(
	ctx context.Context,
	content string,
	design *design.DesignResult,
	projectID uint,
) (string, error) {
	outputPath := filepath.Join(o.outputDir, fmt.Sprintf("project_%d.epub", projectID))

	// Simplified stub - demonstrates integration point
	// In real implementation, would call epub.Generator
	
	// Create a placeholder ZIP file (ePub is a ZIP)
	placeholder := []byte("PK\x03\x04") // ZIP magic bytes
	if err := os.WriteFile(outputPath, placeholder, 0644); err != nil {
		return "", fmt.Errorf("failed to write ePub: %w", err)
	}

	return outputPath, nil
}

// splitIntoChapters splits content into chapters (simplified version)
func (o *BookOrchestrator) splitIntoChapters(content string) []string {
	// For now, return as single chapter
	// TODO: Implement smart chapter detection
	return []string{content}
}

// validateOutputs performs final validation on all generated files
func (o *BookOrchestrator) validateOutputs(result *GenerationResult) error {
	for format, path := range result.OutputFiles {
		// Check file exists
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("%s file not found: %w", format, err)
		}

		// Check file size
		if info.Size() == 0 {
			return fmt.Errorf("%s file is empty", format)
		}

		// Update metrics
		if result.Metrics != nil {
			result.Metrics.FileSize += info.Size()
		}

		// Format-specific validation
		switch format {
		case "pdf":
			if err := o.validatePDF(path); err != nil {
				return fmt.Errorf("PDF validation failed: %w", err)
			}
		case "epub":
			if err := o.validateEPUBFile(path); err != nil {
				return fmt.Errorf("ePub validation failed: %w", err)
			}
		}
	}

	return nil
}

// validatePDF validates PDF file integrity
func (o *BookOrchestrator) validatePDF(path string) error {
	// Basic validation: check PDF magic bytes
	data := make([]byte, 4)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Read(data); err != nil {
		return err
	}

	if string(data) != "%PDF" {
		return fmt.Errorf("invalid PDF header")
	}

	return nil
}

// validateEPUBFile validates ePub file integrity
func (o *BookOrchestrator) validateEPUBFile(path string) error {
	// Basic validation: check ZIP magic bytes
	data := make([]byte, 4)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Read(data); err != nil {
		return err
	}

	// Check for ZIP signature (PK)
	if string(data[0:2]) != "PK" {
		return fmt.Errorf("invalid ePub/ZIP header")
	}

	return nil
}

// GetProgress returns the current progress of a generation job
func (o *BookOrchestrator) GetProgress(ctx context.Context, projectID uint) (*GenerationProgress, error) {
	// TODO: Implement progress tracking with Redis
	return &GenerationProgress{
		ProjectID:    projectID,
		Status:       "processing",
		CurrentStage: "rendering",
		Progress:     75,
	}, nil
}

// GenerationProgress tracks the current state of generation
type GenerationProgress struct {
	ProjectID    uint
	Status       string
	CurrentStage string
	Progress     int
	Message      string
}

// CancelGeneration cancels an in-progress generation
func (o *BookOrchestrator) CancelGeneration(ctx context.Context, projectID uint) error {
	// TODO: Implement cancellation logic
	return nil
}
