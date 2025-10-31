package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JuanCS-Dev/typecraft/internal/service"
	"github.com/gin-gonic/gin"
)

// ProjectHandler lida com requisições HTTP para projetos
type ProjectHandler struct {
	service *service.ProjectService
}

// NewProjectHandler cria uma nova instância do handler
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		service: service.NewProjectService(),
	}
}

// CreateProject godoc
// @Summary Criar novo projeto
// @Tags projects
// @Accept json
// @Produce json
// @Param request body service.CreateProjectRequest true "Dados do projeto"
// @Success 201 {object} domain.Project
// @Router /api/v1/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req service.CreateProjectRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// UserID padrão até implementação de autenticação (Sprint 3-4)
	userID := "default_user"
	
	project, err := h.service.CreateProject(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, project)
}

// GetProject godoc
// @Summary Buscar projeto por ID
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Project
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	
	project, err := h.service.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, project)
}

// ListProjects godoc
// @Summary Listar projetos
// @Tags projects
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Itens por página" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	// UserID padrão até implementação de autenticação (Sprint 3-4)
	userID := "default_user"
	
	projects, total, err := h.service.ListProjects(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":      projects,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateProject godoc
// @Summary Atualizar projeto
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body map[string]interface{} true "Campos para atualizar"
// @Success 200 {object} domain.Project
// @Router /api/v1/projects/{id} [patch]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	project, err := h.service.UpdateProject(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, project)
}

// DeleteProject godoc
// @Summary Deletar projeto
// @Tags projects
// @Param id path string true "Project ID"
// @Success 204
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.service.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}

// UploadManuscript godoc
// @Summary Upload do manuscrito
// @Tags projects
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Project ID"
// @Param file formData file true "Arquivo do manuscrito"
// @Success 200 {object} domain.Project
// @Router /api/v1/projects/{id}/upload [post]
func (h *ProjectHandler) UploadManuscript(c *gin.Context) {
	id := c.Param("id")
	
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo é obrigatório"})
		return
	}
	
	// Upload para MinIO/S3 (Sprint 3-4)
	// Por enquanto, apenas simular
	manuscriptURL := fmt.Sprintf("s3://typecraft-files/manuscripts/%s/%s", id, file.Filename)
	
	if err := h.service.SetManuscriptURL(id, manuscriptURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	project, _ := h.service.GetProject(id)
	c.JSON(http.StatusOK, project)
}

// ProcessProject godoc
// @Summary Iniciar processamento do projeto
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/process [post]
func (h *ProjectHandler) ProcessProject(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.service.StartProcessing(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Processamento iniciado",
		"project_id": id,
	})
}

// GetProjectJobs godoc
// @Summary Listar jobs do projeto
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} domain.Job
// @Router /api/v1/projects/{id}/jobs [get]
func (h *ProjectHandler) GetProjectJobs(c *gin.Context) {
	id := c.Param("id")
	
	jobs, err := h.service.GetProjectJobs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, jobs)
}
