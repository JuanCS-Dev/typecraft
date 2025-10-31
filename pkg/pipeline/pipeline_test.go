package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JuanCS-Dev/typecraft/pkg/typography"
)

func TestHTMLGenerator(t *testing.T) {
	styleEngine := typography.NewStyleEngine()
	gen := NewHTMLGenerator(styleEngine, nil)

	sections := []BookSection{
		{
			Title:   "Capítulo 1",
			Content: "<p>Conteúdo do primeiro capítulo.</p>",
			Type:    "chapter",
			Number:  1,
		},
	}

	metadata := map[string]interface{}{
		"title":  "Livro de Teste",
		"author": "Autor Teste",
	}

	html, err := gen.GenerateHTML(sections, metadata)
	if err != nil {
		t.Fatalf("Erro ao gerar HTML: %v", err)
	}

	if html == "" {
		t.Error("HTML gerado está vazio")
	}

	if !contains(html, "Livro de Teste") {
		t.Error("HTML não contém o título do livro")
	}

	if !contains(html, "Capítulo 1") {
		t.Error("HTML não contém o título do capítulo")
	}
}

func TestProcessChapter(t *testing.T) {
	styleEngine := typography.NewStyleEngine()
	gen := NewHTMLGenerator(styleEngine, nil)

	rawContent := `Este é um teste de conteúdo.

Este é outro parágrafo.`

	section, err := gen.ProcessChapter(rawContent, 1)
	if err != nil {
		t.Fatalf("Erro ao processar capítulo: %v", err)
	}

	if section.Number != 1 {
		t.Errorf("Número do capítulo incorreto: got %d, want 1", section.Number)
	}

	if section.Type != "chapter" {
		t.Errorf("Tipo incorreto: got %s, want chapter", section.Type)
	}

	if !contains(section.Content, "<p>") {
		t.Error("Conteúdo não foi convertido para HTML")
	}
}

func TestGeneratePagedJS(t *testing.T) {
	styleEngine := typography.NewStyleEngine()
	gen := NewHTMLGenerator(styleEngine, nil)

	sections := []BookSection{
		{
			Title:   "Teste",
			Content: "<p>Conteúdo</p>",
			Type:    "chapter",
			Number:  1,
		},
	}

	metadata := map[string]interface{}{
		"title":  "Teste",
		"author": "Autor",
	}

	html, err := gen.GeneratePagedJS(sections, metadata)
	if err != nil {
		t.Fatalf("Erro ao gerar HTML com Paged.js: %v", err)
	}

	if !contains(html, "pagedjs") {
		t.Error("HTML não contém referência ao Paged.js")
	}

	if !contains(html, "paged.polyfill.js") {
		t.Error("HTML não contém o polyfill do Paged.js")
	}
}

func TestPipeline_ProcessBook(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração em modo short")
	}

	// Setup
	tmpDir, err := os.MkdirTemp("", "typecraft-test-*")
	if err != nil {
		t.Fatalf("Erro ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Cria arquivo de teste
	testFile := filepath.Join(tmpDir, "capitulo1.txt")
	testContent := "Este é um capítulo de teste.\n\nCom múltiplos parágrafos."
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Erro ao criar arquivo de teste: %v", err)
	}

	// Cria pipeline
	styleEngine := typography.NewStyleEngine()
	outputDir := filepath.Join(tmpDir, "output")
	
	pipeline, err := NewPipeline(styleEngine, nil, outputDir)
	if err != nil {
		t.Fatalf("Erro ao criar pipeline: %v", err)
	}
	defer pipeline.Cleanup()

	// Configuração
	config := ProcessBookConfig{
		InputFiles: []string{testFile},
		Title:      "Livro de Teste",
		Author:     "Autor de Teste",
		PageSize:   "A5",
	}

	// Processa
	result, err := pipeline.ProcessBook(config)
	
	// Se pagedjs-cli não estiver disponível, aceita o erro
	if err != nil && contains(err.Error(), "pagedjs-cli") {
		t.Skip("pagedjs-cli não disponível, pulando teste")
		return
	}
	
	if err != nil {
		t.Fatalf("Erro ao processar livro: %v", err)
	}

	// Verificações
	if !result.Success {
		t.Error("Processamento não foi bem sucedido")
	}

	if result.HTMLPath == "" {
		t.Error("Caminho do HTML não foi definido")
	}

	// Verifica se HTML foi criado
	if _, err := os.Stat(result.HTMLPath); os.IsNotExist(err) {
		t.Error("Arquivo HTML não foi criado")
	}

	// Lê e verifica conteúdo do HTML
	htmlContent, err := os.ReadFile(result.HTMLPath)
	if err != nil {
		t.Fatalf("Erro ao ler HTML: %v", err)
	}

	if !contains(string(htmlContent), "Livro de Teste") {
		t.Error("HTML não contém título do livro")
	}

	// Imprime relatório
	t.Log(result.Report())
}

func TestFontSubsetter_UniqueRunes(t *testing.T) {
	input := "aabbccddee"
	expected := "abcde"
	
	result := uniqueRunes(input)
	
	if len(result) != len(expected) {
		t.Errorf("Tamanho incorreto: got %d, want %d", len(result), len(expected))
	}

	for _, r := range expected {
		if !contains(result, string(r)) {
			t.Errorf("Caractere %c não encontrado no resultado", r)
		}
	}
}

func TestFontSubsetter_ExtractTextFromHTML(t *testing.T) {
	subsetter := &FontSubsetter{}
	
	html := "<html><body><p>Teste</p><div>Mais texto</div></body></html>"
	
	text := subsetter.ExtractTextFromHTML(html)
	
	if !contains(text, "Teste") {
		t.Error("Texto não foi extraído corretamente")
	}
	
	if contains(text, "<p>") {
		t.Error("Tags HTML não foram removidas")
	}
}

func BenchmarkProcessChapter(b *testing.B) {
	styleEngine := typography.NewStyleEngine()
	gen := NewHTMLGenerator(styleEngine, nil)
	
	content := "Este é um conteúdo de teste.\n\nCom múltiplos parágrafos."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gen.ProcessChapter(content, 1)
	}
}

func BenchmarkGenerateHTML(b *testing.B) {
	styleEngine := typography.NewStyleEngine()
	gen := NewHTMLGenerator(styleEngine, nil)
	
	sections := []BookSection{
		{
			Title:   "Capítulo 1",
			Content: "<p>Conteúdo</p>",
			Type:    "chapter",
			Number:  1,
		},
	}
	
	metadata := map[string]interface{}{
		"title":  "Teste",
		"author": "Autor",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = gen.GenerateHTML(sections, metadata)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
