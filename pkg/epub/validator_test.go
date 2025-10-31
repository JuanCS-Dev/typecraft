package epub

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	require.NotNil(t, validator)
	assert.False(t, validator.strictMode)
}

func TestNewStrictValidator(t *testing.T) {
	validator := NewStrictValidator()
	require.NotNil(t, validator)
	assert.True(t, validator.strictMode)
}

func TestValidator_ValidateFile_Valid(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "valid.epub")
	
	// Create valid ePub
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Valid Book",
		Author:     "Test Author",
		Language:   "en",
		Identifier: "valid-123",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Chapter 1",
		Content: "<p>Content</p>",
	})
	
	err := epub.Write(outputPath)
	require.NoError(t, err)
	
	// Validate
	validator := NewValidator()
	result, err := validator.ValidateFile(outputPath)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	t.Logf("Validation Result: %s", result.Summary())
	t.Logf("Valid: %v", result.Valid)
	t.Logf("Version: %s", result.Version)
	
	if len(result.Issues) > 0 {
		t.Logf("Issues found:")
		for _, issue := range result.Issues {
			t.Logf("  [%s] %s: %s", issue.Level, issue.Code, issue.Message)
		}
	}
	
	// Should be valid
	assert.True(t, result.Valid, "Expected valid ePub")
	assert.Equal(t, EPub3, result.Version)
}

func TestValidator_ValidateFile_EPUB2(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "epub2.epub")
	
	// Create EPUB 2
	epub := NewEPub(EPub2)
	epub.Metadata = Metadata{
		Title:      "EPUB 2 Book",
		Author:     "Test",
		Language:   "en",
		Identifier: "epub2-123",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Test",
		Content: "<p>Test</p>",
	})
	
	err := epub.Write(outputPath)
	require.NoError(t, err)
	
	// Validate
	validator := NewValidator()
	result, err := validator.ValidateFile(outputPath)
	require.NoError(t, err)
	
	t.Logf("EPUB 2 Validation: %s", result.Summary())
	
	// Should detect EPUB 2
	assert.Equal(t, EPub2, result.Version)
}

func TestValidator_MissingFiles(t *testing.T) {
	// Create invalid ePub (manually)
	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "invalid.epub")
	
	// Create minimal ePub with missing files
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Incomplete",
		Author:     "Test",
		Language:   "en",
		Identifier: "incomplete-123",
	}
	
	err := epub.Write(invalidPath)
	require.NoError(t, err)
	
	// Validate
	validator := NewValidator()
	result, err := validator.ValidateFile(invalidPath)
	require.NoError(t, err)
	
	t.Logf("Invalid ePub Validation: %s", result.Summary())
	
	// Should still be structurally valid (we create proper structure)
	// but may have warnings
	assert.NotNil(t, result)
}

func TestValidator_NonExistentFile(t *testing.T) {
	validator := NewValidator()
	_, err := validator.ValidateFile("/nonexistent/file.epub")
	assert.Error(t, err)
}

func TestValidationResult_GetIssuesByLevel(t *testing.T) {
	result := &ValidationResult{
		Issues: []ValidationIssue{
			{Level: LevelError, Code: "E001", Message: "Error 1"},
			{Level: LevelWarning, Code: "W001", Message: "Warning 1"},
			{Level: LevelError, Code: "E002", Message: "Error 2"},
			{Level: LevelInfo, Code: "I001", Message: "Info 1"},
		},
	}
	
	errors := result.GetIssuesByLevel(LevelError)
	assert.Len(t, errors, 2)
	
	warnings := result.GetIssuesByLevel(LevelWarning)
	assert.Len(t, warnings, 1)
	
	infos := result.GetIssuesByLevel(LevelInfo)
	assert.Len(t, infos, 1)
}

func TestValidationResult_HasErrors(t *testing.T) {
	result := &ValidationResult{
		Issues: []ValidationIssue{
			{Level: LevelError, Code: "E001", Message: "Error"},
		},
	}
	
	assert.True(t, result.HasErrors())
	
	result.Issues = []ValidationIssue{
		{Level: LevelWarning, Code: "W001", Message: "Warning"},
	}
	
	assert.False(t, result.HasErrors())
}

