package html

import (
	"bytes"
	"fmt"
	"text/template"
)

// HTMLTemplate define um template HTML5 completo
type HTMLTemplate struct {
	Title       string
	Author      string
	Language    string
	CSS         string
	Content     string
	Metadata    map[string]string
}

// TemplateGenerator gerencia templates HTML
type TemplateGenerator struct {
	baseTemplate *template.Template
}

// NewTemplateGenerator cria um gerador de templates
func NewTemplateGenerator() *TemplateGenerator {
	tmpl := template.Must(template.New("html").Parse(baseHTML5Template))
	return &TemplateGenerator{
		baseTemplate: tmpl,
	}
}

// Generate gera HTML completo a partir do template
func (tg *TemplateGenerator) Generate(data HTMLTemplate) (string, error) {
	var buf bytes.Buffer
	err := tg.baseTemplate.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar HTML: %w", err)
	}
	return buf.String(), nil
}

// baseHTML5Template é o template HTML5 base compatível com Paged.js
const baseHTML5Template = `<!DOCTYPE html>
<html lang="{{.Language}}">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="generator" content="Typecraft Editorial Automation System">
  
  {{if .Author}}<meta name="author" content="{{.Author}}">{{end}}
  
  <title>{{.Title}}</title>
  
  <!-- Paged.js for CSS Paged Media -->
  <script src="https://unpkg.com/pagedjs/dist/paged.polyfill.js"></script>
  
  <!-- Custom CSS -->
  <style>
{{.CSS}}
  </style>
  
  <!-- Additional metadata -->
  {{range $key, $value := .Metadata}}
  <meta name="{{$key}}" content="{{$value}}">
  {{end}}
</head>
<body>
  <main class="content">
{{.Content}}
  </main>
</body>
</html>`
