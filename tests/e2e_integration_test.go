// Package tests contains End-to-End integration tests for the complete system.
// This validates the COMPLETE INTEGRATION as per Dia 7, Sprint 7-8.
package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"github.com/JuanCS-Dev/typecraft/internal/service"
	"github.com/JuanCS-Dev/typecraft/pkg/design"
)

// TestE2E_CompleteBookGeneration tests the complete workflow
func TestE2E_CompleteBookGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	
	// Setup test environment
	tmpDir := t.TempDir()
	
	// Create test content
	contentPath := filepath.Join(tmpDir, "manuscript.md")
	content := generateTestManuscript()
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test content: %v", err)
	}

	// Setup test project
	projectRepo := setupTestProjectRepository(t)
	project := &domain.Project{
		ID:         1,
		Title:      "E2E Test Book",
		Author:     "Test Author",
		Genre:      "Fiction",
		Language:   "pt",
		PageFormat: "6x9",
	}
	if err := projectRepo.Create(ctx, project); err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Setup analysis client
	analysisClient := setupTestAnalysisClient(t)

	// Create orchestrator
	orchestrator := service.NewBookOrchestrator(
		projectRepo,
		analysisClient,
		tmpDir,
	)

	// Execute generation
	req := &service.GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf", "epub"},
	}

	t.Logf("Starting book generation...")
	start := time.Now()

	result, err := orchestrator.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	duration := time.Since(start)
	t.Logf("Generation completed in %v", duration)

	// Validate results
	validateGenerationResult(t, result, tmpDir)

	// Performance assertions
	if duration > 5*time.Minute {
		t.Errorf("Generation took too long: %v (expected < 5 minutes)", duration)
	}

	t.Logf("✓ E2E test PASSED")
	t.Logf("  Pipeline: %s", result.Pipeline)
	t.Logf("  Outputs: %v", result.OutputFiles)
	t.Logf("  Duration: %v", duration)
	t.Logf("  File Size: %d bytes", result.Metrics.FileSize)
}

// TestE2E_LaTeXPipeline tests the LaTeX pipeline specifically
func TestE2E_LaTeXPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	tmpDir := t.TempDir()

	// Create academic content (should trigger LaTeX)
	contentPath := filepath.Join(tmpDir, "academic.md")
	content := generateAcademicManuscript()
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create content: %v", err)
	}

	projectRepo := setupTestProjectRepository(t)
	project := &domain.Project{
		ID:         2,
		Title:      "Academic Book",
		Author:     "Prof. Test",
		Genre:      "Academic",
		PageFormat: "6x9",
	}
	projectRepo.Create(ctx, project)

	analysisClient := setupTestAnalysisClient(t)
	orchestrator := service.NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &service.GenerationRequest{
		ProjectID:     2,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	result, err := orchestrator.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify LaTeX pipeline was selected
	if result.Pipeline != "latex" {
		t.Errorf("Expected LaTeX pipeline for academic content, got %s", result.Pipeline)
	}

	validateGenerationResult(t, result, tmpDir)
}

// TestE2E_HTMLPipeline tests the HTML/CSS pipeline specifically
func TestE2E_HTMLPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	tmpDir := t.TempDir()

	// Create visual content (should trigger HTML)
	contentPath := filepath.Join(tmpDir, "visual.md")
	content := generateVisualManuscript()
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create content: %v", err)
	}

	projectRepo := setupTestProjectRepository(t)
	project := &domain.Project{
		ID:         3,
		Title:      "Visual Book",
		Author:     "Artist Test",
		Genre:      "Art",
		PageFormat: "8.5x11",
	}
	projectRepo.Create(ctx, project)

	analysisClient := setupTestAnalysisClient(t)
	orchestrator := service.NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &service.GenerationRequest{
		ProjectID:     3,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	result, err := orchestrator.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify HTML pipeline was selected
	if result.Pipeline != "html" {
		t.Errorf("Expected HTML pipeline for visual content, got %s", result.Pipeline)
	}

	validateGenerationResult(t, result, tmpDir)
}

