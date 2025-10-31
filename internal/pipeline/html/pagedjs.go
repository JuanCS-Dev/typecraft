package html

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// PagedJSRenderer renders HTML to PDF using Paged.js CLI
type PagedJSRenderer struct {
	timeout time.Duration
}

// NewPagedJSRenderer creates a new Paged.js renderer with default timeout
func NewPagedJSRenderer() *PagedJSRenderer {
	return &PagedJSRenderer{
		timeout: 60 * time.Second,
	}
}

// RenderOptions holds options for PDF rendering
type RenderOptions struct {
	HTMLPath   string
	OutputPath string
	Timeout    time.Duration
}

// RenderToPDF converts HTML to PDF using pagedjs-cli
func (r *PagedJSRenderer) RenderToPDF(ctx context.Context, opts RenderOptions) error {
	// Validate input file exists
	if _, err := os.Stat(opts.HTMLPath); os.IsNotExist(err) {
		return fmt.Errorf("HTML file not found: %s", opts.HTMLPath)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(opts.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Set timeout from options or use default
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = r.timeout
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Build pagedjs-cli command
	// pagedjs-cli converts HTML to PDF with CSS Paged Media support
	cmd := exec.CommandContext(timeoutCtx, "pagedjs-cli",
		opts.HTMLPath,
		"-o", opts.OutputPath,
	)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("PDF rendering timed out after %v: %w", timeout, err)
		}
		return fmt.Errorf("pagedjs-cli failed: %w\nOutput: %s", err, string(output))
	}

	// Verify output file was created
	if _, err := os.Stat(opts.OutputPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF was not generated at %s", opts.OutputPath)
	}

	return nil
}

// RenderWithConfig renders HTML to PDF with full configuration
type RenderConfig struct {
	HTMLPath      string
	OutputPath    string
	Timeout       time.Duration
	AdditionalCSS []string
}

// RenderWithConfigToPDF renders with additional CSS files
func (r *PagedJSRenderer) RenderWithConfigToPDF(ctx context.Context, cfg RenderConfig) error {
	// For now, delegate to simple render
	// Future: inject additional CSS into HTML before rendering
	opts := RenderOptions{
		HTMLPath:   cfg.HTMLPath,
		OutputPath: cfg.OutputPath,
		Timeout:    cfg.Timeout,
	}

	return r.RenderToPDF(ctx, opts)
}
