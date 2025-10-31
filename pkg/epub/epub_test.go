package epub

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEPub(t *testing.T) {
	epub := NewEPub(EPub3)
	require.NotNil(t, epub)
	assert.Equal(t, EPub3, epub.Version)
	assert.Empty(t, epub.Chapters)
	assert.Empty(t, epub.Fonts)
	assert.Empty(t, epub.Images)
}

func TestEPub_AddChapter(t *testing.T) {
	epub := NewEPub(EPub3)
	
	// Add chapter with auto-generated ID and filename
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>This is chapter 1</p>",
	})
	
	assert.Len(t, epub.Chapters, 1)
	assert.Equal(t, "chapter1", epub.Chapters[0].ID)
	assert.Equal(t, "chapter1.xhtml", epub.Chapters[0].FileName)
	
	// Add chapter with custom ID and filename
	epub.AddChapter(Chapter{
		ID:       "intro",
		Title:    "Introduction",
		Content:  "<p>Introduction</p>",
		FileName: "intro.xhtml",
	})
	
	assert.Len(t, epub.Chapters, 2)
	assert.Equal(t, "intro", epub.Chapters[1].ID)
	assert.Equal(t, "intro.xhtml", epub.Chapters[1].FileName)
}

func TestEPub_Write(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test.epub")
	
	// Create ePub
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Test Book",
		Author:     "Test Author",
		Language:   "en",
		Identifier: "test-123",
		Publisher:  "Test Publisher",
		Description: "A test book",
		Subject:    []string{"Fiction", "Test"},
		Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Rights:     "Public Domain",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>This is the first chapter.</p>",
	})
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 2",
		Content: "<p>This is the second chapter.</p>",
	})
	
	// Write ePub
	err := epub.Write(outputPath)
	require.NoError(t, err)
	
	// Verify file exists
	_, err = os.Stat(outputPath)
	require.NoError(t, err)
	
	// Verify it's a valid ZIP
	reader, err := zip.OpenReader(outputPath)
	require.NoError(t, err)
	defer reader.Close()
	
	// Check for required files
	requiredFiles := []string{
		"mimetype",
		"META-INF/container.xml",
		"OEBPS/content.opf",
		"OEBPS/toc.ncx",
		"OEBPS/nav.xhtml",
		"OEBPS/Text/chapter1.xhtml",
		"OEBPS/Text/chapter2.xhtml",
		"OEBPS/Styles/style.css",
	}
	
	foundFiles := make(map[string]bool)
	for _, file := range reader.File {
		foundFiles[file.Name] = true
	}
	
	for _, required := range requiredFiles {
		assert.True(t, foundFiles[required], "Missing required file: %s", required)
	}
	
	// Verify mimetype is first and uncompressed
	assert.Equal(t, "mimetype", reader.File[0].Name)
	assert.Equal(t, zip.Store, reader.File[0].Method, "mimetype should be uncompressed")
	
	// Verify mimetype content
	mimetypeFile, err := reader.File[0].Open()
	require.NoError(t, err)
	defer mimetypeFile.Close()
	
	mimetypeBytes, err := io.ReadAll(mimetypeFile)
	require.NoError(t, err)
	assert.Equal(t, "application/epub+zip", string(mimetypeBytes))
}

func TestEPub_Write_EPUB2(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test2.epub")
	
	epub := NewEPub(EPub2)
	epub.Metadata = Metadata{
		Title:      "Test EPUB 2",
		Author:     "Test Author",
		Language:   "en",
		Identifier: "test-epub2",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>EPUB 2 content</p>",
	})
	
	err := epub.Write(outputPath)
	require.NoError(t, err)
	
	// Verify file exists
	_, err = os.Stat(outputPath)
	require.NoError(t, err)
	
	// EPUB 2 should not have nav.xhtml
	reader, err := zip.OpenReader(outputPath)
	require.NoError(t, err)
	defer reader.Close()
	
	foundNav := false
	for _, file := range reader.File {
		if file.Name == "OEBPS/nav.xhtml" {
			foundNav = true
		}
	}
	
	// EPUB 2 shouldn't have nav.xhtml, but EPUB 3 should
	// Our implementation creates it for EPUB 3 only
	if epub.Version == EPub2 {
		assert.False(t, foundNav, "EPUB 2 should not have nav.xhtml")
	}
}

func TestEPub_WithCustomCSS(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "custom-css.epub")
	
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Custom CSS Test",
		Author:     "Test",
		Language:   "en",
		Identifier: "css-test",
	}
	
	epub.CSS = `
body {
  font-family: Arial, sans-serif;
  color: #333;
}
`
	
	epub.AddChapter(Chapter{
		Title:   "Test",
		Content: "<p>Test</p>",
	})
	
	err := epub.Write(outputPath)
	require.NoError(t, err)
	
	// Verify custom CSS was written
	reader, err := zip.OpenReader(outputPath)
	require.NoError(t, err)
	defer reader.Close()
	
	for _, file := range reader.File {
		if file.Name == "OEBPS/Styles/style.css" {
			f, err := file.Open()
			require.NoError(t, err)
			defer f.Close()
			
			contentBytes, err := io.ReadAll(f)
			require.NoError(t, err)
			
			assert.Contains(t, string(contentBytes), "Arial, sans-serif")
			return
		}
	}
	
	t.Fatal("CSS file not found")
}

