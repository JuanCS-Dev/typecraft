package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PDFGenerator gera PDFs usando Paged.js CLI
type PDFGenerator struct {
	tempDir string
}

// NewPDFGenerator cria um novo gerador de PDF
func NewPDFGenerator() (*PDFGenerator, error) {
	tempDir, err := os.MkdirTemp("", "typecraft-pdf-*")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório temporário: %w", err)
	}

	return &PDFGenerator{
		tempDir: tempDir,
	}, nil
}

// GeneratePDF converte HTML em PDF usando pagedjs-cli
func (p *PDFGenerator) GeneratePDF(htmlContent string, outputPath string) error {
	// Cria arquivo HTML temporário
	htmlPath := filepath.Join(p.tempDir, "input.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("erro ao escrever HTML temporário: %w", err)
	}

	// Executa pagedjs-cli
	cmd := exec.Command("pagedjs-cli",
		htmlPath,
		"-o", outputPath,
		"--timeout", "120000",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao executar pagedjs-cli: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GeneratePDFWithOptions gera PDF com opções customizadas
func (p *PDFGenerator) GeneratePDFWithOptions(htmlContent string, outputPath string, opts PDFOptions) error {
	htmlPath := filepath.Join(p.tempDir, "input.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("erro ao escrever HTML temporário: %w", err)
	}

	args := []string{
		htmlPath,
		"-o", outputPath,
		"--timeout", "120000",
	}

	// Adiciona opções
	if opts.PageSize != "" {
		args = append(args, "-s", opts.PageSize)
	}

	if opts.Landscape {
		args = append(args, "-l")
	}

	if opts.OutlineTags != "" {
		args = append(args, "--outline-tags", opts.OutlineTags)
	}

	if len(opts.AdditionalStyles) > 0 {
		for _, style := range opts.AdditionalStyles {
			args = append(args, "--style", style)
		}
	}

	cmd := exec.Command("pagedjs-cli", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao executar pagedjs-cli: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Cleanup remove arquivos temporários
func (p *PDFGenerator) Cleanup() error {
	if p.tempDir != "" {
		return os.RemoveAll(p.tempDir)
	}
	return nil
}

// PDFOptions opções de geração de PDF
type PDFOptions struct {
	PageSize         string   // A4, A5, Letter, etc.
	Landscape        bool
	OutlineTags      string   // "h1,h2" para gerar índice
	AdditionalStyles []string // Caminhos para CSS adicional
}
