package pipeline

import (
	"strings"
	"testing"

	"github.com/JuanCS-Dev/typecraft/internal/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPipelineSelector(t *testing.T) {
	selector := NewPipelineSelector()
	require.NotNil(t, selector)
	require.NotNil(t, selector.academicDetector)
	assert.Equal(t, DefaultThresholds(), selector.thresholds)
}

func TestPipelineSelector_Select(t *testing.T) {
	selector := NewPipelineSelector()
	contentAnalyzer := analyzer.NewContentAnalyzer()

	t.Run("academic paper with many equations - should choose LaTeX", func(t *testing.T) {
		content := `
# Research Paper on Mathematical Models

## Abstract
This comprehensive study examines advanced mathematical models in computational science.

## Introduction
Mathematical modeling has become essential in modern research. The methodology
employs rigorous statistical analysis and experimental validation.

## Methodology
The research implements several key equations. The fundamental relationship is:
$$f(x) = ax^2 + bx + c$$

The trigonometric functions are defined as:
$$g(x) = \sin(x) + \cos(x)$$

Euler's formula appears throughout:
$$h(x) = e^{ix}$$

Logarithmic relationships:
$$i(x) = \log(x)$$

Square root functions:
$$j(x) = \sqrt{x}$$

Power functions with variable exponents:
$$k(x) = x^n$$

Integration formulas:
$$l(x) = \int f(x) dx$$

Differential equations:
$$m(x) = \frac{dy}{dx}$$

Summation notation:
$$n(x) = \sum_{i=1}^{n} x_i$$

Product notation:
$$o(x) = \prod_{i=1}^{n} x_i$$

Limit definitions:
$$p(x) = \lim_{x \to \infty} f(x)$$

The analysis demonstrates significant correlations. The hypothesis is validated
through experimental data. The results indicate strong statistical significance.

## Results
The findings show conclusive evidence. The data analysis reveals patterns.
The methodology proves effective in complex scenarios.

## Conclusion
Therefore, we conclude this approach is mathematically sound.
The research contributes to the body of scientific literature.
`

		analysis, err := contentAnalyzer.Analyze(content)
		require.NoError(t, err)

		decision, err := selector.Select(analysis)
		require.NoError(t, err)
		require.NotNil(t, decision)

		t.Logf("Decision: %s (%.2f confidence)", decision.Pipeline, decision.Confidence)
		t.Logf("Scores: LaTeX=%.2f, HTML=%.2f", decision.Scores.LaTeXScore, decision.Scores.HTMLScore)
		t.Logf("Reasons: %v", decision.Reasons)

		assert.Equal(t, PipelineLaTeX, decision.Pipeline)
		assert.Greater(t, decision.Confidence, 0.5)
		assert.NotEmpty(t, decision.Reasons)
		assert.Greater(t, decision.Scores.LaTeXScore, decision.Scores.HTMLScore)
	})

	t.Run("image-heavy book - should choose HTML", func(t *testing.T) {
		content := `
# Photography Guide

Chapter 1: Basics

![Photo 1](photo1.jpg)
![Photo 2](photo2.jpg)
![Photo 3](photo3.jpg)
![Photo 4](photo4.jpg)
![Photo 5](photo5.jpg)
![Photo 6](photo6.jpg)
![Photo 7](photo7.jpg)
![Photo 8](photo8.jpg)
![Photo 9](photo9.jpg)
![Photo 10](photo10.jpg)
![Photo 11](photo11.jpg)
![Photo 12](photo12.jpg)
![Photo 13](photo13.jpg)
![Photo 14](photo14.jpg)
![Photo 15](photo15.jpg)

Chapter 2: Advanced

More photos here with casual explanations.
This is easy to read and understand.
The images are the main content.
`

		analysis, err := contentAnalyzer.Analyze(content)
		require.NoError(t, err)

		decision, err := selector.Select(analysis)
		require.NoError(t, err)

		t.Logf("Decision: %s (%.2f confidence)", decision.Pipeline, decision.Confidence)
		t.Logf("Scores: LaTeX=%.2f, HTML=%.2f", decision.Scores.LaTeXScore, decision.Scores.HTMLScore)
		t.Logf("Reasons: %v", decision.Reasons)

		assert.Equal(t, PipelineHTML, decision.Pipeline)
		assert.Greater(t, decision.Confidence, 0.5)
		assert.Greater(t, decision.Scores.HTMLScore, decision.Scores.LaTeXScore)
	})

	t.Run("fiction novel - should choose HTML", func(t *testing.T) {
		content := `
# The Detective's Case

Chapter 1

Detective Sarah walked into the room. The case was mysterious.
She examined the clues carefully. The story was beginning.

Chapter 2

The investigation continued. More secrets were revealed.
The plot thickened as she discovered new evidence.
`

		analysis, err := contentAnalyzer.Analyze(content)
		require.NoError(t, err)

		decision, err := selector.Select(analysis)
		require.NoError(t, err)

		t.Logf("Decision: %s (%.2f confidence)", decision.Pipeline, decision.Confidence)
		t.Logf("Reasons: %v", decision.Reasons)

		assert.Equal(t, PipelineHTML, decision.Pipeline)
	})

	t.Run("technical documentation - borderline", func(t *testing.T) {
		content := `
# API Documentation

## Overview

The system implements a RESTful architecture.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api | List all |
| POST | /api | Create new |

The data format is JSON.

Some equations: $x = y + z$

The methodology follows best practices.
`

		analysis, err := contentAnalyzer.Analyze(content)
		require.NoError(t, err)

		decision, err := selector.Select(analysis)
		require.NoError(t, err)

		t.Logf("Decision: %s (%.2f confidence)", decision.Pipeline, decision.Confidence)
		t.Logf("Scores: LaTeX=%.2f, HTML=%.2f", decision.Scores.LaTeXScore, decision.Scores.HTMLScore)

		// Pode ser LaTeX ou HTML dependendo dos scores
		assert.Contains(t, []PipelineType{PipelineLaTeX, PipelineHTML}, decision.Pipeline)
	})

	t.Run("nil analysis - should error", func(t *testing.T) {
		_, err := selector.Select(nil)
		assert.Error(t, err)
	})
}