func TestEPub_ChapterWrapping(t *testing.T) {
	epub := NewEPub(EPub3)
	
	chapter := Chapter{
		Title:    "Test Chapter",
		Content:  "<p>Test content</p>",
		FileName: "test.xhtml",
	}
	
	wrapped := epub.wrapChapterHTML(chapter)
	
	assert.Contains(t, wrapped, "<?xml version")
	assert.Contains(t, wrapped, "<!DOCTYPE html>")
	assert.Contains(t, wrapped, "<html xmlns=")
	assert.Contains(t, wrapped, "<title>Test Chapter</title>")
	assert.Contains(t, wrapped, "<h1>Test Chapter</h1>")
	assert.Contains(t, wrapped, "<p>Test content</p>")
	assert.Contains(t, wrapped, "epub:type=\"chapter\"")
}

func TestOPFGenerator(t *testing.T) {
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:       "Test Book",
		Author:      "John Doe",
		Language:    "en",
		Identifier:  "test-123",
		Publisher:   "Test Pub",
		Description: "A test",
		Subject:     []string{"Fiction"},
		Date:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>Content</p>",
	})
	
	generator := NewOPFGenerator(epub)
	opf := generator.Generate()
	
	// Verify OPF structure
	assert.Contains(t, opf, "<?xml version")
	assert.Contains(t, opf, "<package xmlns")
	assert.Contains(t, opf, "version=\"3.0\"")
	assert.Contains(t, opf, "<metadata")
	assert.Contains(t, opf, "<manifest")
	assert.Contains(t, opf, "<spine")
	
	// Verify metadata
	assert.Contains(t, opf, "<dc:title>Test Book</dc:title>")
	assert.Contains(t, opf, "<dc:creator")
	assert.Contains(t, opf, "John Doe")
	assert.Contains(t, opf, "<dc:language>en</dc:language>")
	assert.Contains(t, opf, "<dc:identifier")
	assert.Contains(t, opf, "test-123")
	
	// Verify manifest items
	assert.Contains(t, opf, "id=\"nav\"")
	assert.Contains(t, opf, "id=\"ncx\"")
	assert.Contains(t, opf, "id=\"css\"")
	assert.Contains(t, opf, "id=\"chapter1\"")
	
	// Verify spine
	assert.Contains(t, opf, "<itemref idref=\"chapter1\"")
}

func TestNCXGenerator(t *testing.T) {
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Test Book",
		Identifier: "test-123",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>Content</p>",
	})
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 2",
		Content: "<p>More content</p>",
	})
	
	generator := NewNCXGenerator(epub)
	ncx := generator.Generate()
	
	// Verify NCX structure
	assert.Contains(t, ncx, "<?xml version")
	assert.Contains(t, ncx, "<!DOCTYPE ncx")
	assert.Contains(t, ncx, "<ncx xmlns")
	assert.Contains(t, ncx, "<head>")
	assert.Contains(t, ncx, "<docTitle>")
	assert.Contains(t, ncx, "<navMap>")
	
	// Verify metadata
	assert.Contains(t, ncx, "dtb:uid")
	assert.Contains(t, ncx, "test-123")
	
	// Verify chapters
	assert.Contains(t, ncx, "<navPoint id=\"chapter1\"")
	assert.Contains(t, ncx, "playOrder=\"1\"")
	assert.Contains(t, ncx, "<text>Chapter 1</text>")
	assert.Contains(t, ncx, "<navPoint id=\"chapter2\"")
	assert.Contains(t, ncx, "playOrder=\"2\"")
}

func TestNavGenerator(t *testing.T) {
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title: "Test Book",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>Content</p>",
	})
	
	generator := NewNavGenerator(epub)
	nav := generator.Generate()
	
	// Verify nav structure
	assert.Contains(t, nav, "<?xml version")
	assert.Contains(t, nav, "<!DOCTYPE html>")
	assert.Contains(t, nav, "<html xmlns")
	assert.Contains(t, nav, "xmlns:epub")
	assert.Contains(t, nav, "<nav epub:type=\"toc\"")
	assert.Contains(t, nav, "<h1>Table of Contents</h1>")
	
	// Verify landmarks
	assert.Contains(t, nav, "<nav epub:type=\"landmarks\"")
	
	// Verify chapter link
	assert.Contains(t, nav, "<a href=\"Text/chapter1.xhtml\">Chapter 1</a>")
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"<script>", "&lt;script&gt;"},
		{"Tom & Jerry", "Tom &amp; Jerry"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"'single'", "&apos;single&apos;"},
		{"<>&\"'", "&lt;&gt;&amp;&quot;&apos;"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeXML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultCSS(t *testing.T) {
	css := defaultCSS()
	
	assert.Contains(t, css, "@charset \"UTF-8\"")
	assert.Contains(t, css, "body {")
	assert.Contains(t, css, "font-family: serif")
	assert.Contains(t, css, "text-align: justify")
	assert.Contains(t, css, "h1, h2, h3")
	assert.Contains(t, css, "text-indent")
}

func BenchmarkEPub_Write(b *testing.B) {
	tempDir := b.TempDir()
	
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Benchmark Book",
		Author:     "Benchmark Author",
		Language:   "en",
		Identifier: "bench-123",
	}
	
	for i := 1; i <= 10; i++ {
		epub.AddChapter(Chapter{
			Title:   fmt.Sprintf("Chapter %d", i),
			Content: "<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>",
		})
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tempDir, fmt.Sprintf("bench-%d.epub", i))
		_ = epub.Write(outputPath)
		os.Remove(outputPath) // Cleanup
	}
}
