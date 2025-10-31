package repository

import (
	"testing"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnalysisRepository_Cache(t *testing.T) {
	// Skip if DB not available
	if testing.Short() {
		t.Skip("Skipping repository test in short mode")
	}
	
	repo := NewAnalysisRepository()
	if repo.db == nil {
		t.Skip("Database not available")
	}
	
	// Create test analysis
	projectID := uuid.New().String()
	analysis := &domain.AIAnalysis{
		ID:                  uuid.New().String(),
		ProjectID:           projectID,
		Genre:               "fiction",
		GenreConfidence:     0.95,
		SubGenres:           []string{"literary", "contemporary"},
		RecommendedPipeline: "html",
		WordCount:           10000,
		EstimatedPages:      50,
		AnalyzedAt:          time.Now(),
	}
	
	analysis.Tone = domain.ToneAnalysis{
		Primary:    "formal",
		Formality:  0.8,
		Emotion:    "serene",
		Confidence: 0.9,
	}
	
	analysis.Complexity = domain.ComplexityMetrics{
		AvgSentenceLength:  15.5,
		VocabularyRichness: 0.7,
		SyntaxComplexity:   0.6,
		TechnicalDensity:   0.2,
		ReadingLevel:       "college",
	}
	
	// Test Save
	err := repo.Save(analysis)
	assert.NoError(t, err)
	assert.NotEmpty(t, analysis.ID)
	
	// Test GetByProjectID
	retrieved, err := repo.GetByProjectID(projectID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, analysis.Genre, retrieved.Genre)
	assert.Equal(t, analysis.GenreConfidence, retrieved.GenreConfidence)
	assert.Equal(t, len(analysis.SubGenres), len(retrieved.SubGenres))
	
	// Test GetCachedAnalysis - should find it (within 24h)
	cached, err := repo.GetCachedAnalysis(projectID, 24*time.Hour)
	assert.NoError(t, err)
	assert.NotNil(t, cached)
	assert.Equal(t, analysis.ID, cached.ID)
	
	// Test GetCachedAnalysis - should not find if expired
	cached, err = repo.GetCachedAnalysis(projectID, 1*time.Nanosecond)
	assert.NoError(t, err)
	assert.Nil(t, cached)
	
	// Test CountByProject
	count, err := repo.CountByProject(projectID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	
	// Cleanup
	err = repo.Delete(analysis.ID)
	assert.NoError(t, err)
}

func TestAnalysisRepository_ListByProject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping repository test in short mode")
	}
	
	repo := NewAnalysisRepository()
	if repo.db == nil {
		t.Skip("Database not available")
	}
	
	projectID := uuid.New().String()
	
	// Create multiple analyses
	for i := 0; i < 3; i++ {
		analysis := &domain.AIAnalysis{
			ID:                  uuid.New().String(),
			ProjectID:           projectID,
			Genre:               "technical",
			GenreConfidence:     0.9,
			RecommendedPipeline: "latex",
			WordCount:           5000,
			EstimatedPages:      25,
			AnalyzedAt:          time.Now().Add(-time.Duration(i) * time.Hour),
		}
		
		err := repo.Save(analysis)
		assert.NoError(t, err)
	}
	
	// List all
	analyses, err := repo.ListByProject(projectID, 0)
	assert.NoError(t, err)
	assert.Len(t, analyses, 3)
	
	// Should be ordered by most recent first
	assert.True(t, analyses[0].AnalyzedAt.After(analyses[1].AnalyzedAt))
	assert.True(t, analyses[1].AnalyzedAt.After(analyses[2].AnalyzedAt))
	
	// List with limit
	analyses, err = repo.ListByProject(projectID, 2)
	assert.NoError(t, err)
	assert.Len(t, analyses, 2)
	
	// Cleanup
	for _, a := range analyses {
		repo.Delete(a.ID)
	}
}
