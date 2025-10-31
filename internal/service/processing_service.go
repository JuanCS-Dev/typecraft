package service

import (
	"fmt"
	"path/filepath"

	"github.com/JuanCS-Dev/typecraft/pkg/converter"
	"github.com/JuanCS-Dev/typecraft/pkg/renderer"
)

// ProcessingService lida com o pipeline de processamento de livros
type ProcessingService struct {
	pandoc *converter.PandocConverter
	latex  *renderer.LatexRenderer
}

// NewProcessingService cria uma nova instância do serviço
func NewProcessingService() (*ProcessingService, error) {
	pandoc, err := converter.NewPandocConverter()
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar Pandoc: %w", err)
	}
	
	latex, err := renderer.NewLatexRenderer()
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar LaTeX: %w", err)
	}
	
	return &ProcessingService{
		pandoc: pandoc,
		latex:  latex,
	}, nil
}

// ConvertManuscript converte o manuscrito para o formato intermediário (Markdown)
func (s *ProcessingService) ConvertManuscript(inputPath, outputDir string) (string, error) {
	// Determinar extensão do arquivo
	ext := filepath.Ext(inputPath)
	baseName := filepath.Base(inputPath[:len(inputPath)-len(ext)])
	outputPath := filepath.Join(outputDir, baseName+".md")
	
	switch ext {
	case ".docx":
		err := s.pandoc.DocxToMarkdown(inputPath, outputPath)
		if err != nil {
			return "", fmt.Errorf("erro ao converter DOCX: %w", err)
		}
	case ".md", ".markdown":
		// Já está em markdown, apenas copiar
		outputPath = inputPath
	default:
		return "", fmt.Errorf("formato não suportado: %s", ext)
	}
	
	return outputPath, nil
}

// GeneratePDF gera PDF a partir de Markdown
func (s *ProcessingService) GeneratePDF(markdownPath, outputPath string, options PDFOptions) error {
	// Construir opções do Pandoc
	pandocOptions := []string{
		"--pdf-engine=lualatex",
	}
	
	// Configurações de página
	if options.PageSize != "" {
		pandocOptions = append(pandocOptions, "-V", "geometry:"+options.PageSize)
	}
	
	// Margem
	if options.Margin != "" {
		pandocOptions = append(pandocOptions, "-V", "geometry:margin="+options.Margin)
	}
	
	// Fonte
	if options.FontFamily != "" {
		pandocOptions = append(pandocOptions, "-V", "mainfont="+options.FontFamily)
	}
	
	if options.FontSize != "" {
		pandocOptions = append(pandocOptions, "-V", "fontsize="+options.FontSize)
	}
	
	// Outras opções
	if options.TOC {
		pandocOptions = append(pandocOptions, "--toc")
	}
	
	if options.NumberSections {
		pandocOptions = append(pandocOptions, "--number-sections")
	}
	
	// Executar conversão
	err := s.pandoc.Convert(converter.ConvertRequest{
		InputFile:  markdownPath,
		OutputFile: outputPath,
		FromFormat: "markdown",
		ToFormat:   "pdf",
		Options:    pandocOptions,
	})
	
	if err != nil {
		return fmt.Errorf("erro ao gerar PDF: %w", err)
	}
	
	return nil
}

// PDFOptions define opções de renderização do PDF
type PDFOptions struct {
	PageSize       string // "letterpaper", "a4paper", etc
	Margin         string // "1in", "2cm", etc
	FontFamily     string // "Times New Roman", "Arial", etc
	FontSize       string // "10pt", "11pt", "12pt"
	TOC            bool   // Índice
	NumberSections bool   // Numerar seções
}

// DefaultPDFOptions retorna opções padrão para KDP (6x9 inches)
func DefaultPDFOptions() PDFOptions {
	return PDFOptions{
		PageSize:       "paperwidth=6in,paperheight=9in",
		Margin:         "0.75in",
		FontFamily:     "", // Usar fonte padrão do sistema
		FontSize:       "11pt",
		TOC:            true,
		NumberSections: true,
	}
}

// IngramSparkPDFOptions retorna opções para IngramSpark
func IngramSparkPDFOptions() PDFOptions {
	return PDFOptions{
		PageSize:       "paperwidth=6in,paperheight=9in",
		Margin:         "0.5in",
		FontFamily:     "", // Usar fonte padrão do sistema
		FontSize:       "11pt",
		TOC:            true,
		NumberSections: true,
	}
}

// ProcessFullPipeline executa o pipeline completo: conversão + renderização
func (s *ProcessingService) ProcessFullPipeline(inputPath, outputDir string, options PDFOptions) (string, error) {
	// 1. Converter para Markdown
	markdownPath, err := s.ConvertManuscript(inputPath, outputDir)
	if err != nil {
		return "", fmt.Errorf("erro na conversão: %w", err)
	}
	
	// 2. Gerar PDF
	baseName := filepath.Base(markdownPath[:len(markdownPath)-3])
	pdfPath := filepath.Join(outputDir, baseName+".pdf")
	
	err = s.GeneratePDF(markdownPath, pdfPath, options)
	if err != nil {
		return "", fmt.Errorf("erro na renderização: %w", err)
	}
	
	return pdfPath, nil
}