func TestPipelineSelector_CalculateScores(t *testing.T) {
	selector := NewPipelineSelector()

	t.Run("high LaTeX score", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			EquationCount: 15,
			TableCount:    8,
			Complexity:    0.8,
			Tone: analyzer.ToneProfile{
				Academic: 0.9,
			},
			Formality: 0.85,
		}

		academicScore := &analyzer.AcademicScore{
			IsAcademic:      true,
			Confidence:      0.9,
			EquationDensity: 0.8,
		}

		latexScore := selector.calculateLaTeXScore(analysis, academicScore)
		htmlScore := selector.calculateHTMLScore(analysis, academicScore)

		t.Logf("LaTeX Score: %.2f, HTML Score: %.2f", latexScore, htmlScore)

		assert.Greater(t, latexScore, 0.7)
		assert.Greater(t, latexScore, htmlScore)
	})

	t.Run("high HTML score", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			ImageCount:   25,
			ImageRatio:   0.15,
			Complexity:   0.3,
			EquationCount: 0,
			Tone: analyzer.ToneProfile{
				Creative: 0.8,
				Casual:   0.7,
			},
		}

		academicScore := &analyzer.AcademicScore{
			IsAcademic: false,
			Confidence: 0.1,
		}

		latexScore := selector.calculateLaTeXScore(analysis, academicScore)
		htmlScore := selector.calculateHTMLScore(analysis, academicScore)

		t.Logf("LaTeX Score: %.2f, HTML Score: %.2f", latexScore, htmlScore)

		assert.Greater(t, htmlScore, 0.7)
		assert.Greater(t, htmlScore, latexScore)
	})

	t.Run("balanced scores", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			EquationCount: 3,
			ImageCount:    5,
			TableCount:    2,
			ImageRatio:    0.05,
			Complexity:    0.5,
			Tone: analyzer.ToneProfile{
				Academic: 0.5,
				Creative: 0.5,
			},
		}

		academicScore := &analyzer.AcademicScore{
			IsAcademic:      false,
			Confidence:      0.4,
			EquationDensity: 0.1,
		}

		latexScore := selector.calculateLaTeXScore(analysis, academicScore)
		htmlScore := selector.calculateHTMLScore(analysis, academicScore)

		t.Logf("LaTeX Score: %.2f, HTML Score: %.2f", latexScore, htmlScore)

		// Scores should be relatively close
		diff := latexScore - htmlScore
		assert.LessOrEqual(t, diff, 0.3)
		assert.GreaterOrEqual(t, diff, -0.3)
	})
}

