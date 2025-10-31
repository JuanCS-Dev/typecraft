package latex

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument(ClassArticle)
	require.NotNil(t, doc)
	assert.Equal(t, ClassArticle, doc.Class)
	assert.Empty(t, doc.Options)
	assert.Empty(t, doc.Packages)
	assert.Empty(t, doc.Content)
}

func TestDocument_AddOption(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddOption("12pt").AddOption("a4paper")
	
	assert.Len(t, doc.Options, 2)
	assert.Contains(t, doc.Options, "12pt")
	assert.Contains(t, doc.Options, "a4paper")
}

func TestDocument_AddPackage(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddPackage("geometry", "margin=1in")
	doc.AddPackage("graphicx")
	
	assert.Len(t, doc.Packages, 2)
	assert.Equal(t, "geometry", doc.Packages[0].Name)
	assert.Equal(t, []string{"margin=1in"}, doc.Packages[0].Options)
	assert.Equal(t, "graphicx", doc.Packages[1].Name)
	assert.Empty(t, doc.Packages[1].Options)
}

func TestDocument_SetMetadata(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.SetMetadata(DocumentMetadata{
		Title:  "Test Document",
		Author: "John Doe",
		Date:   "2025-01-01",
	})
	
	assert.Equal(t, "Test Document", doc.Metadata.Title)
	assert.Equal(t, "John Doe", doc.Metadata.Author)
	assert.Equal(t, "2025-01-01", doc.Metadata.Date)
}

func TestDocument_AddContent(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddContent("Paragraph 1")
	doc.AddContent("Paragraph 2")
	
	assert.Len(t, doc.Content, 2)
	assert.Equal(t, "Paragraph 1", doc.Content[0])
	assert.Equal(t, "Paragraph 2", doc.Content[1])
}

func TestDocument_AddSection(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddSection("Introduction", "This is the introduction.")
	
	assert.Len(t, doc.Content, 1)
	assert.Contains(t, doc.Content[0], "\\section{Introduction}")
	assert.Contains(t, doc.Content[0], "This is the introduction")
}

func TestDocument_AddSubsection(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddSubsection("Background", "Some background info.")
	
	assert.Len(t, doc.Content, 1)
	assert.Contains(t, doc.Content[0], "\\subsection{Background}")
	assert.Contains(t, doc.Content[0], "Some background info")
}

func TestDocument_AddChapter(t *testing.T) {
	t.Run("book class", func(t *testing.T) {
		doc := NewDocument(ClassBook)
		doc.AddChapter("Chapter 1", "Content of chapter 1")
		
		assert.Len(t, doc.Content, 1)
		assert.Contains(t, doc.Content[0], "\\chapter{Chapter 1}")
	})
	
	t.Run("article class ignores chapters", func(t *testing.T) {
		doc := NewDocument(ClassArticle)
		doc.AddChapter("Chapter 1", "Content")
		
		assert.Empty(t, doc.Content)
	})
}

func TestDocument_Generate_Minimal(t *testing.T) {
	doc := NewDocument(ClassArticle)
	result := doc.Generate()
	
	assert.Contains(t, result, "\\documentclass{article}")
	assert.Contains(t, result, "\\begin{document}")
	assert.Contains(t, result, "\\end{document}")
}

func TestDocument_Generate_WithOptions(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddOption("12pt")
	doc.AddOption("a4paper")
	
	result := doc.Generate()
	
	assert.Contains(t, result, "\\documentclass[12pt,a4paper]{article}")
}

func TestDocument_Generate_WithPackages(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddPackage("inputenc", "utf8")
	doc.AddPackage("graphicx")
	
	result := doc.Generate()
	
	assert.Contains(t, result, "\\usepackage[utf8]{inputenc}")
	assert.Contains(t, result, "\\usepackage{graphicx}")
}

func TestDocument_Generate_WithMetadata(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.SetMetadata(DocumentMetadata{
		Title:  "My Document",
		Author: "Jane Smith",
		Date:   "\\today",
	})
	
	result := doc.Generate()
	
	assert.Contains(t, result, "\\title{My Document}")
	assert.Contains(t, result, "\\author{Jane Smith}")
	assert.Contains(t, result, "\\date{\\today}")
	assert.Contains(t, result, "\\maketitle")
}

func TestDocument_Generate_WithContent(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddContent("First paragraph.")
	doc.AddSection("Introduction", "Intro text.")
	doc.AddContent("Final paragraph.")
	
	result := doc.Generate()
	
	assert.Contains(t, result, "First paragraph")
	assert.Contains(t, result, "\\section{Introduction}")
	assert.Contains(t, result, "Intro text")
	assert.Contains(t, result, "Final paragraph")
}

