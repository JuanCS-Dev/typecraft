package latex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplate(t *testing.T) {
	tmpl := NewTemplate("Hello {{name}}")
	require.NotNil(t, tmpl)
	assert.Equal(t, "Hello {{name}}", tmpl.content)
}

func TestTemplate_SetVariable(t *testing.T) {
	tmpl := NewTemplate("")
	tmpl.SetVariable("name", "World")
	
	val, ok := tmpl.GetVariable("name")
	assert.True(t, ok)
	assert.Equal(t, "World", val)
}

func TestTemplate_SetVariables(t *testing.T) {
	tmpl := NewTemplate("")
	tmpl.SetVariables(map[string]interface{}{
		"name":  "John",
		"age":   30,
		"city":  "NYC",
	})
	
	name, _ := tmpl.GetVariable("name")
	age, _ := tmpl.GetVariable("age")
	city, _ := tmpl.GetVariable("city")
	
	assert.Equal(t, "John", name)
	assert.Equal(t, 30, age)
	assert.Equal(t, "NYC", city)
}

func TestTemplate_VariableSubstitution(t *testing.T) {
	tmpl := NewTemplate("Hello {{name}}, you are {{age}} years old!")
	tmpl.SetVariable("name", "Alice")
	tmpl.SetVariable("age", 25)
	
	result, err := tmpl.Render()
	require.NoError(t, err)
	assert.Equal(t, "Hello Alice, you are 25 years old!", result)
}

func TestTemplate_MissingVariable(t *testing.T) {
	tmpl := NewTemplate("Hello {{name}}!")
	
	result, err := tmpl.Render()
	require.NoError(t, err)
	assert.Equal(t, "Hello {{name}}!", result) // Keeps original
}

func TestTemplate_Conditionals(t *testing.T) {
	t.Run("simple if", func(t *testing.T) {
		tmpl := NewTemplate("Start {{if show}}Content{{end}} End")
		tmpl.SetVariable("show", true)
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Start Content End", result)
	})
	
	t.Run("if false", func(t *testing.T) {
		tmpl := NewTemplate("Start {{if show}}Content{{end}} End")
		tmpl.SetVariable("show", false)
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Start  End", result)
	})
	
	t.Run("if else", func(t *testing.T) {
		tmpl := NewTemplate("{{if premium}}Premium{{else}}Basic{{end}}")
		
		tmpl.SetVariable("premium", true)
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Premium", result)
		
		tmpl.SetVariable("premium", false)
		result, err = tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Basic", result)
	})
	
	t.Run("multiline content", func(t *testing.T) {
		tmpl := NewTemplate(`{{if show}}
Line 1
Line 2
{{end}}`)
		tmpl.SetVariable("show", true)
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Contains(t, result, "Line 1")
		assert.Contains(t, result, "Line 2")
	})
}

func TestTemplate_Loops(t *testing.T) {
	t.Run("string slice", func(t *testing.T) {
		tmpl := NewTemplate("Items: {{range items}}{{.}}, {{end}}")
		tmpl.SetVariable("items", []string{"A", "B", "C"})
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Items: A, B, C, ", result)
	})
	
	t.Run("int slice", func(t *testing.T) {
		tmpl := NewTemplate("Numbers: {{range nums}}{{.}} {{end}}")
		tmpl.SetVariable("nums", []int{1, 2, 3, 4})
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Numbers: 1 2 3 4 ", result)
	})
	
	t.Run("interface slice", func(t *testing.T) {
		tmpl := NewTemplate("{{range items}}Item: {{.}}\n{{end}}")
		tmpl.SetVariable("items", []interface{}{"X", "Y", "Z"})
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Contains(t, result, "Item: X")
		assert.Contains(t, result, "Item: Y")
		assert.Contains(t, result, "Item: Z")
	})
	
	t.Run("empty slice", func(t *testing.T) {
		tmpl := NewTemplate("{{range items}}{{.}}{{end}}")
		tmpl.SetVariable("items", []string{})
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
}

func TestTemplate_Functions(t *testing.T) {
	t.Run("upper function", func(t *testing.T) {
		tmpl := NewTemplate("{{upper text}}")
		tmpl.SetVariable("text", "hello")
		
		for name, fn := range StandardFunctions() {
			tmpl.RegisterFunction(name, fn)
		}
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "HELLO", result)
	})
	
	t.Run("lower function", func(t *testing.T) {
		tmpl := NewTemplate("{{lower text}}")
		tmpl.SetVariable("text", "WORLD")
		
		for name, fn := range StandardFunctions() {
			tmpl.RegisterFunction(name, fn)
		}
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "world", result)
	})
	
	t.Run("title function", func(t *testing.T) {
		tmpl := NewTemplate("{{title text}}")
		tmpl.SetVariable("text", "hello world")
		
		for name, fn := range StandardFunctions() {
			tmpl.RegisterFunction(name, fn)
		}
		
		result, err := tmpl.Render()
		require.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})
}

