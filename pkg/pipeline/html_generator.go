package pipeline

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/JuanCS-Dev/typecraft/pkg/ai"
	"github.com/JuanCS-Dev/typecraft/pkg/typography"
)

// HTMLGenerator gera HTML tipograficamente correto para impressão
type HTMLGenerator struct {
	styleEngine *typography.StyleEngine
	aiClient    *ai.Client
}

// NewHTMLGenerator cria um novo gerador de HTML
func NewHTMLGenerator(styleEngine *typography.StyleEngine, aiClient *ai.Client) *HTMLGenerator {
	return &HTMLGenerator{
		styleEngine: styleEngine,
		aiClient:    aiClient,
	}
}

// BookSection representa uma seção do livro
type BookSection struct {
	Title    string
	Content  string
	Type     string // chapter, preface, appendix
	Number   int
	Metadata map[string]interface{}
}

// GenerateHTML cria HTML estruturado do documento processado
func (h *HTMLGenerator) GenerateHTML(sections []BookSection, metadata map[string]interface{}) (string, error) {
	tmpl, err := template.New("book").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"styleClass": func(t string) string {
			return fmt.Sprintf("section-%s", t)
		},
	}).Parse(bookTemplate)
	
	if err != nil {
		return "", fmt.Errorf("erro ao criar template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Sections": sections,
		"Metadata": metadata,
		"Title":    metadata["title"],
		"Author":   metadata["author"],
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("erro ao executar template: %w", err)
	}

	return buf.String(), nil
}

// ProcessChapter aplica tipografia e formatação a um capítulo
func (h *HTMLGenerator) ProcessChapter(rawContent string, chapterNum int) (BookSection, error) {
	// Aplica regras tipográficas
	styled := h.styleEngine.ApplyRules(rawContent)

	// Enriquece com IA se disponível
	if h.aiClient != nil {
		enriched, err := h.aiClient.EnhanceTypography(styled)
		if err == nil {
			styled = enriched
		}
	}

	// Estrutura o conteúdo em parágrafos
	paragraphs := strings.Split(styled, "\n\n")
	var htmlContent strings.Builder

	for _, p := range paragraphs {
		if strings.TrimSpace(p) == "" {
			continue
		}
		htmlContent.WriteString(fmt.Sprintf("<p>%s</p>\n", p))
	}

	return BookSection{
		Title:   fmt.Sprintf("Capítulo %d", chapterNum),
		Content: htmlContent.String(),
		Type:    "chapter",
		Number:  chapterNum,
	}, nil
}

// GeneratePagedJS cria HTML com estilos Paged.js para paginação
func (h *HTMLGenerator) GeneratePagedJS(sections []BookSection, metadata map[string]interface{}) (string, error) {
	html, err := h.GenerateHTML(sections, metadata)
	if err != nil {
		return "", err
	}

	// Adiciona configuração Paged.js
	pagedConfig := `
<script src="https://unpkg.com/pagedjs/dist/paged.polyfill.js"></script>
<script>
  class PagedHandler extends Paged.Handler {
    constructor(chunker, polisher, caller) {
      super(chunker, polisher, caller);
    }
    
    afterRendered(pages) {
      console.log('Rendered', pages.length, 'pages');
    }
  }
  
  Paged.registerHandlers(PagedHandler);
</script>
`

	// Injeta antes do </body>
	html = strings.Replace(html, "</body>", pagedConfig+"</body>", 1)
	
	return html, nil
}

// ApplyDesignSystem aplica sistema de design gerado por IA
func (h *HTMLGenerator) ApplyDesignSystem(html string, designPrompt string) (string, error) {
	if h.aiClient == nil {
		return html, nil
	}

	// Solicita sistema de design à IA
	design, err := h.aiClient.GenerateDesignSystem(designPrompt)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar design system: %w", err)
	}

	// Injeta CSS gerado
	cssBlock := fmt.Sprintf("<style>\n%s\n</style>", design.CSS)
	html = strings.Replace(html, "</head>", cssBlock+"</head>", 1)

	return html, nil
}

const bookTemplate = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        @page {
            size: 148mm 210mm; /* A5 */
            margin: 20mm 15mm;
            
            @top-center {
                content: "{{.Title}}";
                font-size: 9pt;
                font-family: serif;
            }
            
            @bottom-center {
                content: counter(page);
                font-size: 9pt;
            }
        }
        
        @page :first {
            @top-center { content: none; }
            @bottom-center { content: none; }
        }
        
        body {
            font-family: 'Crimson Pro', 'Crimson Text', 'Garamond', serif;
            font-size: 11pt;
            line-height: 1.6;
            text-align: justify;
            hyphens: auto;
            widows: 2;
            orphans: 2;
        }
        
        h1 {
            font-size: 24pt;
            line-height: 1.2;
            margin-top: 0;
            margin-bottom: 2em;
            text-align: center;
            page-break-before: always;
        }
        
        h2 {
            font-size: 18pt;
            line-height: 1.3;
            margin-top: 2em;
            margin-bottom: 1em;
            page-break-after: avoid;
        }
        
        p {
            margin: 0;
            text-indent: 1.5em;
        }
        
        p:first-of-type,
        h1 + p,
        h2 + p {
            text-indent: 0;
        }
        
        .section-chapter {
            page-break-before: always;
        }
        
        .section-preface,
        .section-appendix {
            page-break-before: always;
        }
        
        /* Tipografia avançada */
        .drop-cap::first-letter {
            font-size: 3em;
            line-height: 0.8;
            float: left;
            margin: 0.1em 0.1em 0 0;
        }
        
        /* Estilos para citações */
        blockquote {
            margin: 1em 2em;
            font-style: italic;
        }
        
        /* Notas de rodapé */
        .footnote {
            font-size: 8pt;
            line-height: 1.4;
        }
    </style>
</head>
<body>
    <div class="title-page">
        <h1>{{.Title}}</h1>
        {{if .Author}}
        <p style="text-align: center; font-size: 14pt;">{{.Author}}</p>
        {{end}}
    </div>
    
    {{range .Sections}}
    <section class="{{styleClass .Type}}">
        <h2>{{.Title}}</h2>
        {{.Content | safeHTML}}
    </section>
    {{end}}
</body>
</html>
`
