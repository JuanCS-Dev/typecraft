package domain

import (
	"time"
)

// ProjectStatus representa o estado atual do processamento
type ProjectStatus string

const (
	StatusCreated    ProjectStatus = "created"
	StatusUploading  ProjectStatus = "uploading"
	StatusAnalyzing  ProjectStatus = "analyzing"
	StatusDesigning  ProjectStatus = "designing"
	StatusRendering  ProjectStatus = "rendering"
	StatusRefining   ProjectStatus = "refining"
	StatusCompleted  ProjectStatus = "completed"
	StatusFailed     ProjectStatus = "failed"
)

// Project representa um projeto de livro
type Project struct {
	ID          string        `json:"id" gorm:"primaryKey"`
	UserID      string        `json:"user_id" gorm:"index;not null"`
	
	// Metadados do livro
	Title       string        `json:"title" gorm:"not null"`
	Author      string        `json:"author" gorm:"not null"`
	Genre       string        `json:"genre"`
	ISBN        string        `json:"isbn"`
	Description string        `json:"description"`
	
	// Configurações
	PageFormat           string    `json:"page_format" gorm:"default:'6x9'"`
	DistributionChannels *[]string `json:"distribution_channels,omitempty" gorm:"type:jsonb"`
	
	// Status do processamento
	Status   ProjectStatus `json:"status" gorm:"default:'created';index"`
	Progress int           `json:"progress" gorm:"default:0"` // 0-100
	
	// Análise de conteúdo (JSON gerado pela IA)
	Analysis *map[string]interface{} `json:"analysis,omitempty" gorm:"type:jsonb"`
	
	// Design gerado (JSON com fontes, cores, layout)
	DesignConfig *map[string]interface{} `json:"design_config,omitempty" gorm:"type:jsonb"`
	
	// URLs dos arquivos (MinIO/S3)
	ManuscriptURL      string `json:"manuscript_url,omitempty"`
	PDFKdpURL          string `json:"pdf_kdp_url,omitempty"`
	PDFIngramSparkURL  string `json:"pdf_ingramspark_url,omitempty"`
	EpubURL            string `json:"epub_url,omitempty"`
	CoverURL           string `json:"cover_url,omitempty"`
	
	// Timestamps
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	// Logs de erro (JSON array)
	ErrorLog *[]string `json:"error_log,omitempty" gorm:"type:jsonb"`
}

// TableName especifica o nome da tabela no banco
func (Project) TableName() string {
	return "projects"
}

// IsCompleted verifica se o projeto foi concluído
func (p *Project) IsCompleted() bool {
	return p.Status == StatusCompleted
}

// IsFailed verifica se o projeto falhou
func (p *Project) IsFailed() bool {
	return p.Status == StatusFailed
}

// CanBeProcessed verifica se o projeto pode ser processado
func (p *Project) CanBeProcessed() bool {
	return p.Status == StatusCreated && p.ManuscriptURL != ""
}

// AddError adiciona um erro ao log
func (p *Project) AddError(err string) {
	if p.ErrorLog == nil {
		p.ErrorLog = &[]string{}
	}
	*p.ErrorLog = append(*p.ErrorLog, err)
}

// SetStatus atualiza o status e o timestamp apropriado
func (p *Project) SetStatus(status ProjectStatus) {
	p.Status = status
	if status == StatusCompleted {
		now := time.Now()
		p.CompletedAt = &now
		p.Progress = 100
	} else if status == StatusFailed {
		now := time.Now()
		p.CompletedAt = &now
	}
}