func TestTemplate_ComplexExample(t *testing.T) {
	template := `
\documentclass{article}
\title{ {{title}} }
\author{ {{author}} }

\begin{document}
\maketitle

{{if hasIntro}}
\section{Introduction}
{{intro}}
{{end}}

\section{Chapters}
{{range chapters}}
\subsection{ {{.}} }
{{end}}

{{if hasConclusion}}
\section{Conclusion}
{{conclusion}}
{{end}}

\end{document}
`
	
	tmpl := NewTemplate(template)
	tmpl.SetVariables(map[string]interface{}{
		"title":         "My Book",
		"author":        "John Doe",
		"hasIntro":      true,
		"intro":         "This is the introduction.",
		"chapters":      []string{"Chapter 1", "Chapter 2", "Chapter 3"},
		"hasConclusion": false,
	})
	
	result, err := tmpl.Render()
	require.NoError(t, err)
	
	assert.Contains(t, result, "My Book")
	assert.Contains(t, result, "John Doe")
	assert.Contains(t, result, "Introduction")
	assert.Contains(t, result, "This is the introduction")
	assert.Contains(t, result, "Chapter 1")
	assert.Contains(t, result, "Chapter 2")
	assert.Contains(t, result, "Chapter 3")
	assert.NotContains(t, result, "Conclusion")
}

func TestEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"50% off", "50\\% off"},
		{"$100", "\\$100"},
		{"C++ & Java", "C++ \\& Java"},
		{"file_name", "file\\_name"},
		{"#hashtag", "\\#hashtag"},
		{"{braces}", "\\{braces\\}"},
		{"~tilde", "\\textasciitilde{}tilde"},
		{"x^2", "x\\textasciicircum{}2"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Escape(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateBuilder(t *testing.T) {
	builder := NewTemplateBuilder()
	builder.
		AddSection("\\documentclass{article}").
		AddSection("\\begin{document}").
		AddVariable("title").
		AddConditional("show", "Content").
		AddLoop("items", "Item: {{.}}\n").
		AddSection("\\end{document}")
	
	template := builder.Build()
	
	assert.Contains(t, template, "\\documentclass{article}")
	assert.Contains(t, template, "{{title}}")
	assert.Contains(t, template, "{{if show}}Content{{end}}")
	assert.Contains(t, template, "{{range items}}Item: {{.}}\n{{end}}")
}

func TestStandardFunctions(t *testing.T) {
	funcs := StandardFunctions()
	
	assert.Contains(t, funcs, "upper")
	assert.Contains(t, funcs, "lower")
	assert.Contains(t, funcs, "title")
	assert.Contains(t, funcs, "escape")
	assert.Contains(t, funcs, "repeat")
	
	t.Run("escape function", func(t *testing.T) {
		result, err := funcs["escape"]("Test & Example")
		require.NoError(t, err)
		assert.Equal(t, "Test \\& Example", result)
	})
	
	t.Run("repeat function", func(t *testing.T) {
		result, err := funcs["repeat"]("*", 5)
		require.NoError(t, err)
		assert.Equal(t, "*****", result)
	})
}

func TestTemplate_IsTruthy(t *testing.T) {
	tmpl := NewTemplate("")
	
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{true, true},
		{false, false},
		{nil, false},
		{"", false},
		{"text", true},
		{0, false},
		{1, true},
		{0.0, false},
		{1.5, true},
		{[]string{}, false},
		{[]string{"item"}, true},
		{[]interface{}{}, false},
		{[]interface{}{"item"}, true},
	}
	
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.value), func(t *testing.T) {
			result := tmpl.isTruthy(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkTemplate_Render(b *testing.B) {
	template := `
Title: {{title}}
Author: {{author}}

{{if hasChapters}}
Chapters:
{{range chapters}}
- {{.}}
{{end}}
{{end}}
`
	
	tmpl := NewTemplate(template)
	tmpl.SetVariables(map[string]interface{}{
		"title":       "Benchmark Book",
		"author":      "Bench Author",
		"hasChapters": true,
		"chapters":    []string{"Chapter 1", "Chapter 2", "Chapter 3"},
	})
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = tmpl.Render()
	}
}

func BenchmarkEscape(b *testing.B) {
	text := "This is a test with special characters: & % $ # _ { } ~ ^"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_ = Escape(text)
	}
}
