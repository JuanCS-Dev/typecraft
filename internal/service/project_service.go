package service

import (
	"fmt"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"github.com/JuanCS-Dev/typecraft/internal/repository"
	"github.com/google/uuid"
)

// ProjectService contém a lógica de negócio para projetos
type ProjectService struct {
	projectRepo *repository.ProjectRepository
	jobRepo     *repository.JobRepository
}

// NewProjectService cria uma nova instância do serviço
func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo: repository.NewProjectRepository(),
		jobRepo:     repository.NewJobRepository(),
	}
}

// CreateProjectRequest representa os dados para criar um projeto
type CreateProjectRequest struct {
	Title                string   `json:"title" binding:"required"`
	Author               string   `json:"author" binding:"required"`
	Genre                string   `json:"genre"`
	ISBN                 string   `json:"isbn"`
	Description          string   `json:"description"`
	PageFormat           string   `json:"page_format"`
	DistributionChannels []string `json:"distribution_channels"`
}

// CreateProject cria um novo projeto
func (s *ProjectService) CreateProject(userID string, req CreateProjectRequest) (*domain.Project, error) {
	// Validações
	if req.Title == "" {
		return nil, fmt.Errorf("título é obrigatório")
	}
	if req.Author == "" {
		return nil, fmt.Errorf("autor é obrigatório")
	}
	
	// Defaults
	if req.PageFormat == "" {
		req.PageFormat = "6x9"
	}
	if len(req.DistributionChannels) == 0 {
		req.DistributionChannels = []string{"kdp"}
	}
	
	// Criar projeto
	project := &domain.Project{
		UserID:               userID,
		Title:                req.Title,
		Author:               req.Author,
		Genre:                req.Genre,
		ISBN:                 req.ISBN,
		Description:          req.Description,
		PageFormat:           req.PageFormat,
		DistributionChannels: &req.DistributionChannels,
		Status:               domain.StatusCreated,
		Progress:             0,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	
	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("erro ao criar projeto: %w", err)
	}
	
	return project, nil
}

// GetProject busca um projeto por ID
func (s *ProjectService) GetProject(id string) (*domain.Project, error) {
	return s.projectRepo.GetByID(id)
}

// ListProjects lista projetos com paginação
func (s *ProjectService) ListProjects(userID string, page, pageSize int) ([]*domain.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	return s.projectRepo.GetAll(userID, pageSize, offset)
}

// UpdateProject atualiza metadados de um projeto
func (s *ProjectService) UpdateProject(id string, updates map[string]interface{}) (*domain.Project, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// Atualizar campos permitidos
	if title, ok := updates["title"].(string); ok && title != "" {
		project.Title = title
	}
	if author, ok := updates["author"].(string); ok && author != "" {
		project.Author = author
	}
	if genre, ok := updates["genre"].(string); ok {
		project.Genre = genre
	}
	if isbn, ok := updates["isbn"].(string); ok {
		project.ISBN = isbn
	}
	if description, ok := updates["description"].(string); ok {
		project.Description = description
	}
	
	project.UpdatedAt = time.Now()
	
	if err := s.projectRepo.Update(project); err != nil {
		return nil, fmt.Errorf("erro ao atualizar projeto: %w", err)
	}
	
	return project, nil
}

// DeleteProject deleta um projeto e seus jobs
func (s *ProjectService) DeleteProject(id string) error {
	// Verificar se existe
	if _, err := s.projectRepo.GetByID(id); err != nil {
		return err
	}
	
	// Deletar jobs associados
	if err := s.jobRepo.DeleteByProjectID(id); err != nil {
		return fmt.Errorf("erro ao deletar jobs do projeto: %w", err)
	}
	
	// Deletar projeto
	if err := s.projectRepo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar projeto: %w", err)
	}
	
	return nil
}

// SetManuscriptURL atualiza a URL do manuscrito
func (s *ProjectService) SetManuscriptURL(projectID, url string) error {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	
	project.ManuscriptURL = url
	project.Status = domain.StatusUploading
	project.UpdatedAt = time.Now()
	
	return s.projectRepo.Update(project)
}

// StartProcessing inicia o processamento de um projeto
func (s *ProjectService) StartProcessing(projectID string) error {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	
	// Validar se pode processar
	if !project.CanBeProcessed() {
		return fmt.Errorf("projeto não pode ser processado (status: %s, manuscript: %s)", 
			project.Status, project.ManuscriptURL)
	}
	
	// Criar jobs de processamento
	jobs := []domain.Job{
		{
			ID:          uuid.New().String(),
			ProjectID:   projectID,
			Type:        domain.JobTypeConvert,
			Status:      domain.JobStatusPending,
			Priority:    10,
			MaxAttempts: 3,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			ProjectID:   projectID,
			Type:        domain.JobTypeAnalyze,
			Status:      domain.JobStatusPending,
			Priority:    9,
			MaxAttempts: 3,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			ProjectID:   projectID,
			Type:        domain.JobTypeDesign,
			Status:      domain.JobStatusPending,
			Priority:    8,
			MaxAttempts: 3,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			ProjectID:   projectID,
			Type:        domain.JobTypeRender,
			Status:      domain.JobStatusPending,
			Priority:    7,
			MaxAttempts: 3,
			CreatedAt:   time.Now(),
		},
	}
	
	// Criar jobs no banco
	for i := range jobs {
		if err := s.jobRepo.Create(&jobs[i]); err != nil {
			return fmt.Errorf("erro ao criar job: %w", err)
		}
	}
	
	// Atualizar status do projeto
	project.Status = domain.StatusAnalyzing
	project.Progress = 10
	project.UpdatedAt = time.Now()
	
	return s.projectRepo.Update(project)
}

// GetProjectJobs retorna todos os jobs de um projeto
func (s *ProjectService) GetProjectJobs(projectID string) ([]*domain.Job, error) {
	return s.jobRepo.GetByProjectID(projectID)
}
