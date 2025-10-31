package repository

import (
	"fmt"

	"github.com/JuanCS-Dev/typecraft/internal/database"
	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"gorm.io/gorm"
)

// JobRepository lida com operações de banco de dados para Jobs
type JobRepository struct {
	db *gorm.DB
}

// NewJobRepository cria uma nova instância do repositório
func NewJobRepository() *JobRepository {
	return &JobRepository{
		db: database.DB,
	}
}

// Create cria um novo job
func (r *JobRepository) Create(job *domain.Job) error {
	if err := r.db.Create(job).Error; err != nil {
		return fmt.Errorf("erro ao criar job: %w", err)
	}
	return nil
}

// GetByID busca um job por ID
func (r *JobRepository) GetByID(id string) (*domain.Job, error) {
	var job domain.Job
	if err := r.db.First(&job, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar job: %w", err)
	}
	return &job, nil
}

// GetByProjectID busca todos os jobs de um projeto
func (r *JobRepository) GetByProjectID(projectID string) ([]*domain.Job, error) {
	var jobs []*domain.Job
	
	if err := r.db.Where("project_id = ?", projectID).
		Order("created_at ASC").
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar jobs do projeto: %w", err)
	}
	
	return jobs, nil
}

// Update atualiza um job
func (r *JobRepository) Update(job *domain.Job) error {
	if err := r.db.Save(job).Error; err != nil {
		return fmt.Errorf("erro ao atualizar job: %w", err)
	}
	return nil
}

// GetPending busca jobs pendentes para processar
func (r *JobRepository) GetPending(limit int) ([]*domain.Job, error) {
	var jobs []*domain.Job
	
	if err := r.db.Where("status = ?", domain.JobStatusPending).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar jobs pendentes: %w", err)
	}
	
	return jobs, nil
}

// GetByStatus busca jobs por status
func (r *JobRepository) GetByStatus(status domain.JobStatus, limit, offset int) ([]*domain.Job, error) {
	var jobs []*domain.Job
	
	if err := r.db.Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar jobs por status: %w", err)
	}
	
	return jobs, nil
}

// GetFailed busca jobs que falharam e podem ser retentados
func (r *JobRepository) GetFailed(limit int) ([]*domain.Job, error) {
	var jobs []*domain.Job
	
	if err := r.db.Where("status = ? AND attempts < max_attempts", domain.JobStatusFailed).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar jobs falhados: %w", err)
	}
	
	return jobs, nil
}

// Delete deleta um job
func (r *JobRepository) Delete(id string) error {
	if err := r.db.Delete(&domain.Job{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("erro ao deletar job: %w", err)
	}
	return nil
}

// DeleteByProjectID deleta todos os jobs de um projeto
func (r *JobRepository) DeleteByProjectID(projectID string) error {
	if err := r.db.Where("project_id = ?", projectID).Delete(&domain.Job{}).Error; err != nil {
		return fmt.Errorf("erro ao deletar jobs do projeto: %w", err)
	}
	return nil
}

// CountByProject conta quantos jobs um projeto tem
func (r *JobRepository) CountByProject(projectID string) (int64, error) {
	var count int64
	
	if err := r.db.Model(&domain.Job{}).
		Where("project_id = ?", projectID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("erro ao contar jobs: %w", err)
	}
	
	return count, nil
}

// CountByStatus conta jobs por status
func (r *JobRepository) CountByStatus(status domain.JobStatus) (int64, error) {
	var count int64
	
	if err := r.db.Model(&domain.Job{}).
		Where("status = ?", status).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("erro ao contar jobs por status: %w", err)
	}
	
	return count, nil
}