// TestE2E_CustomDesign tests custom design options
func TestE2E_CustomDesign(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	tmpDir := t.TempDir()

	contentPath := filepath.Join(tmpDir, "custom.md")
	content := generateTestManuscript()
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create content: %v", err)
	}

	projectRepo := setupTestProjectRepository(t)
	project := &domain.Project{
		ID:         4,
		Title:      "Custom Design Book",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(ctx, project)

	analysisClient := setupTestAnalysisClient(t)
	orchestrator := service.NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	customDesign := &service.DesignOptions{
		BodyFont:    "Garamond",
		HeadingFont: "Futura Bold",
		ColorScheme: []string{"#2C3E50", "#E74C3C"},
		CustomMargins: &design.Margins{
			Top:    25.0,
			Bottom: 25.0,
			Left:   20.0,
			Right:  15.0,
		},
	}

	req := &service.GenerationRequest{
		ProjectID:     4,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
		CustomDesign:  customDesign,
	}

	result, err := orchestrator.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify custom design was applied
	if result.DesignMetadata.Fonts.Body != "Garamond" {
		t.Errorf("Expected custom body font, got %s", result.DesignMetadata.Fonts.Body)
	}

	if result.DesignMetadata.Fonts.Heading != "Futura Bold" {
		t.Errorf("Expected custom heading font, got %s", result.DesignMetadata.Fonts.Heading)
	}

	validateGenerationResult(t, result, tmpDir)
}

// TestE2E_MultipleFormats tests simultaneous PDF and ePub generation
func TestE2E_MultipleFormats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	tmpDir := t.TempDir()

	contentPath := filepath.Join(tmpDir, "multi.md")
	content := generateTestManuscript()
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create content: %v", err)
	}

	projectRepo := setupTestProjectRepository(t)
	project := &domain.Project{
		ID:         5,
		Title:      "Multi-Format Book",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(ctx, project)

	analysisClient := setupTestAnalysisClient(t)
	orchestrator := service.NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &service.GenerationRequest{
		ProjectID:     5,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf", "epub"},
	}

	result, err := orchestrator.Generate(ctx, req)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Verify both formats were generated
	if len(result.OutputFiles) != 2 {
		t.Errorf("Expected 2 output files, got %d", len(result.OutputFiles))
	}

	pdfPath, hasPDF := result.OutputFiles["pdf"]
	if !hasPDF {
		t.Error("Expected PDF output")
	}

	epubPath, hasEPUB := result.OutputFiles["epub"]
	if !hasEPUB {
		t.Error("Expected ePub output")
	}

	// Verify both files exist and are valid
	validatePDFFile(t, pdfPath)
	validateEPUBFile(t, epubPath)
}

// Helper functions

func validateGenerationResult(t *testing.T, result *service.GenerationResult, tmpDir string) {
	t.Helper()

	if !result.Success {
		t.Fatalf("Generation was not successful: %v", result.Error)
	}

	if result.Pipeline == "" {
		t.Error("Pipeline not set")
	}

	if len(result.OutputFiles) == 0 {
		t.Error("No output files generated")
	}

	if result.DesignMetadata == nil {
		t.Error("Design metadata not set")
	}

	if result.Metrics == nil {
		t.Error("Metrics not set")
	}

	// Validate each output file
	for format, path := range result.OutputFiles {
		if !filepath.IsAbs(path) {
			t.Errorf("Output path is not absolute: %s", path)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Output file %s does not exist: %v", format, err)
			continue
		}

		if info.Size() == 0 {
			t.Errorf("Output file %s is empty", format)
		}

		t.Logf("✓ %s file: %s (%d bytes)", format, path, info.Size())
	}

	// Validate metrics
	if result.Metrics.Duration == 0 {
		t.Error("Duration not recorded")
	}

	t.Logf("✓ Metrics: Analysis=%dms, Design=%dms, Rendering=%dms",
		result.Metrics.ContentAnalysisMs,
		result.Metrics.DesignGenerationMs,
		result.Metrics.RenderingMs)
}

