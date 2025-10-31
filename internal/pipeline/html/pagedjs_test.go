package html

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPagedJSRenderer_RenderToPDF(t *testing.T) {
	renderer := NewPagedJSRenderer()

	// Create test HTML file
	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "test.html")
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Test Document</title>
    <script src="https://unpkg.com/pagedjs/dist/paged.polyfill.js"></script>
    <style>
        @page {
            size: A4;
            margin: 1in;
        }
        body {
            font-family: serif;
            font-size: 12pt;
            line-height: 1.5;
        }
    </style>
</head>
<body>
    <h1>Test Chapter</h1>
    <p>This is a test paragraph for PDF generation.</p>
</body>
</html>`

	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("Failed to create test HTML: %v", err)
	}

	tests := []struct {
		name    string
		opts    RenderOptions
		wantErr bool
	}{
		{
			name: "successful render",
			opts: RenderOptions{
				HTMLPath:   htmlPath,
				OutputPath: filepath.Join(tmpDir, "output.pdf"),
				Timeout:    30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "non-existent HTML file",
			opts: RenderOptions{
				HTMLPath:   filepath.Join(tmpDir, "nonexistent.html"),
				OutputPath: filepath.Join(tmpDir, "output2.pdf"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := renderer.RenderToPDF(ctx, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderToPDF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify output file exists on success
			if !tt.wantErr {
				if _, err := os.Stat(tt.opts.OutputPath); os.IsNotExist(err) {
					t.Errorf("Output PDF was not created: %s", tt.opts.OutputPath)
				}

				// Check file size is > 0
				info, err := os.Stat(tt.opts.OutputPath)
				if err != nil {
					t.Errorf("Failed to stat output file: %v", err)
				}
				if info.Size() == 0 {
					t.Errorf("Output PDF is empty")
				}
			}
		})
	}
}

func TestPagedJSRenderer_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	renderer := NewPagedJSRenderer()

	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "test.html")
	htmlContent := `<!DOCTYPE html>
<html><body><h1>Test</h1></body></html>`

	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		t.Fatalf("Failed to create test HTML: %v", err)
	}

	opts := RenderOptions{
		HTMLPath:   htmlPath,
		OutputPath: filepath.Join(tmpDir, "output.pdf"),
		Timeout:    1 * time.Nanosecond, // Extremely short timeout
	}

	ctx := context.Background()
	err := renderer.RenderToPDF(ctx, opts)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestNewPagedJSRenderer(t *testing.T) {
	renderer := NewPagedJSRenderer()

	if renderer == nil {
		t.Fatal("NewPagedJSRenderer() returned nil")
	}

	if renderer.timeout != 60*time.Second {
		t.Errorf("Expected default timeout 60s, got %v", renderer.timeout)
	}
}
