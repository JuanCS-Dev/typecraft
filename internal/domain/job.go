package domain

import "time"

// Job representa uma tarefa assíncrona de processamento
type Job struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ProjectID   string    `json:"project_id" gorm:"index;not null"`
	Type        JobType   `json:"type" gorm:"not null"`
	Status      JobStatus `json:"status" gorm:"default:'pending';index"`
	Priority    int       `json:"priority" gorm:"default:5"` // 1-10, maior = mais prioritário
	Payload     *map[string]interface{} `json:"payload,omitempty" gorm:"type:jsonb"`
	Result      *map[string]interface{} `json:"result,omitempty" gorm:"type:jsonb"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
	Attempts    int       `json:"attempts" gorm:"default:0"`
	MaxAttempts int       `json:"max_attempts" gorm:"default:3"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// JobType define os tipos de jobs disponíveis
type JobType string

const (
	JobTypeConvert  JobType = "convert"  // Conversão de formato (docx → markdown)
	JobTypeAnalyze  JobType = "analyze"  // Análise de conteúdo (IA/PLN)
	JobTypeDesign   JobType = "design"   // Geração de design (IA)
	JobTypeRender   JobType = "render"   // Renderização (LaTeX/HTML → PDF)
	JobTypeRefine   JobType = "refine"   // Refinamento tipográfico (IA)
	JobTypeExport   JobType = "export"   // Exportação multi-formato
)

// JobStatus define os estados possíveis de um job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// TableName especifica o nome da tabela no banco
func (Job) TableName() string {
	return "jobs"
}

// IsTerminal verifica se o job está em estado final
func (j *Job) IsTerminal() bool {
	return j.Status == JobStatusCompleted || 
	       j.Status == JobStatusFailed || 
	       j.Status == JobStatusCancelled
}

// CanRetry verifica se o job pode ser tentado novamente
func (j *Job) CanRetry() bool {
	return j.Status == JobStatusFailed && j.Attempts < j.MaxAttempts
}

// MarkStarted marca o job como iniciado
func (j *Job) MarkStarted() {
	j.Status = JobStatusRunning
	now := time.Now()
	j.StartedAt = &now
	j.Attempts++
}

// MarkCompleted marca o job como concluído com sucesso
func (j *Job) MarkCompleted(result map[string]interface{}) {
	j.Status = JobStatusCompleted
	j.Result = &result
	now := time.Now()
	j.CompletedAt = &now
}

// MarkFailed marca o job como falho
func (j *Job) MarkFailed(err error) {
	j.Status = JobStatusFailed
	j.ErrorMsg = err.Error()
	now := time.Now()
	j.CompletedAt = &now
}
