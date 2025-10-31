// Package e2e provides end-to-end tests for the Typecraft pipeline
// Conformidade: Constituição Vértice v3.0 - P2 (Validação Preventiva)
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/pipeline/html"
)

func TestFullPipelineHTML(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Check pagedjs-cli availability
	if _, err := exec.LookPath("pagedjs-cli"); err != nil {
		t.Skip("pagedjs-cli not found. Run: npm install -g pagedjs-cli")
	}

	// Setup
	fixtureDir := "fixtures"
	outputDir := t.TempDir()
	
	inputFile := filepath.Join(fixtureDir, "sample_book.md")
	htmlFile := filepath.Join(outputDir, "book.html")
	pdfFile := filepath.Join(outputDir, "book.pdf")

	// Test data
	metadata := html.Metadata{
		Title:    "O Algoritmo do Destino",
		Author:   "Dr. Marcus Chen",
		Language: "pt-BR",
		Subject:  "Ficção Científica",
		Keywords: "IA, consciência, ética, futuro",
	}

	// Step 1: Convert Markdown to HTML
	t.Run("ConvertMarkdownToHTML", func(t *testing.T) {
		converter := html.NewConverter()
		err := converter.ConvertFile(inputFile, htmlFile, metadata)
		if err != nil {
			t.Fatalf("Failed to convert Markdown to HTML: %v", err)
		}

		// Verify HTML file was created
		if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
			t.Fatal("HTML file was not created")
		}

		// Check file size
		info, err := os.Stat(htmlFile)
		if err != nil {
			t.Fatalf("Failed to stat HTML file: %v", err)
		}
		if info.Size() == 0 {
			t.Fatal("HTML file is empty")
		}

		t.Logf("✅ HTML created: %s (%d bytes)", htmlFile, info.Size())
	})

	// Step 2: Render HTML to PDF with Paged.js
	t.Run("RenderHTMLToPDF", func(t *testing.T) {
		renderer := html.NewPagedJSRenderer()
		err := renderer.RenderPDF(htmlFile, pdfFile)
		if err != nil {
			t.Fatalf("Failed to render PDF: %v", err)
		}

		// Verify PDF file was created
		if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
			t.Fatal("PDF file was not created")
		}

		// Check file size
		info, err := os.Stat(pdfFile)
		if err != nil {
			t.Fatalf("Failed to stat PDF file: %v", err)
		}
		if info.Size() == 0 {
			t.Fatal("PDF file is empty")
		}
		if info.Size() < 1000 {
			t.Fatalf("PDF file seems too small: %d bytes", info.Size())
		}

		t.Logf("✅ PDF created: %s (%d bytes)", pdfFile, info.Size())
	})

	// Step 3: Validate PDF structure
	t.Run("ValidatePDF", func(t *testing.T) {
		// Read PDF file header
		file, err := os.Open(pdfFile)
		if err != nil {
			t.Fatalf("Failed to open PDF: %v", err)
		}
		defer file.Close()

		header := make([]byte, 5)
		n, err := file.Read(header)
		if err != nil {
			t.Fatalf("Failed to read PDF header: %v", err)
		}
		if n != 5 {
			t.Fatalf("Expected to read 5 bytes, got %d", n)
		}

		// Check PDF magic number
		if string(header) != "%PDF-" {
			t.Fatalf("Invalid PDF header: %s", string(header))
		}

		t.Logf("✅ PDF is valid")
	})
}

func TestFullPipelineWithTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	if _, err := exec.LookPath("pagedjs-cli"); err != nil {
		t.Skip("pagedjs-cli not found")
	}

	// Setup
	fixtureDir := "fixtures"
	outputDir := t.TempDir()
	
	inputFile := filepath.Join(fixtureDir, "sample_book.md")
	htmlFile := filepath.Join(outputDir, "book.html")
	pdfFile := filepath.Join(outputDir, "book.pdf")

	metadata := html.Metadata{
		Title:    "O Algoritmo do Destino",
		Author:   "Dr. Marcus Chen",
		Language: "pt-BR",
		Subject:  "Ficção Científica",
		Keywords: "IA, consciência, ética, futuro",
	}

	// Custom CSS variables
	cssVars := map[string]string{
		"--primary-color":   "#2C3E50",
		"--secondary-color": "#34495E",
		"--accent-color":    "#E74C3C",
		"--body-font":       "Source Serif Pro, serif",
		"--heading-font":    "Playfair Display, serif",
	}

	// Convert with custom template
	t.Run("ConvertWithCustomCSS", func(t *testing.T) {
		converter := html.NewConverter()
		converter.SetCSSVariables(cssVars)
		
		err := converter.ConvertFile(inputFile, htmlFile, metadata)
		if err != nil {
			t.Fatalf("Failed to convert: %v", err)
		}

		// Verify HTML contains custom CSS
		content, err := os.ReadFile(htmlFile)
		if err != nil {
			t.Fatalf("Failed to read HTML: %v", err)
		}

		htmlStr := string(content)
		if !strings.Contains(htmlStr, "--primary-color") {
			t.Error("HTML doesn't contain custom CSS variables")
		}

		t.Logf("✅ Custom CSS applied")
	})

	// Render to PDF
	t.Run("RenderCustomPDF", func(t *testing.T) {
		renderer := html.NewPagedJSRenderer()
		err := renderer.RenderPDF(htmlFile, pdfFile)
		if err != nil {
			t.Fatalf("Failed to render PDF: %v", err)
		}

		info, err := os.Stat(pdfFile)
		if err != nil {
			t.Fatalf("Failed to stat PDF: %v", err)
		}

		t.Logf("✅ Custom PDF created: %d bytes", info.Size())
	})
}

func TestPerformanceBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping benchmark in short mode")
	}

	if _, err := exec.LookPath("pagedjs-cli"); err != nil {
		t.Skip("pagedjs-cli not found")
	}

	fixtureDir := "fixtures"
	inputFile := filepath.Join(fixtureDir, "sample_book.md")

	metadata := html.Metadata{
		Title:  "Benchmark Test",
		Author: "Test Author",
	}

	// Benchmark HTML conversion
	t.Run("BenchmarkHTMLConversion", func(t *testing.T) {
		outputDir := t.TempDir()
		htmlFile := filepath.Join(outputDir, "benchmark.html")

		start := time.Now()
		converter := html.NewConverter()
		err := converter.ConvertFile(inputFile, htmlFile, metadata)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Conversion failed: %v", err)
		}

		t.Logf("⏱️  HTML conversion: %v", elapsed)

		// Performance assertion
		if elapsed > 5*time.Second {
			t.Errorf("HTML conversion too slow: %v (expected < 5s)", elapsed)
		}
	})

	// Benchmark PDF rendering
	t.Run("BenchmarkPDFRendering", func(t *testing.T) {
		outputDir := t.TempDir()
		htmlFile := filepath.Join(outputDir, "benchmark.html")
		pdfFile := filepath.Join(outputDir, "benchmark.pdf")

		// First create HTML
		converter := html.NewConverter()
		err := converter.ConvertFile(inputFile, htmlFile, metadata)
		if err != nil {
			t.Fatalf("HTML conversion failed: %v", err)
		}

		// Then benchmark PDF
		start := time.Now()
		renderer := html.NewPagedJSRenderer()
		err = renderer.RenderPDF(htmlFile, pdfFile)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("PDF rendering failed: %v", err)
		}

		t.Logf("⏱️  PDF rendering: %v", elapsed)

		// Performance assertion
		if elapsed > 30*time.Second {
			t.Errorf("PDF rendering too slow: %v (expected < 30s)", elapsed)
		}
	})
}