func TestPipelineSelector_CustomThresholds(t *testing.T) {
	customThresholds := PipelineThresholds{
		MathEquationsMin: 5, // Lower threshold
		AcademicConfMin:  0.5,
		ImageRatioMin:    0.2, // Higher threshold
		NeutralZone:      0.1,
	}

	selector := NewPipelineSelectorWithThresholds(customThresholds)
	assert.Equal(t, customThresholds, selector.thresholds)

	// Test that custom thresholds are used
	analysis := &analyzer.ContentAnalysis{
		EquationCount: 6, // Above custom threshold (5) but below default (10)
		Complexity:    0.6,
	}

	decision, err := selector.Select(analysis)
	require.NoError(t, err)

	t.Logf("Decision with custom thresholds: %s", decision.Pipeline)
	t.Logf("Scores: LaTeX=%.2f, HTML=%.2f", decision.Scores.LaTeXScore, decision.Scores.HTMLScore)
}

func TestPipelineSelector_GetReasons(t *testing.T) {
	selector := NewPipelineSelector()

	t.Run("LaTeX reasons", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			EquationCount: 15,
			TableCount:    8,
			Complexity:    0.8,
		}

		academicScore := &analyzer.AcademicScore{
			IsAcademic:      true,
			Confidence:      0.85,
			EquationDensity: 0.7,
		}

		reasons := selector.getLaTeXReasons(analysis, academicScore)
		
		assert.NotEmpty(t, reasons)
		t.Logf("LaTeX Reasons: %v", reasons)
		
		// Should mention high equation count
		hasEquationReason := false
		for _, r := range reasons {
			if strings.Contains(r, "equation") {
				hasEquationReason = true
				break
			}
		}
		assert.True(t, hasEquationReason)
	})

	t.Run("HTML reasons", func(t *testing.T) {
		analysis := &analyzer.ContentAnalysis{
			ImageCount: 25,
			ImageRatio: 0.15,
			Tone: analyzer.ToneProfile{
				Creative: 0.8,
				Casual:   0.7,
			},
		}

		academicScore := &analyzer.AcademicScore{
			IsAcademic: false,
		}

		reasons := selector.getHTMLReasons(analysis, academicScore)
		
		assert.NotEmpty(t, reasons)
		t.Logf("HTML Reasons: %v", reasons)
		
		// Should mention images
		hasImageReason := false
		for _, r := range reasons {
			if strings.Contains(r, "image") || strings.Contains(r, "Image") {
				hasImageReason = true
				break
			}
		}
		assert.True(t, hasImageReason)
	})
}

func TestDefaultThresholds(t *testing.T) {
	thresholds := DefaultThresholds()
	
	assert.Equal(t, 10, thresholds.MathEquationsMin)
	assert.Equal(t, 0.7, thresholds.AcademicConfMin)
	assert.Equal(t, 5, thresholds.ComplexTablesMin)
	assert.Equal(t, 0.1, thresholds.ImageRatioMin)
	assert.Equal(t, 3, thresholds.InteractiveMin)
	assert.Equal(t, 0.2, thresholds.NeutralZone)
}

func BenchmarkPipelineSelector_Select(b *testing.B) {
	selector := NewPipelineSelector()
	contentAnalyzer := analyzer.NewContentAnalyzer()

	content := `
# Technical Paper

The algorithm implements f(x) = x^2. The data shows significant results.
According to research [1], this approach is effective.

![Figure 1](fig1.png)
![Figure 2](fig2.png)

| Method | Result |
|--------|--------|
| A | 95% |
| B | 92% |
`

	analysis, _ := contentAnalyzer.Analyze(content)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = selector.Select(analysis)
	}
}
