package ai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxChars int
		want     string
	}{
		{
			name:     "text shorter than max",
			text:     "Short text",
			maxChars: 100,
			want:     "Short text",
		},
		{
			name:     "text exactly at max",
			text:     "Exact",
			maxChars: 5,
			want:     "Exact",
		},
		{
			name:     "truncate at period",
			text:     "First sentence. Second sentence. Third.",
			maxChars: 25,
			want:     "First sentence.",
		},
		{
			name:     "truncate with ellipsis when no good break",
			text:     "NoBreaksHereAtAll",
			maxChars: 10,
			want:     "NoBreaksHe...",
		},
		{
			name:     "long text truncates correctly",
			text:     strings.Repeat("word ", 1000),
			maxChars: 100,
			want:     func() string {
				long := strings.Repeat("word ", 1000)
				return long[:100] + "..."
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateText(tt.text, tt.maxChars)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildAnalysisPrompt(t *testing.T) {
	text := "This is a sample text for analysis."
	
	prompt := buildAnalysisPrompt(text)
	
	assert.Contains(t, prompt, "MANUSCRITO:")
	assert.Contains(t, prompt, text)
	assert.Contains(t, prompt, "JSON")
	assert.Contains(t, prompt, "gênero")
	assert.Contains(t, prompt, "pipeline")
}

func TestBuildAnalysisPromptTruncation(t *testing.T) {
	longText := ""
	for i := 0; i < 10000; i++ {
		longText += "word "
	}
	
	prompt := buildAnalysisPrompt(longText)
	
	// Prompt não deve ser gigante
	assert.Less(t, len(prompt), 7000, "Prompt should be truncated")
	assert.Contains(t, prompt, "MANUSCRITO:")
}
