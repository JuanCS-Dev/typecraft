package renderer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LatexRenderer lida com renderização de LaTeX para PDF
type LatexRenderer struct {
	lualatexPath string
	pdflatexPath string
}

// NewLatexRenderer cria uma nova instância do renderizador
func NewLatexRenderer() (*LatexRenderer, error) {
	// Preferir LuaLaTeX (melhor suporte a Unicode)
	lualatexPath, _ := exec.LookPath("lualatex")
	pdflatexPath, _ := exec.LookPath("pdflatex")
	
	if lualatexPath == "" && pdflatexPath == "" {
		return nil, fmt.Errorf("nenhum engine LaTeX encontrado (lualatex ou pdflatex)")
	}
	
	return &LatexRenderer{
		lualatexPath: lualatexPath,
		pdflatexPath: pdflatexPath,
	}, nil
}

// RenderRequest representa uma requisição de renderização
type RenderRequest struct {
	TexFile    string
	OutputDir  string
	Engine     string // "lualatex" ou "pdflatex"
	Runs       int    // Número de execuções (para resolver referências)
	Options    []string
}

// Render executa a renderização LaTeX
func (r *LatexRenderer) Render(req RenderRequest) (string, error) {
	// Verificar se arquivo .tex existe
	if _, err := os.Stat(req.TexFile); os.IsNotExist(err) {
		return "", fmt.Errorf("arquivo .tex não existe: %s", req.TexFile)
	}
	
	// Determinar engine
	engine := req.Engine
	if engine == "" {
		if r.lualatexPath != "" {
			engine = "lualatex"
		} else {
			engine = "pdflatex"
		}
	}
	
	var enginePath string
	switch engine {
	case "lualatex":
		enginePath = r.lualatexPath
	case "pdflatex":
		enginePath = r.pdflatexPath
	default:
		return "", fmt.Errorf("engine desconhecido: %s", engine)
	}
	
	if enginePath == "" {
		return "", fmt.Errorf("engine %s não disponível", engine)
	}
	
	// Determinar diretório de saída
	outputDir := req.OutputDir
	if outputDir == "" {
		outputDir = filepath.Dir(req.TexFile)
	}
	
	// Criar diretório de saída se não existir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de saída: %w", err)
	}
	
	// Número de execuções (padrão: 2 para resolver referências)
	runs := req.Runs
	if runs == 0 {
		runs = 2
	}
	
	// Construir argumentos base
	baseArgs := []string{
		"-interaction=nonstopmode",
		"-halt-on-error",
		"-output-directory=" + outputDir,
	}
	
	// Adicionar opções extras
	baseArgs = append(baseArgs, req.Options...)
	baseArgs = append(baseArgs, req.TexFile)
	
	// Executar múltiplas vezes
	var lastStderr bytes.Buffer
	for i := 0; i < runs; i++ {
		cmd := exec.Command(enginePath, baseArgs...)
		
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		
		if err := cmd.Run(); err != nil {
			lastStderr = stderr
			if i == runs-1 {
				return "", fmt.Errorf("erro ao executar %s (run %d/%d): %w\nStderr: %s", 
					engine, i+1, runs, err, stderr.String())
			}
		}
	}
	
	// Determinar nome do PDF gerado
	baseName := strings.TrimSuffix(filepath.Base(req.TexFile), ".tex")
	pdfPath := filepath.Join(outputDir, baseName+".pdf")
	
	// Verificar se PDF foi criado
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF não foi gerado: %s\nStderr: %s", 
			pdfPath, lastStderr.String())
	}
	
	return pdfPath, nil
}

// RenderToPDF renderiza um arquivo .tex para PDF (método simplificado)
func (r *LatexRenderer) RenderToPDF(texPath, outputDir string) (string, error) {
	return r.Render(RenderRequest{
		TexFile:   texPath,
		OutputDir: outputDir,
		Engine:    "lualatex",
		Runs:      2,
	})
}

// RenderWithOptions renderiza com opções customizadas
func (r *LatexRenderer) RenderWithOptions(texPath, outputDir string, options []string) (string, error) {
	return r.Render(RenderRequest{
		TexFile:   texPath,
		OutputDir: outputDir,
		Engine:    "lualatex",
		Runs:      2,
		Options:   options,
	})
}

// CleanAuxFiles remove arquivos auxiliares gerados pelo LaTeX
func (r *LatexRenderer) CleanAuxFiles(texPath string) error {
	dir := filepath.Dir(texPath)
	baseName := strings.TrimSuffix(filepath.Base(texPath), ".tex")
	
	// Extensões de arquivos auxiliares
	auxExtensions := []string{
		".aux", ".log", ".out", ".toc",
		".lof", ".lot", ".fls", ".fdb_latexmk",
		".synctex.gz", ".blg", ".bbl",
	}
	
	for _, ext := range auxExtensions {
		auxFile := filepath.Join(dir, baseName+ext)
		if _, err := os.Stat(auxFile); err == nil {
			os.Remove(auxFile)
		}
	}
	
	return nil
}

// GetVersion retorna a versão do engine LaTeX
func (r *LatexRenderer) GetVersion() (string, error) {
	enginePath := r.lualatexPath
	if enginePath == "" {
		enginePath = r.pdflatexPath
	}
	
	cmd := exec.Command(enginePath, "--version")
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("erro ao obter versão do LaTeX: %w", err)
	}
	
	// Primeira linha contém a versão
	version := strings.Split(stdout.String(), "\n")[0]
	return version, nil
}

// HasEngine verifica se um engine específico está disponível
func (r *LatexRenderer) HasEngine(engine string) bool {
	switch engine {
	case "lualatex":
		return r.lualatexPath != ""
	case "pdflatex":
		return r.pdflatexPath != ""
	default:
		return false
	}
}
