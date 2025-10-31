package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JuanCS-Dev/typecraft/pkg/ai"
	"github.com/JuanCS-Dev/typecraft/pkg/typography"
)

// Pipeline orquestra todo o processo de geração do livro
type Pipeline struct {
	htmlGen      *HTMLGenerator
	pdfGen       *PDFGenerator
	fontSubset   *FontSubsetter
	styleEngine  *typography.StyleEngine
	aiClient     *ai.Client
	outputDir    string
}

// NewPipeline cria uma nova pipeline de processamento
func NewPipeline(styleEngine *typography.StyleEngine, aiClient *ai.Client, outputDir string) (*Pipeline, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de saída: %w", err)
	}

	pdfGen, err := NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar gerador PDF: %w", err)
	}

	fontSubset, err := NewFontSubsetter()
	if err != nil {
		// Font subsetting é opcional
		fmt.Println("⚠️  Font subsetting não disponível (requer fonttools)")
		fontSubset = nil
	}

	return &Pipeline{
		htmlGen:     NewHTMLGenerator(styleEngine, aiClient),
		pdfGen:      pdfGen,
		fontSubset:  fontSubset,
		styleEngine: styleEngine,
		aiClient:    aiClient,
		outputDir:   outputDir,
	}, nil
}

// ProcessBookConfig configuração do livro
type ProcessBookConfig struct {
	InputFiles   []string
	Title        string
	Author       string
	DesignPrompt string
	PageSize     string
	FontDir      string
}

// ProcessBook processa o livro completo
func (p *Pipeline) ProcessBook(config ProcessBookConfig) (*ProcessResult, error) {
	fmt.Println("🚀 Iniciando pipeline de processamento...")
	startTime := time.Now()

	result := &ProcessResult{
		Steps: make(map[string]StepResult),
	}

	// 1. Carregar e processar capítulos
	fmt.Println("📖 Processando capítulos...")
	sections, err := p.loadAndProcessChapters(config.InputFiles)
	if err != nil {
		return nil, fmt.Errorf("erro ao processar capítulos: %w", err)
	}
	result.Steps["chapters"] = StepResult{Success: true, Duration: time.Since(startTime)}

	// 2. Gerar HTML base
	fmt.Println("🎨 Gerando HTML estruturado...")
	metadata := map[string]interface{}{
		"title":  config.Title,
		"author": config.Author,
	}

	html, err := p.htmlGen.GeneratePagedJS(sections, metadata)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar HTML: %w", err)
	}

	// 3. Aplicar design system com IA (se configurado)
	if config.DesignPrompt != "" && p.aiClient != nil {
		fmt.Println("✨ Gerando design system com IA...")
		html, err = p.htmlGen.ApplyDesignSystem(html, config.DesignPrompt)
		if err != nil {
			fmt.Printf("⚠️  Erro ao gerar design system: %v\n", err)
		} else {
			result.Steps["design"] = StepResult{Success: true, Duration: time.Since(startTime)}
		}
	}

	// 4. Salvar HTML
	htmlPath := filepath.Join(p.outputDir, "book.html")
	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("erro ao salvar HTML: %w", err)
	}
	result.HTMLPath = htmlPath
	fmt.Printf("✅ HTML salvo em: %s\n", htmlPath)

	// 5. Font subsetting (se disponível e configurado)
	if p.fontSubset != nil && config.FontDir != "" {
		fmt.Println("🔤 Otimizando fontes...")
		fontOutputDir := filepath.Join(p.outputDir, "fonts")
		text := p.fontSubset.ExtractTextFromHTML(html)
		
		if err := p.fontSubset.SubsetFontFamily(config.FontDir, text, fontOutputDir); err != nil {
			fmt.Printf("⚠️  Erro ao otimizar fontes: %v\n", err)
		} else {
			// Gera CSS para as fontes
			fontCSS, err := p.fontSubset.GenerateFontFaceCSS(fontOutputDir, "BookFont")
			if err == nil {
				cssPath := filepath.Join(fontOutputDir, "fonts.css")
				os.WriteFile(cssPath, []byte(fontCSS), 0644)
				fmt.Printf("✅ Fontes otimizadas em: %s\n", fontOutputDir)
				result.Steps["fonts"] = StepResult{Success: true, Duration: time.Since(startTime)}
			}
		}
	}

	// 6. Gerar PDF
	fmt.Println("📄 Gerando PDF...")
	pdfPath := filepath.Join(p.outputDir, "book.pdf")
	pdfOpts := PDFOptions{
		PageSize:    config.PageSize,
		OutlineTags: "h1,h2",
	}

	if err := p.pdfGen.GeneratePDFWithOptions(html, pdfPath, pdfOpts); err != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", err)
	}
	result.PDFPath = pdfPath
	result.Steps["pdf"] = StepResult{Success: true, Duration: time.Since(startTime)}
	fmt.Printf("✅ PDF gerado em: %s\n", pdfPath)

	result.Success = true
	result.TotalDuration = time.Since(startTime)

	return result, nil
}

// loadAndProcessChapters carrega e processa arquivos de entrada
func (p *Pipeline) loadAndProcessChapters(inputFiles []string) ([]BookSection, error) {
	var sections []BookSection

	for i, file := range inputFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler %s: %w", file, err)
		}

		section, err := p.htmlGen.ProcessChapter(string(content), i+1)
		if err != nil {
			return nil, fmt.Errorf("erro ao processar %s: %w", file, err)
		}

		// Define título baseado no nome do arquivo se não tiver
		if section.Title == "" {
			baseName := filepath.Base(file)
			section.Title = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		}

		sections = append(sections, section)
	}

	return sections, nil
}

// Cleanup limpa recursos
func (p *Pipeline) Cleanup() error {
	if p.pdfGen != nil {
		return p.pdfGen.Cleanup()
	}
	return nil
}

// ProcessResult resultado do processamento
type ProcessResult struct {
	Success       bool
	HTMLPath      string
	PDFPath       string
	Steps         map[string]StepResult
	TotalDuration time.Duration
}

// StepResult resultado de uma etapa
type StepResult struct {
	Success  bool
	Duration time.Duration
	Error    error
}

// Report gera relatório do processamento
func (r *ProcessResult) Report() string {
	var report strings.Builder
	
	report.WriteString("\n═══════════════════════════════════════════\n")
	report.WriteString("       📚 RELATÓRIO DE PROCESSAMENTO\n")
	report.WriteString("═══════════════════════════════════════════\n\n")

	if r.Success {
		report.WriteString("✅ Status: SUCESSO\n\n")
	} else {
		report.WriteString("❌ Status: FALHA\n\n")
	}

	report.WriteString("📄 Arquivos gerados:\n")
	if r.HTMLPath != "" {
		report.WriteString(fmt.Sprintf("   HTML: %s\n", r.HTMLPath))
	}
	if r.PDFPath != "" {
		report.WriteString(fmt.Sprintf("   PDF:  %s\n", r.PDFPath))
	}

	report.WriteString("\n⏱️  Etapas executadas:\n")
	for step, result := range r.Steps {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		report.WriteString(fmt.Sprintf("   %s %-12s: %v\n", status, step, result.Duration.Round(time.Millisecond)))
	}

	report.WriteString(fmt.Sprintf("\n🎯 Tempo total: %v\n", r.TotalDuration.Round(time.Millisecond)))
	report.WriteString("\n═══════════════════════════════════════════\n")

	return report.String()
}
