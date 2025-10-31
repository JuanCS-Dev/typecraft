package html

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FontSubsetter handles font subsetting using Python fonttools
type FontSubsetter struct {
	pythonPath string
	scriptPath string
}

// NewFontSubsetter creates a new font subsetter
func NewFontSubsetter(scriptPath string) *FontSubsetter {
	pythonPath := "python3"
	if customPath := os.Getenv("PYTHON_PATH"); customPath != "" {
		pythonPath = customPath
	}

	return &FontSubsetter{
		pythonPath: pythonPath,
		scriptPath: scriptPath,
	}
}

// SubsetOptions holds options for font subsetting
type SubsetOptions struct {
	InputFont  string
	OutputFont string
	Text       string
	Unicodes   []rune
}

// SubsetFont subsets a font to include only specified characters
func (fs *FontSubsetter) SubsetFont(ctx context.Context, opts SubsetOptions) error {
	// Validate input font exists
	if _, err := os.Stat(opts.InputFont); os.IsNotExist(err) {
		return fmt.Errorf("font file not found: %s", opts.InputFont)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(opts.OutputFont)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build pyftsubset command (from fonttools)
	// pyftsubset reduces font file size by including only used glyphs
	args := []string{
		opts.InputFont,
		fmt.Sprintf("--output-file=%s", opts.OutputFont),
		"--flavor=woff2", // Modern web format
		"--layout-features=*",
		"--desubroutinize",
	}

	if opts.Text != "" {
		args = append(args, fmt.Sprintf("--text=%s", opts.Text))
	}

	cmd := exec.CommandContext(ctx, "pyftsubset", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pyftsubset failed: %w\nOutput: %s", err, string(output))
	}

	// Verify output file was created
	if _, err := os.Stat(opts.OutputFont); os.IsNotExist(err) {
		return fmt.Errorf("subset font was not generated at %s", opts.OutputFont)
	}

	return nil
}

// ExtractUsedChars extracts all unique characters from HTML content
func ExtractUsedChars(htmlContent string) string {
	charMap := make(map[rune]bool)

	for _, char := range htmlContent {
		charMap[char] = true
	}

	var chars []rune
	for char := range charMap {
		chars = append(chars, char)
	}

	return string(chars)
}

// SubsetFromHTML subsets font based on HTML content
func (fs *FontSubsetter) SubsetFromHTML(ctx context.Context, fontPath, htmlContent, outputPath string) error {
	text := ExtractUsedChars(htmlContent)

	return fs.SubsetFont(ctx, SubsetOptions{
		InputFont:  fontPath,
		OutputFont: outputPath,
		Text:       text,
	})
}
