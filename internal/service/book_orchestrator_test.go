package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
)

// mockAnalysisClient implements AnalysisClient for testing
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

// mockProjectRepository implements ProjectRepository for testing
type mockProjectRepository struct {
	projects map[uint]*domain.Project
}

func newMockProjectRepository() *mockProjectRepository {
	return &mockProjectRepository{
		projects: make(map[uint]*domain.Project),
	}
}

func (m *mockProjectRepository) GetByID(ctx context.Context, id uint) (*domain.Project, error) {
	if project, ok := m.projects[id]; ok {
		return project, nil
	}
	return nil, domain.ErrProjectNotFound
}

func (m *mockProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) Delete(ctx context.Context, id uint) error {
	delete(m.projects, id)
	return nil
}

func (m *mockProjectRepository) List(ctx context.Context) ([]*domain.Project, error) {
	var projects []*domain.Project
	for _, p := range m.projects {
		projects = append(projects, p)
	}
	return projects, nil
}

// Helper functions for tests

func setupTestEnvironment(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "typecraft-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func createTestContent(t *testing.T, dir string) string {
	contentPath := filepath.Join(dir, "test-content.md")
	content := `# Test Book

## Chapter 1: Introduction

This is a test book with multiple chapters and some complex content.

### Section 1.1

Lorem ipsum dolor sit amet, consectetur adipiscing elit.

### Section 1.2

More content with **bold** and *italic* text.

## Chapter 2: Advanced Topics

This chapter covers more advanced topics.

$$
E = mc^2
$$

Some mathematical equations and formulas.

| Name | Age | City |
|------|-----|------|
| John | 30  | NYC  |
| Jane | 25  | LA   |

A table with data.
`

	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test content: %v", err)
	}

	return contentPath
}

// Test Cases

func TestBookOrchestrator_Generate_BasicFlow(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	contentPath := createTestContent(t, tmpDir)

	// Setup mocks
	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Test Book",
		Author:     "Test Author",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(context.Background(), project)

	analysisClient := &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Fiction",
			Tone:       "Neutral",
			Complexity: 0.5,
			HasMath:    true,
			ImageCount: 0,
			TableCount: 1,
		},
	}

	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	// Test generation
	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	result, err := orchestrator.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success=true, got false")
	}

	if result.Pipeline == "" {
		t.Errorf("Expected pipeline to be selected")
	}

	if len(result.OutputFiles) == 0 {
		t.Errorf("Expected output files to be generated")
	}

	// Verify metrics
	if result.Metrics.Duration == 0 {
		t.Errorf("Expected duration to be recorded")
	}
}

func TestBookOrchestrator_PipelineSelection_LaTeX(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	contentPath := createTestContent(t, tmpDir)

	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Math Book",
		Genre:      "Academic",
		PageFormat: "6x9",
	}
	projectRepo.Create(context.Background(), project)

	// Analysis indicating LaTeX is better (high math content)
	analysisClient := &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Academic",
			Tone:       "Formal",
			Complexity: 0.8,
			HasMath:    true,
			ImageCount: 0,
			TableCount: 5,
			CodeBlocks: 10,
		},
	}

	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	result, err := orchestrator.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.Pipeline != "latex" {
		t.Errorf("Expected LaTeX pipeline for math-heavy content, got %s", result.Pipeline)
	}
}

func TestBookOrchestrator_PipelineSelection_HTML(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	contentPath := createTestContent(t, tmpDir)

	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Visual Book",
		Genre:      "Art",
		PageFormat: "8.5x11",
	}
	projectRepo.Create(context.Background(), project)

	// Analysis indicating HTML is better (high image content)
	analysisClient := &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Art",
			Tone:       "Creative",
			Complexity: 0.3,
			HasMath:    false,
			ImageCount: 25,
			TableCount: 0,
		},
	}

	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	result, err := orchestrator.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.Pipeline != "html" {
		t.Errorf("Expected HTML pipeline for image-heavy content, got %s", result.Pipeline)
	}
}

func TestBookOrchestrator_MultiFormat(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	contentPath := createTestContent(t, tmpDir)

	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Multi-Format Book",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(context.Background(), project)

	analysisClient := &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Fiction",
			Tone:       "Neutral",
			Complexity: 0.5,
		},
	}

	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf", "epub"},
	}

	result, err := orchestrator.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(result.OutputFiles) != 2 {
		t.Errorf("Expected 2 output files, got %d", len(result.OutputFiles))
	}

	if _, hasPDF := result.OutputFiles["pdf"]; !hasPDF {
		t.Errorf("Expected PDF output")
	}

	if _, hasEPUB := result.OutputFiles["epub"]; !hasEPUB {
		t.Errorf("Expected ePub output")
	}
}

func TestBookOrchestrator_CustomDesign(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	contentPath := createTestContent(t, tmpDir)

	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Custom Design Book",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(context.Background(), project)

	analysisClient := &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Fiction",
			Tone:       "Neutral",
			Complexity: 0.5,
		},
	}

	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	customDesign := &DesignOptions{
		BodyFont:    "Garamond",
		HeadingFont: "Futura",
		ColorScheme: []string{"#FF5733", "#C70039"},
	}

	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
		CustomDesign:  customDesign,
	}

	result, err := orchestrator.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.DesignMetadata.Fonts.Body != "Garamond" {
		t.Errorf("Expected custom body font Garamond, got %s", result.DesignMetadata.Fonts.Body)
	}

	if result.DesignMetadata.Fonts.Heading != "Futura" {
		t.Errorf("Expected custom heading font Futura, got %s", result.DesignMetadata.Fonts.Heading)
	}
}

func TestBookOrchestrator_InvalidProject(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	contentPath := createTestContent(t, tmpDir)

	projectRepo := newMockProjectRepository()
	analysisClient := &mockAnalysisClient{}
	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &GenerationRequest{
		ProjectID:     999, // Non-existent project
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	result, err := orchestrator.Generate(context.Background(), req)
	if err == nil {
		t.Error("Expected error for invalid project")
	}

	if result.Success {
		t.Error("Expected success=false for invalid project")
	}
}

func TestBookOrchestrator_EmptyContent(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create empty content file
	contentPath := filepath.Join(tmpDir, "empty.md")
	if err := os.WriteFile(contentPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty content: %v", err)
	}

	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Empty Book",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(context.Background(), project)

	analysisClient := &mockAnalysisClient{}
	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	_, err := orchestrator.Generate(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty content")
	}
}

// Benchmark tests

func BenchmarkOrchestrator_Generate(b *testing.B) {
	tmpDir, cleanup := setupTestEnvironment(&testing.T{})
	defer cleanup()

	contentPath := createTestContent(&testing.T{}, tmpDir)

	projectRepo := newMockProjectRepository()
	project := &domain.Project{
		ID:         1,
		Title:      "Benchmark Book",
		Genre:      "Fiction",
		PageFormat: "6x9",
	}
	projectRepo.Create(context.Background(), project)

	analysisClient := &mockAnalysisClient{
		analysis: &domain.Analysis{
			Genre:      "Fiction",
			Complexity: 0.5,
		},
	}

	orchestrator := NewBookOrchestrator(projectRepo, analysisClient, tmpDir)

	req := &GenerationRequest{
		ProjectID:     1,
		ContentPath:   contentPath,
		OutputFormats: []string{"pdf"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orchestrator.Generate(context.Background(), req)
	}
}
