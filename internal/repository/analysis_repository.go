// Package repository provides data access layer for AI analyses
// Following Constituição Vértice v3.0 - Artigo I, Seção 2
package repository

import (
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/database"
	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"gorm.io/gorm"
)

// AnalysisRepository handles database operations for AI analyses
type AnalysisRepository struct {
	db *gorm.DB
}

// NewAnalysisRepository creates a new analysis repository
func NewAnalysisRepository() *AnalysisRepository {
	return &AnalysisRepository{
		db: database.DB,
	}
}

// Save saves an AI analysis to the database
func (r *AnalysisRepository) Save(analysis *domain.AIAnalysis) error {
	if err := analysis.BeforeCreate(); err != nil {
		return err
	}
	return r.db.Create(analysis).Error
}

// GetByProjectID retrieves the most recent analysis for a project
func (r *AnalysisRepository) GetByProjectID(projectID string) (*domain.AIAnalysis, error) {
	var analysis domain.AIAnalysis
	err := r.db.Where("project_id = ?", projectID).
		Order("analyzed_at DESC").
		First(&analysis).Error
	
	if err != nil {
		return nil, err
	}
	
	// Trigger AfterFind hook
	if err := analysis.AfterFind(); err != nil {
		return nil, err
	}
	
	return &analysis, nil
}

// GetByID retrieves an analysis by its ID
func (r *AnalysisRepository) GetByID(id string) (*domain.AIAnalysis, error) {
	var analysis domain.AIAnalysis
	err := r.db.First(&analysis, "id = ?", id).Error
	
	if err != nil {
		return nil, err
	}
	
	if err := analysis.AfterFind(); err != nil {
		return nil, err
	}
	
	return &analysis, nil
}

// GetCachedAnalysis retrieves a cached analysis if it exists and is still valid
// Returns nil if no valid cache exists
func (r *AnalysisRepository) GetCachedAnalysis(projectID string, maxAge time.Duration) (*domain.AIAnalysis, error) {
	var analysis domain.AIAnalysis
	cutoff := time.Now().Add(-maxAge)
	
	err := r.db.Where("project_id = ? AND analyzed_at > ?", projectID, cutoff).
		Order("analyzed_at DESC").
		First(&analysis).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, nil // No cache found, not an error
	}
	
	if err != nil {
		return nil, err
	}
	
	if err := analysis.AfterFind(); err != nil {
		return nil, err
	}
	
	return &analysis, nil
}

// Update updates an existing analysis
func (r *AnalysisRepository) Update(analysis *domain.AIAnalysis) error {
	if err := analysis.BeforeCreate(); err != nil {
		return err
	}
	return r.db.Save(analysis).Error
}

// Delete deletes an analysis by ID
func (r *AnalysisRepository) Delete(id string) error {
	return r.db.Delete(&domain.AIAnalysis{}, "id = ?", id).Error
}

// ListByProject lists all analyses for a project, ordered by most recent
func (r *AnalysisRepository) ListByProject(projectID string, limit int) ([]*domain.AIAnalysis, error) {
	var analyses []*domain.AIAnalysis
	query := r.db.Where("project_id = ?", projectID).
		Order("analyzed_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&analyses).Error
	if err != nil {
		return nil, err
	}
	
	// Trigger AfterFind for each
	for _, analysis := range analyses {
		if err := analysis.AfterFind(); err != nil {
			return nil, err
		}
	}
	
	return analyses, nil
}

// CountByProject counts total analyses for a project
func (r *AnalysisRepository) CountByProject(projectID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.AIAnalysis{}).
		Where("project_id = ?", projectID).
		Count(&count).Error
	return count, err
}

// GetTotalTokensUsed calculates total tokens used for a project
func (r *AnalysisRepository) GetTotalTokensUsed(projectID string) (int, error) {
	var total int
	err := r.db.Model(&domain.AIAnalysis{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(SUM(tokens_used), 0)").
		Scan(&total).Error
	return total, err
}

// DeleteOldAnalyses deletes analyses older than the specified duration
func (r *AnalysisRepository) DeleteOldAnalyses(maxAge time.Duration) (int64, error) {
	cutoff := time.Now().Add(-maxAge)
	result := r.db.Where("analyzed_at < ?", cutoff).
		Delete(&domain.AIAnalysis{})
	return result.RowsAffected, result.Error
}
