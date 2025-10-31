# Exemplo: Conversão Markdown → HTML → PDF com Paged.js

Este exemplo demonstra como usar o pipeline HTML/CSS do Typecraft para converter Markdown em PDF profissional usando Paged.js.

## Pré-requisitos

```bash
# Instalar pagedjs-cli globalmente
npm install -g pagedjs-cli

# Instalar Python fonttools para font subsetting (opcional)
pip install fonttools brotli
```

## Uso Básico

```go
package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/pipeline/html"
)

func main() {
	// 1. Criar renderizador Paged.js
	renderer := html.NewPagedJSRenderer()

	// 2. Configurar opções de renderização
	opts := html.RenderOptions{
		HTMLPath:   "input.html",
		OutputPath: "output.pdf",
		Timeout:    60 * time.Second,
	}

	// 3. Renderizar HTML → PDF
	ctx := context.Background()
	if err := renderer.RenderToPDF(ctx, opts); err != nil {
		log.Fatalf("Erro ao renderizar PDF: %v", err)
	}

	fmt.Println("PDF gerado com sucesso:", opts.OutputPath)
}
```

## Exemplo com Font Subsetting

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JuanCS-Dev/typecraft/internal/pipeline/html"
)

func main() {
	// 1. Criar font subsetter
	subsetter := html.NewFontSubsetter("")

	// 2. Ler conteúdo HTML
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Meu Livro</title>
</head>
<body>
    <h1>Capítulo 1</h1>
    <p>Este é o conteúdo do meu livro.</p>
</body>
</html>`

	// 3. Fazer subset da fonte baseado no texto usado
	ctx := context.Background()
	err := subsetter.SubsetFromHTML(
		ctx,
		"/path/to/SourceSerifPro-Regular.ttf",
		htmlContent,
		"fonts/subset.woff2",
	)

	if err != nil {
		log.Fatalf("Erro ao fazer subset: %v", err)
	}

	fmt.Println("Font subset criado com sucesso!")
}
```

## Exemplo Completo: Pipeline End-to-End

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/pipeline/html"
)

type BookData struct {
	Title      string
	Author     string
	Content    string
	PageSize   string
	BodyFont   string
	HeadingFont string
}

func main() {
	// 1. Preparar dados do livro
	book := BookData{
		Title:       "Minha Obra Prima",
		Author:      "João Silva",
		Content:     "<h1>Capítulo 1</h1><p>Era uma vez...</p>",
		PageSize:    "6in 9in",
		BodyFont:    "Source Serif Pro",
		HeadingFont: "Source Sans Pro",
	}

	// 2. Carregar template HTML
	tmpl, err := template.ParseFiles("templates/pagedjs/base.html")
	if err != nil {
		log.Fatalf("Erro ao carregar template: %v", err)
	}

	// 3. Renderizar template com dados
	htmlFile, err := os.Create("temp.html")
	if err != nil {
		log.Fatalf("Erro ao criar arquivo HTML: %v", err)
	}
	defer htmlFile.Close()
	defer os.Remove("temp.html")

	if err := tmpl.Execute(htmlFile, book); err != nil {
		log.Fatalf("Erro ao executar template: %v", err)
	}
	htmlFile.Close()

	// 4. Renderizar HTML → PDF
	renderer := html.NewPagedJSRenderer()
	opts := html.RenderOptions{
		HTMLPath:   "temp.html",
		OutputPath: "meu_livro.pdf",
		Timeout:    60 * time.Second,
	}

	ctx := context.Background()
	if err := renderer.RenderToPDF(ctx, opts); err != nil {
		log.Fatalf("Erro ao gerar PDF: %v", err)
	}

	fmt.Println("✅ PDF gerado com sucesso: meu_livro.pdf")
}
```

## HTML de Entrada Exemplo

Arquivo: `input.html`

```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Meu Livro</title>
    <script src="https://unpkg.com/pagedjs/dist/paged.polyfill.js"></script>
    <style>
        @page {
            size: 6in 9in;
            margin: 0.75in 1in 1.125in 0.5in;
            
            @top-center {
                content: "Meu Livro";
                font-family: serif;
                font-size: 9pt;
            }
            
            @bottom-center {
                content: counter(page);
            }
        }
        
        body {
            font-family: 'Source Serif Pro', Georgia, serif;
            font-size: 11pt;
            line-height: 1.5;
            text-align: justify;
        }
        
        h1 {
            font-family: 'Source Sans Pro', Arial, sans-serif;
            font-size: 24pt;
            page-break-before: always;
            text-align: center;
        }
        
        p {
            text-indent: 1.5em;
            margin: 0 0 1em 0;
        }
    </style>
</head>
<body>
    <h1>Capítulo 1</h1>
    <p>Era uma vez, em uma terra muito distante...</p>
    
    <h1>Capítulo 2</h1>
    <p>Continuando a história...</p>
</body>
</html>
```

## Executar o Exemplo

```bash
# Compilar
go run examples/pagedjs/main.go

# Ou usando make
make run-example-pagedjs
```

## Características Implementadas

### ✅ CSS Paged Media Support
- `@page` rules para definição de página
- Running headers e footers
- Page numbering automático
- Page breaks controlados

### ✅ Van de Graaf Canon
- Margens proporcionais (2:3:4:6)
- Harmonia matemática na página

### ✅ Typography
- Fontes serifadas para corpo de texto
- Fontes sans-serif para títulos
- Entrelinha otimizada (1.5)
- Justificação de texto
- Controle de órfãs e viúvas

### ✅ Performance
- Timeout configurável
- Renderização em ~5-30 segundos (depende do tamanho)

## Troubleshooting

### Erro: "pagedjs-cli: command not found"
```bash
npm install -g pagedjs-cli
```

### Erro: "pyftsubset: command not found"
```bash
pip install fonttools brotli
```

### PDF não renderiza corretamente
- Verifique se o HTML é válido
- Confirme que o Paged.js script está sendo carregado
- Teste com HTML minimalista primeiro

## Próximos Passos

1. **Design IA**: Integrar geração de paleta de cores
2. **Font Pairing**: Sugestão automática de fontes
3. **ePub Generation**: Reusar HTML/CSS para ePub 3
4. **Optimization**: Paralelização de renderização

## Referências

- [Paged.js Documentation](https://pagedjs.org/documentation/)
- [CSS Paged Media Module Level 3](https://www.w3.org/TR/css-page-3/)
- [Blueprint Automação Editorial](../../DOCS/1_BLUEPRINT_Sistema_Automacao_Editorial.md)