func TestDocument_Generate_Complete(t *testing.T) {
	doc := NewDocument(ClassArticle)
	doc.AddOption("12pt")
	doc.AddPackage("inputenc", "utf8")
	doc.AddPackage("geometry", "margin=1in")
	doc.SetMetadata(DocumentMetadata{
		Title:  "Complete Document",
		Author: "Test Author",
		Date:   "2025-01-01",
	})
	doc.AddSection("Introduction", "Introduction text.")
	doc.AddSection("Methods", "Methods text.")
	doc.AddSection("Results", "Results text.")
	doc.AddSection("Conclusion", "Conclusion text.")
	
	result := doc.Generate()
	
	// Check structure
	assert.Contains(t, result, "\\documentclass[12pt]{article}")
	assert.Contains(t, result, "\\usepackage[utf8]{inputenc}")
	assert.Contains(t, result, "\\usepackage[margin=1in]{geometry}")
	assert.Contains(t, result, "\\title{Complete Document}")
	assert.Contains(t, result, "\\author{Test Author}")
	assert.Contains(t, result, "\\date{2025-01-01}")
	assert.Contains(t, result, "\\maketitle")
	assert.Contains(t, result, "\\section{Introduction}")
	assert.Contains(t, result, "\\section{Methods}")
	assert.Contains(t, result, "\\section{Results}")
	assert.Contains(t, result, "\\section{Conclusion}")
	assert.Contains(t, result, "\\end{document}")
	
	// Check order
	titleIdx := strings.Index(result, "\\title{")
	beginIdx := strings.Index(result, "\\begin{document}")
	introIdx := strings.Index(result, "\\section{Introduction}")
	endIdx := strings.Index(result, "\\end{document}")
	
	assert.True(t, titleIdx < beginIdx)
	assert.True(t, beginIdx < introIdx)
	assert.True(t, introIdx < endIdx)
}

func TestDocumentBuilder(t *testing.T) {
	doc := NewDocumentBuilder(ClassBook).
		WithOptions("11pt", "twoside").
		WithPackage("inputenc", "utf8").
		WithPackage("graphicx").
		WithMetadata(DocumentMetadata{
			Title:  "My Book",
			Author: "Author Name",
		}).
		WithChapter("Chapter 1", "Content of chapter 1").
		WithSection("Section 1.1", "Section content").
		WithContent("Additional content").
		Build()
	
	result := doc.Generate()
	
	assert.Contains(t, result, "\\documentclass[11pt,twoside]{book}")
	assert.Contains(t, result, "\\usepackage[utf8]{inputenc}")
	assert.Contains(t, result, "\\title{My Book}")
	assert.Contains(t, result, "\\chapter{Chapter 1}")
	assert.Contains(t, result, "\\section{Section 1.1}")
	assert.Contains(t, result, "Additional content")
}

func TestStandardPackages(t *testing.T) {
	packages := StandardPackages()
	
	assert.NotEmpty(t, packages)
	
	names := make([]string, 0, len(packages))
	for _, pkg := range packages {
		names = append(names, pkg.Name)
	}
	
	assert.Contains(t, names, "inputenc")
	assert.Contains(t, names, "fontenc")
	assert.Contains(t, names, "babel")
	assert.Contains(t, names, "geometry")
	assert.Contains(t, names, "graphicx")
	assert.Contains(t, names, "hyperref")
	assert.Contains(t, names, "amsmath")
}

func TestAcademicPackages(t *testing.T) {
	packages := AcademicPackages()
	
	assert.NotEmpty(t, packages)
	
	names := make([]string, 0, len(packages))
	for _, pkg := range packages {
		names = append(names, pkg.Name)
	}
	
	// Should include standard packages
	assert.Contains(t, names, "inputenc")
	assert.Contains(t, names, "graphicx")
	
	// Should include academic-specific packages
	assert.Contains(t, names, "natbib")
	assert.Contains(t, names, "algorithm")
	assert.Contains(t, names, "listings")
	assert.Contains(t, names, "booktabs")
}

func TestBookPackages(t *testing.T) {
	packages := BookPackages()
	
	assert.NotEmpty(t, packages)
	
	names := make([]string, 0, len(packages))
	for _, pkg := range packages {
		names = append(names, pkg.Name)
	}
	
	// Should include standard packages
	assert.Contains(t, names, "inputenc")
	
	// Should include book-specific packages
	assert.Contains(t, names, "fancyhdr")
	assert.Contains(t, names, "titlesec")
	assert.Contains(t, names, "tocloft")
}

func TestDocumentClass(t *testing.T) {
	tests := []struct {
		class    DocumentClass
		expected string
	}{
		{ClassArticle, "article"},
		{ClassBook, "book"},
		{ClassReport, "report"},
		{ClassBeamer, "beamer"},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			doc := NewDocument(tt.class)
			result := doc.Generate()
			assert.Contains(t, result, fmt.Sprintf("\\documentclass{%s}", tt.expected))
		})
	}
}

func TestDocument_ChainedCalls(t *testing.T) {
	doc := NewDocument(ClassArticle).
		AddOption("12pt").
		AddPackage("inputenc", "utf8").
		SetMetadata(DocumentMetadata{Title: "Test"}).
		AddSection("Section 1", "Content 1").
		AddSection("Section 2", "Content 2")
	
	result := doc.Generate()
	
	assert.Contains(t, result, "12pt")
	assert.Contains(t, result, "inputenc")
	assert.Contains(t, result, "Test")
	assert.Contains(t, result, "Section 1")
	assert.Contains(t, result, "Section 2")
}

func BenchmarkDocument_Generate(b *testing.B) {
	doc := NewDocument(ClassArticle)
	doc.AddOption("12pt")
	doc.AddPackage("inputenc", "utf8")
	doc.AddPackage("geometry", "margin=1in")
	doc.SetMetadata(DocumentMetadata{
		Title:  "Benchmark Document",
		Author: "Bench",
		Date:   "\\today",
	})
	
	for i := 1; i <= 10; i++ {
		doc.AddSection(fmt.Sprintf("Section %d", i), "Content")
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_ = doc.Generate()
	}
}