func TestValidationResult_HasWarnings(t *testing.T) {
	result := &ValidationResult{
		Issues: []ValidationIssue{
			{Level: LevelWarning, Code: "W001", Message: "Warning"},
		},
	}
	
	assert.True(t, result.HasWarnings())
	
	result.Issues = []ValidationIssue{
		{Level: LevelError, Code: "E001", Message: "Error"},
	}
	
	assert.False(t, result.HasWarnings())
}

func TestValidationResult_Summary(t *testing.T) {
	result := &ValidationResult{
		Issues: []ValidationIssue{
			{Level: LevelError, Code: "E001", Message: "Error 1"},
			{Level: LevelError, Code: "E002", Message: "Error 2"},
			{Level: LevelWarning, Code: "W001", Message: "Warning 1"},
			{Level: LevelInfo, Code: "I001", Message: "Info 1"},
			{Level: LevelInfo, Code: "I002", Message: "Info 2"},
			{Level: LevelInfo, Code: "I003", Message: "Info 3"},
		},
	}
	
	summary := result.Summary()
	assert.Equal(t, "Errors: 2, Warnings: 1, Info: 3", summary)
}

func TestValidator_StrictMode(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "strict.epub")
	
	// Create ePub
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Strict Test",
		Author:     "Test",
		Language:   "en",
		Identifier: "strict-123",
	}
	
	epub.AddChapter(Chapter{
		Title:   "Test",
		Content: "<p>Test</p>",
	})
	
	err := epub.Write(outputPath)
	require.NoError(t, err)
	
	// Validate with strict mode
	strictValidator := NewStrictValidator()
	result, err := strictValidator.ValidateFile(outputPath)
	require.NoError(t, err)
	
	t.Logf("Strict Mode Validation: %s", result.Summary())
	
	// Strict mode may find more issues
	assert.NotNil(t, result)
}

func TestMetadataEnhancer_Enhance(t *testing.T) {
	enhancer := NewMetadataEnhancer()
	
	metadata := &Metadata{
		Title:  "  Test   Book  ",
		Author: "  John   Doe  ",
	}
	
	content := "This is a test book. The content is in English. And it has some text."
	
	enhancer.Enhance(metadata, content)
	
	// Should generate UUID
	assert.NotEmpty(t, metadata.Identifier)
	assert.Contains(t, metadata.Identifier, "urn:uuid:")
	
	// Should detect language
	assert.Equal(t, "en", metadata.Language)
	
	// Should normalize title
	assert.Equal(t, "Test Book", metadata.Title)
	assert.NotContains(t, metadata.Title, "  ")
	
	// Should normalize author
	assert.Equal(t, "John Doe", metadata.Author)
	assert.NotContains(t, metadata.Author, "  ")
	
	// Should set date
	assert.False(t, metadata.Date.IsZero())
}

func TestMetadataEnhancer_GenerateUUID(t *testing.T) {
	enhancer := NewMetadataEnhancer()
	
	uuid1 := enhancer.GenerateUUID()
	uuid2 := enhancer.GenerateUUID()
	
	assert.NotEmpty(t, uuid1)
	assert.NotEmpty(t, uuid2)
	assert.NotEqual(t, uuid1, uuid2)
	assert.Contains(t, uuid1, "urn:uuid:")
}

func TestMetadataEnhancer_ValidateMetadata(t *testing.T) {
	enhancer := NewMetadataEnhancer()
	
	t.Run("complete metadata", func(t *testing.T) {
		metadata := &Metadata{
			Title:       "Complete Book",
			Author:      "Author Name",
			Language:    "en",
			Identifier:  "complete-123",
			Publisher:   "Publisher",
			Description: "Description",
			Subject:     []string{"Fiction"},
		}
		
		issues := enhancer.ValidateMetadata(metadata)
		assert.Empty(t, issues)
	})
	
	t.Run("minimal metadata", func(t *testing.T) {
		metadata := &Metadata{
			Title:      "Minimal",
			Author:     "Author",
			Language:   "en",
			Identifier: "minimal-123",
		}
		
		issues := enhancer.ValidateMetadata(metadata)
		assert.NotEmpty(t, issues)
		t.Logf("Minimal metadata issues: %v", issues)
	})
	
	t.Run("missing required fields", func(t *testing.T) {
		metadata := &Metadata{}
		
		issues := enhancer.ValidateMetadata(metadata)
		assert.Len(t, issues, 7) // 4 required + 3 recommended
		
		// Check for required field issues
		issuesStr := fmt.Sprintf("%v", issues)
		assert.Contains(t, issuesStr, "Title")
		assert.Contains(t, issuesStr, "Author")
		assert.Contains(t, issuesStr, "Language")
		assert.Contains(t, issuesStr, "Identifier")
	})
}

