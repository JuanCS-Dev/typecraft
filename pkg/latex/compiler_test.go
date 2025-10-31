package latex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCompiler(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	assert.NotEmpty(t, c.workDir)
	assert.Equal(t, "pdflatex", c.engine)
	assert.False(t, c.keepTemp)
	assert.Equal(t, 30*time.Second, c.timeout)

	// Verifica se diretório foi criado
	_, err = os.Stat(c.workDir)
	assert.NoError(t, err)
}

func TestCompiler_Options(t *testing.T) {
	c, err := NewCompiler(
		WithEngine("xelatex"),
		WithTimeout(60*time.Second),
		WithKeepTemp(true),
		WithExtraArgs("-shell-escape"),
	)
	require.NoError(t, err)
	defer c.Cleanup()

	assert.Equal(t, "xelatex", c.engine)
	assert.Equal(t, 60*time.Second, c.timeout)
	assert.True(t, c.keepTemp)
	assert.Contains(t, c.extraArgs, "-shell-escape")
}

func TestCompiler_CompileSimple(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	latex := `\documentclass{article}
\begin{document}
Hello, World!
\end{document}`

	result, err := c.Compile(latex)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.PDFPath)
	assert.Greater(t, result.Duration.Milliseconds(), int64(0))

	// Verifica se PDF existe
	_, err = os.Stat(result.PDFPath)
	assert.NoError(t, err)
}

func TestCompiler_CompileWithMetadata(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	latex := `\documentclass{article}
\usepackage[utf8]{inputenc}
\title{Test Document}
\author{Test Author}
\date{\today}

\begin{document}
\maketitle

\section{Introduction}
This is a test document.

\end{document}`

	result, err := c.Compile(latex)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.PDFPath)

	// Valida PDF
	err = ValidatePDF(result.PDFPath)
	assert.NoError(t, err)
}

func TestCompiler_CompileWithSections(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	doc := NewDocument(ClassArticle).
		AddOption("12pt").
		AddPackage("inputenc", "utf8").
		SetTitle("Test").
		SetAuthor("Author").
		AddSection("Introduction", "This is the intro.").
		AddSection("Methods", "This is the methods.").
		AddSection("Results", "This is the results.")

	latex := doc.Generate()

	result, err := c.Compile(latex)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.PDFPath)

	// Verifica tamanho mínimo
	info, err := os.Stat(result.PDFPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(1000))
}

func TestCompiler_CompileError(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	// LaTeX inválido
	latex := `\documentclass{article}
\begin{document}
\undefined_command
\end{document}`

	result, err := c.Compile(latex)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Errors)
}

