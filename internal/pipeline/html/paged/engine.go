package paged

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Engine gerencia a conversão HTML para PDF usando Paged.js
// Conformidade: CONSTITUIÇÃO VÉRTICE v3.0
type Engine struct {
	nodeModulesPath string
	tempDir         string
}

// Config configurações do engine Paged.js
type Config struct {
	// Caminho para node_modules (opcional, usa detecção automática)
	NodeModulesPath string
	// Diretório temporário para arquivos intermediários
	TempDir string
}

// PagedOutput resultado da conversão
type PagedOutput struct {
	PDFPath   string
	HTMLPath  string
	PageCount int
	FileSize  int64
	Warnings  []string
}

// NewEngine cria uma nova instância do Paged.js engine
func NewEngine(cfg Config) (*Engine, error) {
	// Detectar node_modules se não especificado
	nodeModules := cfg.NodeModulesPath
	if nodeModules == "" {
		detected, err := detectNodeModules()
		if err != nil {
			return nil, fmt.Errorf("node_modules não encontrado: %w", err)
		}
		nodeModules = detected
	}

	// Validar que pagedjs existe
	pagedPath := filepath.Join(nodeModules, "pagedjs", "package.json")
	if _, err := os.Stat(pagedPath); err != nil {
		return nil, fmt.Errorf("pagedjs não instalado em %s: %w", nodeModules, err)
	}

	// Setup temp dir
	tempDir := cfg.TempDir
	if tempDir == "" {
		tempDir = filepath.Join(os.TempDir(), "typecraft-paged")
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("falha criar temp dir: %w", err)
	}

	log.Info().
		Str("node_modules", nodeModules).
		Str("temp_dir", tempDir).
		Msg("Paged.js engine inicializado")

	return &Engine{
		nodeModulesPath: nodeModules,
		tempDir:         tempDir,
	}, nil
}

// ConvertToPDF converte HTML para PDF usando Paged.js
func (e *Engine) ConvertToPDF(ctx context.Context, htmlContent string, options PageOptions) (*PagedOutput, error) {
	// Criar arquivo HTML temporário
	htmlFile := filepath.Join(e.tempDir, fmt.Sprintf("input-%s.html", generateID()))
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
		return nil, fmt.Errorf("falha escrever HTML: %w", err)
	}
	defer os.Remove(htmlFile)

	// Preparar arquivo CSS se fornecido
	if options.CustomCSS != "" {
		cssFile := filepath.Join(e.tempDir, "custom.css")
		if err := os.WriteFile(cssFile, []byte(options.CustomCSS), 0644); err != nil {
			return nil, fmt.Errorf("falha escrever CSS: %w", err)
		}
		defer os.Remove(cssFile)
		htmlContent = injectCSSLink(htmlContent, cssFile)
		if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
			return nil, fmt.Errorf("falha atualizar HTML com CSS: %w", err)
		}
	}

	// Arquivo de saída
	outputPDF := filepath.Join(e.tempDir, fmt.Sprintf("output-%s.pdf", generateID()))

	// Executar pagedjs-cli
	if err := e.runPagedJS(ctx, htmlFile, outputPDF, options); err != nil {
		return nil, fmt.Errorf("falha executar pagedjs: %w", err)
	}

	// Obter informações do arquivo
	stat, err := os.Stat(outputPDF)
	if err != nil {
		return nil, fmt.Errorf("falha obter info PDF: %w", err)
	}

	result := &PagedOutput{
		PDFPath:  outputPDF,
		HTMLPath: htmlFile,
		FileSize: stat.Size(),
		Warnings: []string{},
	}

	log.Info().
		Str("pdf", outputPDF).
		Int64("size", result.FileSize).
		Msg("PDF gerado com sucesso")

	return result, nil
}

// runPagedJS executa o comando pagedjs-cli
func (e *Engine) runPagedJS(ctx context.Context, inputHTML, outputPDF string, options PageOptions) error {
	// Construir comando
	// Usando npx para executar pagedjs-cli
	args := []string{
		"pagedjs-cli",
		inputHTML,
		"-o", outputPDF,
	}

	// Adicionar opções
	if options.Format != "" {
		args = append(args, "--format", options.Format)
	}
	if options.Landscape {
		args = append(args, "--landscape")
	}
	if options.Timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", options.Timeout))
	}

	cmd := exec.CommandContext(ctx, "npx", args...)
	cmd.Dir = filepath.Dir(e.nodeModulesPath)

	// Capturar output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Falha executar pagedjs-cli")
		return fmt.Errorf("pagedjs-cli error: %w\nOutput: %s", err, output)
	}

	log.Debug().
		Str("output", string(output)).
		Msg("pagedjs-cli executado")

	return nil
}

// detectNodeModules tenta encontrar node_modules no projeto
func detectNodeModules() (string, error) {
	// Procurar node_modules a partir do diretório atual
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Procurar até 3 níveis acima
	for i := 0; i < 3; i++ {
		nodeModules := filepath.Join(cwd, "node_modules")
		if _, err := os.Stat(nodeModules); err == nil {
			return nodeModules, nil
		}
		cwd = filepath.Dir(cwd)
	}

	return "", fmt.Errorf("node_modules não encontrado")
}

// injectCSSLink injeta link CSS no HTML
func injectCSSLink(html, cssPath string) string {
	cssLink := fmt.Sprintf(`<link rel="stylesheet" href="%s">`, cssPath)
	// Inserir antes de </head>
	if strings.Contains(html, "</head>") {
		return strings.Replace(html, "</head>", cssLink+"\n</head>", 1)
	}
	// Se não tem head, adicionar no início
	return cssLink + "\n" + html
}

// generateID gera um ID único simples
func generateID() string {
	return fmt.Sprintf("%d", os.Getpid())
}

// Cleanup remove arquivos temporários
func (e *Engine) Cleanup() error {
	return os.RemoveAll(e.tempDir)
}