func validatePDFFile(t *testing.T, path string) {
	t.Helper()

	data := make([]byte, 4)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open PDF: %v", err)
	}
	defer f.Close()

	if _, err := f.Read(data); err != nil {
		t.Fatalf("Failed to read PDF header: %v", err)
	}

	if string(data) != "%PDF" {
		t.Error("Invalid PDF header")
	}
}

func validateEPUBFile(t *testing.T, path string) {
	t.Helper()

	// Basic validation: file exists and has .epub extension
	if filepath.Ext(path) != ".epub" {
		t.Errorf("Invalid ePub file extension: %s", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("ePub file not found: %v", err)
	}

	if info.Size() == 0 {
		t.Error("ePub file is empty")
	}
}

// Test content generators

func generateTestManuscript() string {
	return `# The Complete Test Book

## Chapter 1: Introduction

This is a comprehensive test manuscript designed to exercise all features of the book generation system.

### Section 1.1: Typography

The quick brown fox jumps over the lazy dog. This sentence contains every letter of the alphabet and is commonly used for testing fonts.

**Bold text** and *italic text* should be properly rendered.

### Section 1.2: Lists

Unordered list:
- Item 1
- Item 2
- Item 3

Ordered list:
1. First item
2. Second item
3. Third item

## Chapter 2: Advanced Features

### Mathematical Equations

The famous equation: $E = mc^2$

### Tables

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Data 1   | Data 2   | Data 3   |
| Data 4   | Data 5   | Data 6   |

### Code Blocks

` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

## Chapter 3: Conclusion

This concludes our test manuscript. The system should handle all these elements correctly.
`
}

func generateAcademicManuscript() string {
	return `# Research Paper: Advanced Topics

## Abstract

This academic paper contains mathematical formulas and technical content.

## Section 1: Mathematical Analysis

Consider the integral:

$$\int_{0}^{\infty} e^{-x^2} dx = \frac{\sqrt{\pi}}{2}$$

## Section 2: Algorithm

The algorithm complexity is $O(n \log n)$.

` + "```python" + `
def quicksort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    return quicksort(left) + middle + quicksort(right)
` + "```" + `

## Conclusion

Mathematical content requires LaTeX for proper rendering.
`
}

func generateVisualManuscript() string {
	return `# Visual Design Book

## Chapter 1: Layout Principles

This book contains many visual elements and should use the HTML pipeline.

![Placeholder Image 1](image1.jpg)

### Design Patterns

![Placeholder Image 2](image2.jpg)

Text mixed with images requires flexible layout.

![Placeholder Image 3](image3.jpg)

## Chapter 2: Color Theory

More visual content with multiple images per page.

![Placeholder Image 4](image4.jpg)
![Placeholder Image 5](image5.jpg)
`
}

// Mock implementations

func setupTestProjectRepository(t *testing.T) domain.ProjectRepository {
	return &mockProjectRepo{
		projects: make(map[uint]*domain.Project),
	}
}

func setupTestAnalysisClient(t *testing.T) service.AnalysisClient {
	return &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Fiction",
			Tone:       "Neutral",
			Complexity: 0.5,
			HasMath:    false,
			ImageCount: 0,
			TableCount: 1,
		},
	}
}

type mockProjectRepo struct {
	projects map[uint]*domain.Project
}

func (m *mockProjectRepo) GetByID(ctx context.Context, id uint) (*domain.Project, error) {
	if p, ok := m.projects[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("project not found")
}

func (m *mockProjectRepo) Create(ctx context.Context, project *domain.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepo) Update(ctx context.Context, project *domain.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepo) Delete(ctx context.Context, id uint) error {
	delete(m.projects, id)
	return nil
}

func (m *mockProjectRepo) List(ctx context.Context) ([]*domain.Project, error) {
	var projects []*domain.Project
	for _, p := range m.projects {
		projects = append(projects, p)
	}
	return projects, nil
}

type mockAnalysisClient struct {
	analysis *domain.Analysis
	err      error
}

func (m *mockAnalysisClient) AnalyzeContent(ctx context.Context, content string) (*domain.Analysis, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.analysis, nil
}
