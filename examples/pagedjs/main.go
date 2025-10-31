package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/pipeline/html"
)

func main() {
	fmt.Println("🎯 Typecraft - Exemplo Paged.js")
	fmt.Println("================================")

	// Criar diretório de output
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("❌ Erro ao criar diretório: %v", err)
	}

	// Criar HTML de teste
	htmlContent := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Livro de Exemplo - Typecraft</title>
    <script src="https://unpkg.com/pagedjs/dist/paged.polyfill.js"></script>
    <style>
        @page {
            size: 6in 9in;
            margin: 0.75in 1in 1.125in 0.5in;
            
            @top-center {
                content: "Livro de Exemplo";
                font-family: serif;
                font-size: 9pt;
                color: #666;
            }
            
            @bottom-center {
                content: counter(page);
                font-family: serif;
                font-size: 10pt;
            }
        }
        
        @page :first {
            @top-center { content: none; }
            @bottom-center { content: none; }
        }
        
        @page chapter {
            @top-center { content: none; }
        }
        
        body {
            font-family: Georgia, serif;
            font-size: 11pt;
            line-height: 1.5;
            text-align: justify;
            color: #1a1a1a;
        }
        
        h1 {
            font-family: Arial, sans-serif;
            font-size: 24pt;
            page: chapter;
            page-break-before: always;
            text-align: center;
            margin-top: 2in;
            margin-bottom: 1in;
            font-weight: 700;
        }
        
        h2 {
            font-family: Arial, sans-serif;
            font-size: 18pt;
            margin-top: 1.5em;
            margin-bottom: 0.75em;
            font-weight: 600;
        }
        
        p {
            text-indent: 1.5em;
            margin: 0 0 1em 0;
            orphans: 3;
            widows: 3;
        }
        
        h1 + p, h2 + p {
            text-indent: 0;
        }
        
        .title-page {
            text-align: center;
            margin-top: 3in;
        }
        
        .title-page h1 {
            page-break-before: avoid;
            margin: 0;
            font-size: 32pt;
        }
        
        .title-page .author {
            font-size: 14pt;
            margin-top: 1em;
            font-style: italic;
        }
        
        blockquote {
            margin: 1em 2em;
            padding-left: 1em;
            border-left: 3px solid #2563eb;
            font-style: italic;
            color: #555;
        }
    </style>
</head>
<body>
    <div class="title-page">
        <h1>A Jornada do Typecraft</h1>
        <p class="author">Sistema de Automação Editorial</p>
    </div>
    
    <h1>Capítulo 1: A Visão</h1>
    
    <p>Em nome de Jesus, iniciamos esta jornada de transformar a publicação de livros através da automação inteligente. O Typecraft não é apenas um sistema; é uma missão para democratizar a beleza tipográfica.</p>
    
    <p>Cada linha de código escrita sob a Constituição Vértice v3.0 é um testemunho de qualidade e excelência. Não aceitamos atalhos, não toleramos placeholders, e sempre honramos o conteúdo com o design que ele merece.</p>
    
    <h2>A Filosofia Tipográfica</h2>
    
    <p>Como Robert Bringhurst ensina, "a tipografia existe para honrar o conteúdo". Este princípio guia cada decisão do nosso sistema de IA, desde a escolha de fontes até o ajuste microtipográfico.</p>
    
    <blockquote>
        A tipografia não é uma forma de autoexpressão do designer, mas um ato de serviço ao texto e ao leitor.
    </blockquote>
    
    <h1>Capítulo 2: A Implementação</h1>
    
    <p>O pipeline HTML/CSS com Paged.js representa o coração do nosso sistema para layouts visualmente ricos. Cada elemento foi cuidadosamente implementado seguindo os princípios do Blueprint:</p>
    
    <p>O Cânone de Van de Graaf define as proporções harmoniosas da página. As margens seguem a razão 2:3:4:6, criando equilíbrio visual e funcional.</p>
    
    <p>O grid de Müller-Brockmann organiza o conteúdo com clareza e consistência. Cada módulo do grid serve a um propósito específico na hierarquia visual.</p>
    
    <h2>Tecnologias Integradas</h2>
    
    <p>Paged.js traz o poder do CSS Paged Media para qualquer navegador moderno. Running headers, page numbers e breaks são todos controlados via CSS declarativo.</p>
    
    <p>Font subsetting via fonttools reduz o tamanho dos arquivos em 60% ou mais, mantendo apenas os glifos realmente usados no documento.</p>
    
    <h1>Capítulo 3: A Vitória</h1>
    
    <p>Seguimos metodicamente, passo a passo, conforme o planejamento. Não nos desviamos do CAMINHO. E alcançamos a vitória em nome de Jesus.</p>
    
    <p>Este PDF que você está lendo foi gerado automaticamente pelo Typecraft, demonstrando que beleza e automação podem coexistir em perfeita harmonia.</p>
    
    <p>A qualidade profissional não é mais privilégio de grandes editoras. Com Typecraft, cada autor pode publicar com a mesma excelência tipográfica.</p>
    
    <h2>Próximos Passos</h2>
    
    <p>A jornada continua. IA para design de cores, sugestão de fontes, e geração de ePub 3 estão no horizonte. Mas hoje, celebramos esta vitória.</p>
    
    <p>Para a honra e glória de Jesus. Amém.</p>
</body>
</html>`

	htmlPath := filepath.Join(outputDir, "example.html")
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		log.Fatalf("❌ Erro ao criar HTML: %v", err)
	}
	fmt.Printf("✅ HTML criado: %s\n", htmlPath)

	// Renderizar para PDF
	fmt.Println("\n📄 Iniciando renderização PDF...")
	renderer := html.NewPagedJSRenderer()

	opts := html.RenderOptions{
		HTMLPath:   htmlPath,
		OutputPath: filepath.Join(outputDir, "example.pdf"),
		Timeout:    60 * time.Second,
	}

	ctx := context.Background()
	startTime := time.Now()

	if err := renderer.RenderToPDF(ctx, opts); err != nil {
		log.Fatalf("❌ Erro ao renderizar PDF: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ PDF gerado com sucesso: %s\n", opts.OutputPath)
	fmt.Printf("⏱️  Tempo de renderização: %.2f segundos\n", duration.Seconds())

	// Verificar tamanho do arquivo
	info, err := os.Stat(opts.OutputPath)
	if err != nil {
		log.Fatalf("❌ Erro ao verificar arquivo: %v", err)
	}
	fmt.Printf("📊 Tamanho do PDF: %.2f KB\n", float64(info.Size())/1024)

	fmt.Println("\n🎉 Exemplo concluído com sucesso!")
	fmt.Println("📖 Abra o arquivo para visualizar o resultado.")
}
