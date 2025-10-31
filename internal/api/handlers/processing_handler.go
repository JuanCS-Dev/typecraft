package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/JuanCS-Dev/typecraft/internal/service"
	"github.com/gin-gonic/gin"
)

// ProcessingHandler lida com requisições de processamento
type ProcessingHandler struct {
	service *service.ProcessingService
}

// NewProcessingHandler cria uma nova instância do handler
func NewProcessingHandler() (*ProcessingHandler, error) {
	processingService, err := service.NewProcessingService()
	if err != nil {
		return nil, err
	}
	
	return &ProcessingHandler{
		service: processingService,
	}, nil
}

// ConvertFile converte um arquivo para outro formato
func (h *ProcessingHandler) ConvertFile(c *gin.Context) {
	// Receber arquivo
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo é obrigatório"})
		return
	}
	
	// Criar diretório temporário
	tempDir := filepath.Join(os.TempDir(), "typecraft", c.GetString("request_id"))
	os.MkdirAll(tempDir, 0755)
	
	// Salvar arquivo temporariamente
	inputPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar arquivo"})
		return
	}
	
	// Converter
	outputPath, err := h.service.ConvertManuscript(inputPath, tempDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Retornar arquivo convertido
	c.FileAttachment(outputPath, filepath.Base(outputPath))
	
	// Limpar depois (em background)
	go func() {
		os.RemoveAll(tempDir)
	}()
}

// GeneratePDF gera PDF a partir de arquivo
func (h *ProcessingHandler) GeneratePDF(c *gin.Context) {
	// Receber arquivo
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo é obrigatório"})
		return
	}
	
	// Opções de PDF
	var pdfOptions service.PDFOptions
	if err := c.ShouldBind(&pdfOptions); err != nil {
		// Usar defaults
		pdfOptions = service.DefaultPDFOptions()
	}
	
	// Criar diretório temporário
	tempDir := filepath.Join(os.TempDir(), "typecraft", fmt.Sprintf("pdf_%d", os.Getpid()))
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	
	// Salvar arquivo temporariamente
	inputPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar arquivo"})
		return
	}
	
	// Processar pipeline completo
	pdfPath, err := h.service.ProcessFullPipeline(inputPath, tempDir, pdfOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Retornar PDF
	c.FileAttachment(pdfPath, filepath.Base(pdfPath))
}

// ProcessManuscript processa um manuscrito completo (endpoint simplificado)
func (h *ProcessingHandler) ProcessManuscript(c *gin.Context) {
	// Receber arquivo
	file, err := c.FormFile("manuscript")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo 'manuscript' é obrigatório"})
		return
	}
	
	// Parâmetros
	format := c.DefaultPostForm("format", "kdp") // kdp ou ingramspark
	
	// Criar diretório temporário
	tempDir := filepath.Join(os.TempDir(), "typecraft", fmt.Sprintf("process_%d", os.Getpid()))
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	
	// Salvar arquivo
	inputPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar arquivo"})
		return
	}
	
	// Escolher opções conforme formato
	var pdfOptions service.PDFOptions
	if format == "ingramspark" {
		pdfOptions = service.IngramSparkPDFOptions()
	} else {
		pdfOptions = service.DefaultPDFOptions()
	}
	
	// Processar
	pdfPath, err := h.service.ProcessFullPipeline(inputPath, tempDir, pdfOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Retornar PDF
	c.FileAttachment(pdfPath, filepath.Base(pdfPath))
}
