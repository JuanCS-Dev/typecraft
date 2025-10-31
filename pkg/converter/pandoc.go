package converter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PandocConverter lida com conversões usando Pandoc
type PandocConverter struct {
	pandocPath string
}

// NewPandocConverter cria uma nova instância do conversor
func NewPandocConverter() (*PandocConverter, error) {
	// Verificar se pandoc está disponível
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, fmt.Errorf("pandoc não encontrado no PATH: %w", err)
	}
	
	return &PandocConverter{
		pandocPath: pandocPath,
	}, nil
}

// ConvertRequest representa uma requisição de conversão
type ConvertRequest struct {
	InputFile  string
	OutputFile string
	FromFormat string
	ToFormat   string
	Options    []string
}

// Convert executa a conversão usando pandoc
func (c *PandocConverter) Convert(req ConvertRequest) error {
	// Verificar se arquivo de entrada existe
	if _, err := os.Stat(req.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("arquivo de entrada não existe: %s", req.InputFile)
	}
	
	// Criar diretório de saída se não existir
	outputDir := filepath.Dir(req.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de saída: %w", err)
	}
	
	// Construir argumentos do comando
	args := []string{
		req.InputFile,
		"-o", req.OutputFile,
	}
	
	if req.FromFormat != "" {
		args = append(args, "-f", req.FromFormat)
	}
	
	if req.ToFormat != "" {
		args = append(args, "-t", req.ToFormat)
	}
	
	// Adicionar opções extras
	args = append(args, req.Options...)
	
	// Executar comando
	cmd := exec.Command(c.pandocPath, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar pandoc: %w\nStderr: %s", err, stderr.String())
	}
	
	// Verificar se arquivo de saída foi criado
	if _, err := os.Stat(req.OutputFile); os.IsNotExist(err) {
		return fmt.Errorf("arquivo de saída não foi criado")
	}
	
	return nil
}

// DocxToMarkdown converte DOCX para Markdown
func (c *PandocConverter) DocxToMarkdown(inputPath, outputPath string) error {
	return c.Convert(ConvertRequest{
		InputFile:  inputPath,
		OutputFile: outputPath,
		FromFormat: "docx",
		ToFormat:   "markdown",
		Options: []string{
			"--extract-media=.",
			"--wrap=none",
		},
	})
}

// MarkdownToLatex converte Markdown para LaTeX
func (c *PandocConverter) MarkdownToLatex(inputPath, outputPath string, templatePath string) error {
	options := []string{
		"--standalone",
		"--toc",
		"--number-sections",
	}
	
	if templatePath != "" {
		options = append(options, "--template="+templatePath)
	}
	
	return c.Convert(ConvertRequest{
		InputFile:  inputPath,
		OutputFile: outputPath,
		FromFormat: "markdown",
		ToFormat:   "latex",
		Options:    options,
	})
}

// MarkdownToPDF converte Markdown diretamente para PDF (usando LaTeX internamente)
func (c *PandocConverter) MarkdownToPDF(inputPath, outputPath string, options []string) error {
	defaultOptions := []string{
		"--pdf-engine=lualatex",
		"--standalone",
		"--toc",
		"--number-sections",
	}
	
	// Merge com opções customizadas
	allOptions := append(defaultOptions, options...)
	
	return c.Convert(ConvertRequest{
		InputFile:  inputPath,
		OutputFile: outputPath,
		FromFormat: "markdown",
		ToFormat:   "pdf",
		Options:    allOptions,
	})
}

// DocxToPDF converte DOCX diretamente para PDF
func (c *PandocConverter) DocxToPDF(inputPath, outputPath string, options []string) error {
	defaultOptions := []string{
		"--pdf-engine=lualatex",
		"--standalone",
		"--toc",
		"--number-sections",
	}
	
	allOptions := append(defaultOptions, options...)
	
	return c.Convert(ConvertRequest{
		InputFile:  inputPath,
		OutputFile: outputPath,
		FromFormat: "docx",
		ToFormat:   "pdf",
		Options:    allOptions,
	})
}

// GetVersion retorna a versão do Pandoc instalada
func (c *PandocConverter) GetVersion() (string, error) {
	cmd := exec.Command(c.pandocPath, "--version")
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("erro ao obter versão do pandoc: %w", err)
	}
	
	// Primeira linha contém a versão
	version := strings.Split(stdout.String(), "\n")[0]
	return version, nil
}

// ListInputFormats retorna os formatos de entrada suportados
func (c *PandocConverter) ListInputFormats() ([]string, error) {
	cmd := exec.Command(c.pandocPath, "--list-input-formats")
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao listar formatos de entrada: %w", err)
	}
	
	formats := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	return formats, nil
}

// ListOutputFormats retorna os formatos de saída suportados
func (c *PandocConverter) ListOutputFormats() ([]string, error) {
	cmd := exec.Command(c.pandocPath, "--list-output-formats")
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao listar formatos de saída: %w", err)
	}
	
	formats := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	return formats, nil
}
