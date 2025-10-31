package repository

import (
	"fmt"

	"github.com/JuanCS-Dev/typecraft/internal/database"
	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"gorm.io/gorm"
)

// ProjectRepository lida com operações de banco de dados para Projects
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository cria uma nova instância do repositório
func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		db: database.DB,
	}
}

// Create cria um novo projeto
func (r *ProjectRepository) Create(project *domain.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return fmt.Errorf("erro ao criar projeto: %w", err)
	}
	return nil
}

// GetByID busca um projeto por ID
func (r *ProjectRepository) GetByID(id string) (*domain.Project, error) {
	var project domain.Project
	if err := r.db.First(&project, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("projeto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar projeto: %w", err)
	}
	return &project, nil
}

// GetAll lista todos os projetos com paginação
func (r *ProjectRepository) GetAll(userID string, limit, offset int) ([]*domain.Project, int64, error) {
	var projects []*domain.Project
	var total int64
	
	query := r.db.Model(&domain.Project{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	
	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar projetos: %w", err)
	}
	
	// Buscar com paginação
	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar projetos: %w", err)
	}
	
	return projects, total, nil
}

// Update atualiza um projeto
func (r *ProjectRepository) Update(project *domain.Project) error {
	if err := r.db.Save(project).Error; err != nil {
		return fmt.Errorf("erro ao atualizar projeto: %w", err)
	}
	return nil
}

// Delete deleta um projeto
func (r *ProjectRepository) Delete(id string) error {
	if err := r.db.Delete(&domain.Project{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("erro ao deletar projeto: %w", err)
	}
	return nil
}

// GetByStatus busca projetos por status
func (r *ProjectRepository) GetByStatus(status domain.ProjectStatus, limit, offset int) ([]*domain.Project, error) {
	var projects []*domain.Project
	
	if err := r.db.Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar projetos por status: %w", err)
	}
	
	return projects, nil
}

// GetProcessable busca projetos prontos para processar
func (r *ProjectRepository) GetProcessable(limit int) ([]*domain.Project, error) {
	var projects []*domain.Project
	
	if err := r.db.Where("status = ? AND manuscript_url != ''", domain.StatusCreated).
		Limit(limit).
		Order("created_at ASC").
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar projetos processáveis: %w", err)
	}
	
	return projects, nil
}

// UpdateStatus atualiza apenas o status de um projeto
func (r *ProjectRepository) UpdateStatus(id string, status domain.ProjectStatus, progress int) error {
	if err := r.db.Model(&domain.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":   status,
			"progress": progress,
		}).Error; err != nil {
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}
	return nil
}
