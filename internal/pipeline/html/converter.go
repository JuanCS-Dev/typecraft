package html

import (
	"bytes"
	"fmt"
	"os/exec"
)

// PandocConverter gerencia a conversão de Markdown para HTML via Pandoc
type PandocConverter struct {
	PandocPath string // Caminho para o executável pandoc
}

// NewPandocConverter cria um novo conversor
func NewPandocConverter() (*PandocConverter, error) {
	// Verificar se pandoc está instalado
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, fmt.Errorf("pandoc não encontrado: %w (instale com: apt install pandoc)", err)
	}

	return &PandocConverter{
		PandocPath: pandocPath,
	}, nil
}

// ConvertOptions configura a conversão
type ConvertOptions struct {
	InputFormat  string            // Ex: "markdown", "markdown+smart"
	OutputFormat string            // Ex: "html5", "html"
	Standalone   bool              // Gera documento HTML completo
	TOC          bool              // Gera índice (table of contents)
	Template     string            // Template customizado (opcional)
	Variables    map[string]string // Variáveis para injeção
	CSSFiles     []string          // Arquivos CSS a incluir
	Metadata     map[string]string // Metadados YAML adicionais
}

// DefaultHTMLOptions retorna opções padrão para conversão HTML
func DefaultHTMLOptions() ConvertOptions {
	return ConvertOptions{
		InputFormat:  "markdown+smart",
		OutputFormat: "html5",
		Standalone:   true,
		TOC:          true,
		Variables:    make(map[string]string),
		Metadata:     make(map[string]string),
	}
}

// Convert converte Markdown para HTML usando Pandoc
//
// Referência: Blueprint Seção 5.1 - "O Princípio da Fonte Única"
// "O Pandoc lerá o arquivo Markdown com seu cabeçalho YAML e o transformará
// no formato intermediário necessário para a etapa de renderização final"
func (pc *PandocConverter) Convert(markdown string, opts ConvertOptions) (string, error) {
	// Construir argumentos do Pandoc
	args := []string{
		"--from", opts.InputFormat,
		"--to", opts.OutputFormat,
	}

	if opts.Standalone {
		args = append(args, "--standalone")
	}

	if opts.TOC {
		args = append(args, "--toc")
		args = append(args, "--toc-depth=3")
	}

	// Template customizado
	if opts.Template != "" {
		args = append(args, "--template", opts.Template)
	}

	// Variáveis
	for key, value := range opts.Variables {
		args = append(args, "--variable", fmt.Sprintf("%s=%s", key, value))
	}

	// Metadados
	for key, value := range opts.Metadata {
		args = append(args, "--metadata", fmt.Sprintf("%s=%s", key, value))
	}

	// Arquivos CSS
	for _, cssFile := range opts.CSSFiles {
		args = append(args, "--css", cssFile)
	}

	// Criar comando
	cmd := exec.Command(pc.PandocPath, args...)

	// Input via stdin
	cmd.Stdin = bytes.NewBufferString(markdown)

	// Capturar output e error
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Executar
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pandoc falhou: %w\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ConvertToHTML5 é um helper para conversão rápida para HTML5
func (pc *PandocConverter) ConvertToHTML5(markdown string) (string, error) {
	opts := DefaultHTMLOptions()
	return pc.Convert(markdown, opts)
}

// ConvertWithCSS converte e injeta arquivos CSS
func (pc *PandocConverter) ConvertWithCSS(markdown string, cssFiles ...string) (string, error) {
	opts := DefaultHTMLOptions()
	opts.CSSFiles = cssFiles
	return pc.Convert(markdown, opts)
}

// GetVersion retorna a versão do Pandoc instalado
func (pc *PandocConverter) GetVersion() (string, error) {
	cmd := exec.Command(pc.PandocPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
