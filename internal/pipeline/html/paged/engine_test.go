package paged

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestPagedEngine testa o engine básico
func TestPagedEngine(t *testing.T) {
	// Skip se não estiver em ambiente com Node.js
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Pulando teste de integração")
	}

	engine, err := NewEngine(Config{})
	if err != nil {
		t.Fatalf("Falha criar engine: %v", err)
	}
	defer engine.Cleanup()

	// HTML simples de teste
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Test</title>
</head>
<body>
    <h1>Test Document</h1>
    <p>This is a test paragraph.</p>
</body>
</html>`

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	options := DefaultPageOptions()
	result, err := engine.ConvertToPDF(ctx, html, options)
	if err != nil {
		t.Fatalf("Falha converter para PDF: %v", err)
	}

	// Verificar que PDF foi criado
	if _, err := os.Stat(result.PDFPath); err != nil {
		t.Errorf("PDF não foi criado: %v", err)
	}

	// Verificar tamanho
	if result.FileSize == 0 {
		t.Error("PDF está vazio")
	}

	t.Logf("PDF gerado: %s (%d bytes)", result.PDFPath, result.FileSize)
}

// TestBookTemplate testa renderização de template
func TestBookTemplate(t *testing.T) {
	tmpl := BookTemplate{
		Title:      "Livro de Teste",
		Author:     "Autor Teste",
		Publisher:  "Editora Teste",
		ISBN:       "123-456-789",
		Copyright:  "© 2025 Todos os direitos reservados",
		Content:    "<h1>Capítulo 1</h1><p>Conteúdo do capítulo.</p>",
		FontFamily: "Georgia, serif",
		FontSize:   "12pt",
	}

	html, err := tmpl.Render()
	if err != nil {
		t.Fatalf("Falha renderizar template: %v", err)
	}

	// Verificar elementos essenciais
	if len(html) == 0 {
		t.Error("HTML está vazio")
	}

	// Verificar presença de elementos
	checks := []string{
		"Livro de Teste",
		"Autor Teste",
		"Capítulo 1",
		"DOCTYPE html",
		"UTF-8",
	}

	for _, check := range checks {
		if !contains(html, check) {
			t.Errorf("HTML não contém: %s", check)
		}
	}

	t.Logf("Template renderizado com %d bytes", len(html))
}

// TestPageOptions testa geração de CSS
func TestPageOptions(t *testing.T) {
	opts := DefaultPageOptions()
	css := opts.GeneratePagedCSS()

	if len(css) == 0 {
		t.Error("CSS está vazio")
	}

	// Verificar elementos CSS essenciais
	checks := []string{
		"@page",
		"margin-top",
		"counter(page",
	}

	for _, check := range checks {
		if !contains(css, check) {
			t.Errorf("CSS não contém: %s", check)
		}
	}

	t.Logf("CSS gerado:\n%s", css)
}

// TestMinimalTemplate testa template mínimo
func TestMinimalTemplate(t *testing.T) {
	content := "<h1>Test</h1><p>Content</p>"
	html := MinimalTemplate(content)

	if len(html) == 0 {
		t.Error("HTML está vazio")
	}

	if !contains(html, content) {
		t.Error("HTML não contém o conteúdo fornecido")
	}

	t.Logf("Template mínimo: %d bytes", len(html))
}

// Helper: contains verifica se string contém substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 len(substr) == 0 || 
		 findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