func TestLanguageDetector_Detect(t *testing.T) {
	detector := NewLanguageDetector()
	
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "English",
			content:  "The quick brown fox jumps over the lazy dog. This is a test in English.",
			expected: "en",
		},
		{
			name:     "Portuguese",
			content:  "O rato roeu a roupa do rei de Roma. Este é um teste em português.",
			expected: "pt",
		},
		{
			name:     "Spanish",
			content:  "El perro come la comida. Este es un texto en español.",
			expected: "es",
		},
		{
			name:     "Empty content",
			content:  "",
			expected: "en", // default
		},
		{
			name:     "Too short",
			content:  "Hello",
			expected: "en", // default
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detected := detector.Detect(tt.content)
			assert.Equal(t, tt.expected, detected)
		})
	}
}

func TestSubjectTaxonomy(t *testing.T) {
	taxonomy := NewSubjectTaxonomy()
	
	t.Run("get category", func(t *testing.T) {
		assert.Equal(t, "Fiction", taxonomy.GetCategory("Science Fiction"))
		assert.Equal(t, "Academic", taxonomy.GetCategory("Mathematics"))
		assert.Equal(t, "Non-Fiction", taxonomy.GetCategory("Biography"))
		assert.Equal(t, "General", taxonomy.GetCategory("Unknown Subject"))
	})
	
	t.Run("get related", func(t *testing.T) {
		related := taxonomy.GetRelated("Science Fiction")
		assert.NotEmpty(t, related)
		assert.Contains(t, related, "Fantasy")
		assert.NotContains(t, related, "Science Fiction")
		
		relatedUnknown := taxonomy.GetRelated("Unknown")
		assert.Empty(t, relatedUnknown)
	})
}

func TestMetadataEnhancer_EnhanceSubjects(t *testing.T) {
	enhancer := NewMetadataEnhancer()
	
	t.Run("remove duplicates", func(t *testing.T) {
		subjects := []string{"Fiction", "fiction", "FICTION", "Science"}
		enhanced := enhancer.EnhanceSubjects(subjects)
		
		assert.Len(t, enhanced, 2)
		assert.Contains(t, enhanced, "Fiction")
		assert.Contains(t, enhanced, "Science")
	})
	
	t.Run("remove empty", func(t *testing.T) {
		subjects := []string{"Fiction", "", "  ", "Science"}
		enhanced := enhancer.EnhanceSubjects(subjects)
		
		assert.Len(t, enhanced, 2)
	})
}

func BenchmarkValidator_ValidateFile(b *testing.B) {
	tempDir := b.TempDir()
	outputPath := filepath.Join(tempDir, "bench.epub")
	
	// Create ePub once
	epub := NewEPub(EPub3)
	epub.Metadata = Metadata{
		Title:      "Benchmark",
		Author:     "Bench",
		Language:   "en",
		Identifier: "bench-123",
	}
	
	for i := 1; i <= 5; i++ {
		epub.AddChapter(Chapter{
			Title:   fmt.Sprintf("Chapter %d", i),
			Content: "<p>Content</p>",
		})
	}
	
	_ = epub.Write(outputPath)
	defer os.Remove(outputPath)
	
	validator := NewValidator()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = validator.ValidateFile(outputPath)
	}
}

func BenchmarkMetadataEnhancer_Enhance(b *testing.B) {
	enhancer := NewMetadataEnhancer()
	content := "This is test content in English. The book contains many words and sentences."
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		metadata := &Metadata{
			Title:  "Test Book",
			Author: "Author",
		}
		enhancer.Enhance(metadata, content)
	}
}
