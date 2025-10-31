package paged

import (
	"bytes"
	"fmt"
	"html/template"
)

// BookTemplate template completo para um livro
// Conformidade: CONSTITUIÇÃO VÉRTICE v3.0
type BookTemplate struct {
	Title       string
	Author      string
	Publisher   string
	ISBN        string
	Copyright   string
	Content     template.HTML
	CustomCSS   string
	FontFamily  string
	FontSize    string
	LineHeight  string
	TextAlign   string
	Hyphenation bool
}

// DefaultBookTemplate retorna template padrão
func DefaultBookTemplate() BookTemplate {
	return BookTemplate{
		FontFamily:  "Georgia, serif",
		FontSize:    "12pt",
		LineHeight:  "1.6",
		TextAlign:   "justify",
		Hyphenation: true,
	}
}

// Render renderiza o template HTML completo
func (bt BookTemplate) Render() (string, error) {
	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <meta name="author" content="{{.Author}}">
    {{if .Publisher}}<meta name="publisher" content="{{.Publisher}}">{{end}}
    {{if .ISBN}}<meta name="isbn" content="{{.ISBN}}">{{end}}
    
    <style>
        /* Reset e base */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        /* Tipografia base */
        html {
            font-size: {{.FontSize}};
            line-height: {{.LineHeight}};
        }
        
        body {
            font-family: {{.FontFamily}};
            text-align: {{.TextAlign}};
            {{if .Hyphenation}}hyphens: auto;
            -webkit-hyphens: auto;
            -moz-hyphens: auto;{{end}}
        }
        
        /* Elementos tipográficos */
        h1, h2, h3, h4, h5, h6 {
            font-weight: bold;
            margin-top: 1.5em;
            margin-bottom: 0.5em;
            line-height: 1.2;
            page-break-after: avoid;
        }
        
        h1 {
            font-size: 2em;
            page-break-before: always;
        }
        
        h2 {
            font-size: 1.5em;
        }
        
        h3 {
            font-size: 1.25em;
        }
        
        p {
            margin-bottom: 1em;
            orphans: 2;
            widows: 2;
        }
        
        /* Parágrafos especiais */
        p.first-paragraph::first-letter {
            font-size: 3em;
            float: left;
            line-height: 0.9;
            margin: 0.1em 0.1em 0 0;
        }
        
        /* Listas */
        ul, ol {
            margin-left: 2em;
            margin-bottom: 1em;
        }
        
        li {
            margin-bottom: 0.5em;
        }
        
        /* Citações */
        blockquote {
            margin: 1.5em 2em;
            font-style: italic;
            border-left: 3px solid #ccc;
            padding-left: 1em;
        }
        
        /* Código */
        code {
            font-family: 'Courier New', monospace;
            background: #f4f4f4;
            padding: 0.2em 0.4em;
            border-radius: 3px;
        }
        
        pre {
            background: #f4f4f4;
            padding: 1em;
            overflow-x: auto;
            margin-bottom: 1em;
        }
        
        /* Imagens */
        img {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 1em auto;
        }
        
        /* Tabelas */
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 1em;
        }
        
        th, td {
            border: 1px solid #ddd;
            padding: 0.5em;
            text-align: left;
        }
        
        th {
            background: #f4f4f4;
            font-weight: bold;
        }
        
        /* Links */
        a {
            color: #0066cc;
            text-decoration: none;
        }
        
        a:hover {
            text-decoration: underline;
        }
        
        /* Quebras de página */
        .page-break {
            page-break-after: always;
        }
        
        .page-break-before {
            page-break-before: always;
        }
        
        .no-break {
            page-break-inside: avoid;
        }
        
        /* Frontmatter */
        .title-page {
            text-align: center;
            page-break-after: always;
            display: flex;
            flex-direction: column;
            justify-content: center;
            min-height: 80vh;
        }
        
        .title-page h1 {
            font-size: 3em;
            margin-bottom: 0.5em;
        }
        
        .title-page .author {
            font-size: 1.5em;
            font-style: italic;
            margin-top: 2em;
        }
        
        .copyright-page {
            page-break-after: always;
            font-size: 0.9em;
            margin-top: 50vh;
        }
        
        /* TOC */
        .toc {
            page-break-after: always;
        }
        
        .toc h2 {
            text-align: center;
            margin-bottom: 2em;
        }
        
        .toc ul {
            list-style: none;
            margin-left: 0;
        }
        
        .toc li {
            margin-bottom: 0.5em;
        }
        
        .toc a {
            text-decoration: none;
            color: inherit;
        }
        
        /* Custom CSS */
        {{.CustomCSS}}
    </style>
</head>
<body>
    <!-- Title Page -->
    {{if .Title}}
    <div class="title-page">
        <h1>{{.Title}}</h1>
        {{if .Author}}<div class="author">{{.Author}}</div>{{end}}
        {{if .Publisher}}<div class="publisher">{{.Publisher}}</div>{{end}}
    </div>
    {{end}}
    
    <!-- Copyright Page -->
    {{if .Copyright}}
    <div class="copyright-page">
        <p>{{.Copyright}}</p>
        {{if .ISBN}}<p>ISBN: {{.ISBN}}</p>{{end}}
    </div>
    {{end}}
    
    <!-- Main Content -->
    <div class="content">
        {{.Content}}
    </div>
</body>
</html>`

	t, err := template.New("book").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("falha parsear template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, bt); err != nil {
		return "", fmt.Errorf("falha executar template: %w", err)
	}

	return buf.String(), nil
}

// MinimalTemplate template mínimo para testes rápidos
func MinimalTemplate(content string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Georgia, serif;
            font-size: 12pt;
            line-height: 1.6;
            margin: 2cm;
        }
        h1, h2, h3 { page-break-after: avoid; }
        p { margin-bottom: 1em; }
    </style>
</head>
<body>
%s
</body>
</html>`, content)
}