func TestCompiler_ParseErrors(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	log := `./document.tex:5: Undefined control sequence
l.5 \undefined
        _command
! Emergency stop.`

	errors := c.parseErrors(log)
	require.NotEmpty(t, errors)

	// Verifica erro de linha
	found := false
	for _, e := range errors {
		if e.Line == 5 && strings.Contains(e.Message, "Undefined control sequence") {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected error at line 5")

	// Verifica erro fatal
	foundFatal := false
	for _, e := range errors {
		if e.Type == "fatal" && strings.Contains(e.Message, "Emergency stop") {
			foundFatal = true
			break
		}
	}
	assert.True(t, foundFatal, "Expected fatal error")
}

func TestCompiler_ParseWarnings(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	log := `LaTeX Warning: Reference 'fig:test' on page 1 undefined.
Package hyperref Warning: Token not allowed in a PDF string.`

	warnings := c.parseWarnings(log)
	require.Len(t, warnings, 2)

	assert.Contains(t, warnings[0], "Reference")
	assert.Contains(t, warnings[1], "Token not allowed")
}

func TestCompiler_CopyPDF(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	latex := `\documentclass{article}
\begin{document}
Test
\end{document}`

	result, err := c.Compile(latex)
	require.NoError(t, err)

	// Copia para destino
	destPath := filepath.Join(os.TempDir(), "test-output.pdf")
	defer os.Remove(destPath)

	err = c.CopyPDF(result, destPath)
	require.NoError(t, err)

	// Verifica se foi copiado
	_, err = os.Stat(destPath)
	assert.NoError(t, err)

	// Valida PDF copiado
	err = ValidatePDF(destPath)
	assert.NoError(t, err)
}

func TestCompiler_Cleanup(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)

	workDir := c.workDir

	// Verifica que existe
	_, err = os.Stat(workDir)
	require.NoError(t, err)

	// Cleanup
	err = c.Cleanup()
	require.NoError(t, err)

	// Verifica que foi removido
	_, err = os.Stat(workDir)
	assert.True(t, os.IsNotExist(err))
}

func TestCompiler_KeepTemp(t *testing.T) {
	c, err := NewCompiler(WithKeepTemp(true))
	require.NoError(t, err)

	workDir := c.workDir

	latex := `\documentclass{article}
\begin{document}
Test
\end{document}`

	_, err = c.Compile(latex)
	require.NoError(t, err)

	// Cleanup
	err = c.Cleanup()
	require.NoError(t, err)

	// Verifica que ainda existe
	_, err = os.Stat(workDir)
	assert.NoError(t, err)

	// Limpa manualmente
	os.RemoveAll(workDir)
}

func TestValidatePDF(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "valid PDF",
			setup: func() string {
				c, _ := NewCompiler()
				defer c.Cleanup()

				latex := `\documentclass{article}
\begin{document}
Test
\end{document}`

				result, _ := c.Compile(latex)
				return result.PDFPath
			},
			wantErr: false,
		},
		{
			name: "non-existent file",
			setup: func() string {
				return "/tmp/non-existent-file.pdf"
			},
			wantErr: true,
		},
		{
			name: "invalid PDF header",
			setup: func() string {
				path := filepath.Join(os.TempDir(), "invalid.pdf")
				os.WriteFile(path, []byte("not a PDF"), 0644)
				return path
			},
			wantErr: true,
		},
		{
			name: "too small",
			setup: func() string {
				path := filepath.Join(os.TempDir(), "small.pdf")
				os.WriteFile(path, []byte("%PDF-"), 0644)
				return path
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer os.Remove(path)

			err := ValidatePDF(path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompiler_WithTemplate(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	// Usa template
	tmpl := NewTemplate(`\documentclass{article}
\begin{document}
\title{ {{title}} }
\author{ {{author}} }
\maketitle

\section{Introduction}
{{content}}

\end{document}`)

	tmpl.SetVariable("title", "Test Document")
	tmpl.SetVariable("author", "Test Author")
	tmpl.SetVariable("content", "This is the content.")

	latex, err := tmpl.Render()
	require.NoError(t, err)

	result, err := c.Compile(latex)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.PDFPath)
}

func TestCompiler_MultipleCompilations(t *testing.T) {
	c, err := NewCompiler()
	require.NoError(t, err)
	defer c.Cleanup()

	// Primeira compilação
	latex1 := `\documentclass{article}
\begin{document}
Document 1
\end{document}`

	result1, err := c.Compile(latex1)
	require.NoError(t, err)
	assert.True(t, result1.Success)

	// Segunda compilação (deve sobrescrever)
	latex2 := `\documentclass{article}
\begin{document}
Document 2
\end{document}`

	result2, err := c.Compile(latex2)
	require.NoError(t, err)
	assert.True(t, result2.Success)

	// Ambos devem ter mesmo workDir
	assert.Equal(t, filepath.Dir(result1.PDFPath), filepath.Dir(result2.PDFPath))
}

// Benchmarks
func BenchmarkCompiler_Compile(b *testing.B) {
	c, err := NewCompiler()
	require.NoError(b, err)
	defer c.Cleanup()

	latex := `\documentclass{article}
\begin{document}
Hello, World!
\end{document}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Compile(latex)
	}
}

func BenchmarkValidatePDF(b *testing.B) {
	c, err := NewCompiler()
	require.NoError(b, err)
	defer c.Cleanup()

	latex := `\documentclass{article}
\begin{document}
Test
\end{document}`

	result, err := c.Compile(latex)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePDF(result.PDFPath)
	}
}
