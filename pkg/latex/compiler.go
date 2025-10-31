package latex

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Compiler compila documentos LaTeX em PDF
type Compiler struct {
	workDir     string
	keepTemp    bool
	timeout     time.Duration
	engine      string // pdflatex, xelatex, lualatex
	extraArgs   []string
}

// CompilerOption configura o compilador
type CompilerOption func(*Compiler)

// NewCompiler cria um novo compilador LaTeX
func NewCompiler(opts ...CompilerOption) (*Compiler, error) {
	workDir, err := os.MkdirTemp("", "typecraft-latex-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	c := &Compiler{
		workDir:  workDir,
		keepTemp: false,
		timeout:  30 * time.Second,
		engine:   "pdflatex",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// WithKeepTemp mantém arquivos temporários para debug
func WithKeepTemp(keep bool) CompilerOption {
	return func(c *Compiler) {
		c.keepTemp = keep
	}
}

// WithTimeout define timeout para compilação
func WithTimeout(timeout time.Duration) CompilerOption {
	return func(c *Compiler) {
		c.timeout = timeout
	}
}

// WithEngine define engine LaTeX (pdflatex, xelatex, lualatex)
func WithEngine(engine string) CompilerOption {
	return func(c *Compiler) {
		c.engine = engine
	}
}

// WithExtraArgs adiciona argumentos extras ao compilador
func WithExtraArgs(args ...string) CompilerOption {
	return func(c *Compiler) {
		c.extraArgs = append(c.extraArgs, args...)
	}
}

// CompileResult resultado da compilação
type CompileResult struct {
	Success    bool
	PDFPath    string
	LogPath    string
	Errors     []CompileError
	Warnings   []string
	Duration   time.Duration
	TempFiles  []string
}

// CompileError erro de compilação LaTeX
type CompileError struct {
	Line    int
	File    string
	Message string
	Type    string // error, warning, fatal
}

// Compile compila documento LaTeX para PDF
func (c *Compiler) Compile(latexContent string) (*CompileResult, error) {
	start := time.Now()
	
	// Cria arquivo .tex
	texPath := filepath.Join(c.workDir, "document.tex")
	if err := os.WriteFile(texPath, []byte(latexContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write tex file: %w", err)
	}

	result := &CompileResult{
		TempFiles: []string{texPath},
	}

	// Compila (2 passes para resolver referências)
	for i := 0; i < 2; i++ {
		if err := c.runCompiler(texPath); err != nil {
			result.Success = false
			result.Duration = time.Since(start)
			
			// Parse errors from log
			logPath := filepath.Join(c.workDir, "document.log")
			if logData, err := os.ReadFile(logPath); err == nil {
				result.LogPath = logPath
				result.Errors = c.parseErrors(string(logData))
				result.Warnings = c.parseWarnings(string(logData))
			}
			
			return result, fmt.Errorf("compilation failed: %w", err)
		}
	}

	// Verifica se PDF foi gerado
	pdfPath := filepath.Join(c.workDir, "document.pdf")
	if _, err := os.Stat(pdfPath); err != nil {
		result.Success = false
		result.Duration = time.Since(start)
		return result, fmt.Errorf("PDF not generated: %w", err)
	}

	result.Success = true
	result.PDFPath = pdfPath
	result.Duration = time.Since(start)
	
	// Lista arquivos temp
	files, _ := filepath.Glob(filepath.Join(c.workDir, "document.*"))
	result.TempFiles = files

	// Parse log para warnings
	logPath := filepath.Join(c.workDir, "document.log")
	if logData, err := os.ReadFile(logPath); err == nil {
		result.LogPath = logPath
		result.Warnings = c.parseWarnings(string(logData))
	}

	return result, nil
}

// runCompiler executa o compilador LaTeX
func (c *Compiler) runCompiler(texPath string) error {
	args := []string{
		"-interaction=nonstopmode",
		"-halt-on-error",
		"-file-line-error",
		"-output-directory=" + c.workDir,
	}
	
	args = append(args, c.extraArgs...)
	args = append(args, texPath)

	cmd := exec.Command(c.engine, args...)
	cmd.Dir = c.workDir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", c.engine, err)
	}

	// Timeout handling
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(c.timeout):
		cmd.Process.Kill()
		return fmt.Errorf("compilation timeout after %v", c.timeout)
	case err := <-done:
		if err != nil {
			return fmt.Errorf("compilation error: %w\n%s", err, stderr.String())
		}
	}

	return nil
}

// parseErrors extrai erros do log LaTeX
func (c *Compiler) parseErrors(log string) []CompileError {
	var errors []CompileError

	// Pattern: ./document.tex:10: Error message
	reError := regexp.MustCompile(`(?m)^\./(.*?):(\d+):\s*(.*)$`)
	matches := reError.FindAllStringSubmatch(log, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			line := 0
			fmt.Sscanf(match[2], "%d", &line)
			
			errors = append(errors, CompileError{
				File:    match[1],
				Line:    line,
				Message: strings.TrimSpace(match[3]),
				Type:    "error",
			})
		}
	}

	// Pattern: ! Error message
	reFatal := regexp.MustCompile(`(?m)^!\s+(.*)$`)
	matches = reFatal.FindAllStringSubmatch(log, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			errors = append(errors, CompileError{
				Message: strings.TrimSpace(match[1]),
				Type:    "fatal",
			})
		}
	}

	return errors
}

// parseWarnings extrai warnings do log LaTeX
func (c *Compiler) parseWarnings(log string) []string {
	var warnings []string

	// Pattern: LaTeX Warning: ...
	reWarning := regexp.MustCompile(`(?m)^LaTeX Warning:\s*(.*)$`)
	matches := reWarning.FindAllStringSubmatch(log, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			warnings = append(warnings, strings.TrimSpace(match[1]))
		}
	}

	// Pattern: Package warning
	rePackage := regexp.MustCompile(`(?m)^Package \w+ Warning:\s*(.*)$`)
	matches = rePackage.FindAllStringSubmatch(log, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			warnings = append(warnings, strings.TrimSpace(match[1]))
		}
	}

	return warnings
}

// CopyPDF copia PDF gerado para destino
func (c *Compiler) CopyPDF(result *CompileResult, destPath string) error {
	if !result.Success {
		return fmt.Errorf("cannot copy PDF: compilation failed")
	}

	data, err := os.ReadFile(result.PDFPath)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}

// Cleanup remove arquivos temporários
func (c *Compiler) Cleanup() error {
	if c.keepTemp {
		return nil
	}

	if c.workDir != "" {
		return os.RemoveAll(c.workDir)
	}

	return nil
}

// GetWorkDir retorna diretório de trabalho
func (c *Compiler) GetWorkDir() string {
	return c.workDir
}

// ValidatePDF valida PDF gerado
func ValidatePDF(pdfPath string) error {
	// Verifica se arquivo existe
	info, err := os.Stat(pdfPath)
	if err != nil {
		return fmt.Errorf("PDF not found: %w", err)
	}

	// Verifica tamanho mínimo (PDF vazio ~500 bytes)
	if info.Size() < 500 {
		return fmt.Errorf("PDF too small: %d bytes", info.Size())
	}

	// Verifica header PDF
	data := make([]byte, 5)
	f, err := os.Open(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	if _, err := f.Read(data); err != nil {
		return fmt.Errorf("failed to read PDF header: %w", err)
	}

	if string(data) != "%PDF-" {
		return fmt.Errorf("invalid PDF header: %s", string(data))
	}

	return nil
}
