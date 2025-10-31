package html

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractUsedChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int
		contains []rune
	}{
		{
			name:     "simple text",
			input:    "Hello World",
			wantLen:  8, // H, e, l, o, space, W, r, d
			contains: []rune{'H', 'e', 'l', 'o', ' ', 'W', 'r', 'd'},
		},
		{
			name:     "with special chars",
			input:    "Hello, World!",
			wantLen:  10,
			contains: []rune{'H', 'e', 'l', 'o', ',', ' ', 'W', 'r', 'd', '!'},
		},
		{
			name:     "unicode chars",
			input:    "Olá Mundo",
			wantLen:  10, // O, l, á (2 chars in Unicode), space, M, u, n, d, o
			contains: []rune{'O', 'l', 'á', ' ', 'M', 'u', 'n', 'd', 'o'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractUsedChars(tt.input)

			// Check length
			if len(result) != tt.wantLen {
				t.Errorf("ExtractUsedChars() length = %d, want %d", len(result), tt.wantLen)
			}

			// Check all expected chars are present
			for _, char := range tt.contains {
				found := false
				for _, r := range result {
					if r == char {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected char %q not found in result", char)
				}
			}
		})
	}
}

func TestFontSubsetter_SubsetFont(t *testing.T) {
	// This test requires pyftsubset to be installed
	// Skip if not available
	if testing.Short() {
		t.Skip("Skipping font subsetting test in short mode")
	}

	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "subset.py")

	subsetter := NewFontSubsetter(scriptPath)

	// We can't test actual font subsetting without a real font file
	// But we can test the interface and error handling
	t.Run("non-existent font", func(t *testing.T) {
		opts := SubsetOptions{
			InputFont:  "/nonexistent/font.ttf",
			OutputFont: filepath.Join(tmpDir, "output.woff2"),
			Text:       "Hello",
		}

		ctx := context.Background()
		err := subsetter.SubsetFont(ctx, opts)

		if err == nil {
			t.Error("Expected error for non-existent font, got nil")
		}
	})
}

func TestFontSubsetter_SubsetFromHTML(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping font subsetting test in short mode")
	}

	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "subset.py")

	subsetter := NewFontSubsetter(scriptPath)

	htmlContent := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body><h1>Hello World</h1><p>Test content.</p></body>
</html>`

	t.Run("extract chars from HTML", func(t *testing.T) {
		chars := ExtractUsedChars(htmlContent)

		// Should contain basic letters, punctuation, and HTML chars
		if len(chars) == 0 {
			t.Error("Expected non-empty character extraction")
		}

		// Should contain basic ASCII letters
		expectedChars := []rune{'H', 'e', 'l', 'o', 'W', 'r', 'd', 'T', 's', 't'}
		for _, expected := range expectedChars {
			found := false
			for _, char := range chars {
				if char == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected char %q not found in extracted chars", expected)
			}
		}
	})

	t.Run("subset from HTML", func(t *testing.T) {
		// Test with non-existent font to verify error handling
		ctx := context.Background()
		err := subsetter.SubsetFromHTML(ctx, "/nonexistent/font.ttf", htmlContent, filepath.Join(tmpDir, "output.woff2"))

		if err == nil {
			t.Error("Expected error for non-existent font, got nil")
		}
	})
}

func TestNewFontSubsetter(t *testing.T) {
	scriptPath := "/tmp/test_script.py"
	subsetter := NewFontSubsetter(scriptPath)

	if subsetter == nil {
		t.Fatal("NewFontSubsetter() returned nil")
	}

	if subsetter.scriptPath != scriptPath {
		t.Errorf("Expected scriptPath %s, got %s", scriptPath, subsetter.scriptPath)
	}

	if subsetter.pythonPath == "" {
		t.Error("pythonPath should not be empty")
	}
}

func TestFontSubsetter_CustomPythonPath(t *testing.T) {
	// Set custom Python path
	customPath := "/usr/local/bin/python3"
	os.Setenv("PYTHON_PATH", customPath)
	defer os.Unsetenv("PYTHON_PATH")

	subsetter := NewFontSubsetter("/tmp/script.py")

	if subsetter.pythonPath != customPath {
		t.Errorf("Expected pythonPath %s, got %s", customPath, subsetter.pythonPath)
	}
}
