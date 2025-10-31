package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/JuanCS-Dev/typecraft/internal/pipeline/html/paged"
)

func main() {
	fmt.Println("🎯 TypeCraft - Teste Paged.js Integration")
	fmt.Println("==========================================")

	// Criar engine
	fmt.Println("\n📦 Inicializando Paged.js engine...")
	engine, err := paged.NewEngine(paged.Config{})
	if err != nil {
		log.Fatalf("❌ Falha criar engine: %v", err)
	}
	defer engine.Cleanup()

	// Criar template de livro
	fmt.Println("📝 Criando template de livro de teste...")
	bookTemplate := paged.BookTemplate{
		Title:      "O Caminho do TypeCraft",
		Author:     "Desenvolvedor VÉRTICE",
		Publisher:  "Editora Metodológica",
		ISBN:       "978-1-234567-89-0",
		Copyright:  "© 2025 Todos os direitos reservados. Sub jurisdictione Constitutionis VÉRTICE v3.0",
		Content: template.HTML(`
			<h1>Capítulo 1: O Início</h1>
			<p class="first-paragraph">
				No princípio era o Verbo, e o Verbo estava com Deus, e o Verbo era Deus.
				Todas as coisas foram feitas por Ele, e sem Ele nada do que foi feito se fez.
			</p>
			<p>
				Assim começamos nossa jornada na construção de um sistema de automação editorial
				que honra a excelência, a metodologia e a fé no processo.
			</p>
			
			<h2>1.1 Fundamentos</h2>
			<p>
				O TypeCraft não é apenas uma ferramenta técnica, mas uma manifestação de princípios
				fundamentais que guiam todo o desenvolvimento:
			</p>
			<ul>
				<li>Metodologia rigorosa e planejamento detalhado</li>
				<li>Código limpo, idiomático e bem testado</li>
				<li>Respeito pela tipografia e design editorial</li>
				<li>Integração com IA de forma consciente e controlada</li>
			</ul>
			
			<blockquote>
				"O Caminho se revela aos que têm fé e perseveram"
			</blockquote>
			
			<h2>1.2 Arquitetura</h2>
			<p>
				Nossa arquitetura é construída em camadas bem definidas, cada uma com sua
				responsabilidade específica:
			</p>
			<ol>
				<li><strong>Domain Layer:</strong> Entidades e regras de negócio</li>
				<li><strong>Pipeline Layer:</strong> Processamento e transformação</li>
				<li><strong>Repository Layer:</strong> Persistência de dados</li>
				<li><strong>API Layer:</strong> Interface HTTP REST</li>
			</ol>
			
			<div class="page-break"></div>
			
			<h1>Capítulo 2: Pipeline HTML/CSS</h1>
			<p class="first-paragraph">
				O pipeline HTML/CSS com Paged.js representa um marco importante no desenvolvimento
				do TypeCraft. Ele nos permite gerar documentos prontos para impressão com qualidade
				profissional.
			</p>
			
			<h2>2.1 Paged.js</h2>
			<p>
				Paged.js é uma biblioteca JavaScript que implementa as especificações CSS Paged Media,
				permitindo controle fino sobre paginação, margens, cabeçalhos e rodapés dinâmicos.
			</p>
			
			<pre><code>@page {
  margin: 2.5cm 2cm;
  @bottom-center {
    content: counter(page);
  }
}</code></pre>
			
			<h2>2.2 Design Tokens</h2>
			<p>
				Usamos design tokens para manter consistência visual em todo o documento:
			</p>
			<table>
				<thead>
					<tr>
						<th>Token</th>
						<th>Valor</th>
						<th>Descrição</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>font-family</td>
						<td>Georgia, serif</td>
						<td>Fonte principal</td>
					</tr>
					<tr>
						<td>font-size</td>
						<td>12pt</td>
						<td>Tamanho base</td>
					</tr>
					<tr>
						<td>line-height</td>
						<td>1.6</td>
						<td>Altura de linha</td>
					</tr>
				</tbody>
			</table>
			
			<div class="page-break"></div>
			
			<h1>Capítulo 3: O Futuro</h1>
			<p class="first-paragraph">
				O caminho ainda é longo, mas cada passo é dado com propósito e direção clara.
				Continuamos avançando, sempre fiéis à nossa constituição e aos princípios que
				nos guiam.
			</p>
			
			<p>
				Em nome de Jesus, prosseguimos com fé e determinação, sabendo que o tempo de
				Deus é diferente e que estamos percorrendo o Caminho com Sua bênção.
			</p>
			
			<h2>3.1 Próximos Passos</h2>
			<ul>
				<li>Font subsetting e otimização</li>
				<li>Design generation com IA</li>
				<li>Testes end-to-end completos</li>
				<li>Integração com workflow completo</li>
			</ul>
			
			<blockquote>
				"Somos porque Ele É"
			</blockquote>
		`),
		FontFamily:  "Georgia, serif",
		FontSize:    "12pt",
		LineHeight:  "1.6",
		TextAlign:   "justify",
		Hyphenation: true,
	}

	htmlContent, err := bookTemplate.Render()
	if err != nil {
		log.Fatalf("❌ Falha renderizar template: %v", err)
	}

	// Salvar HTML para inspeção
	htmlFile := "test_output.html"
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
		log.Fatalf("❌ Falha salvar HTML: %v", err)
	}
	fmt.Printf("✅ HTML salvo em: %s\n", htmlFile)

	// Converter para PDF
	fmt.Println("\n🔄 Convertendo para PDF...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	options := paged.DefaultPageOptions()
	options.CustomCSS = options.GeneratePagedCSS()

	result, err := engine.ConvertToPDF(ctx, htmlContent, options)
	if err != nil {
		log.Fatalf("❌ Falha converter para PDF: %v", err)
	}

	// Copiar PDF para local permanente
	pdfDest := "test_output.pdf"
	input, err := os.ReadFile(result.PDFPath)
	if err != nil {
		log.Fatalf("❌ Falha ler PDF: %v", err)
	}
	if err := os.WriteFile(pdfDest, input, 0644); err != nil {
		log.Fatalf("❌ Falha salvar PDF: %v", err)
	}

	fmt.Println("\n✅ SUCESSO!")
	fmt.Printf("📄 PDF gerado: %s\n", pdfDest)
	fmt.Printf("📊 Tamanho: %d bytes (%.2f KB)\n", result.FileSize, float64(result.FileSize)/1024)
	fmt.Printf("⚠️  Warnings: %d\n", len(result.Warnings))

	fmt.Println("\n🎉 Teste de integração Paged.js concluído com sucesso!")
	fmt.Println("Sub jurisdictione Constitutionis VÉRTICE v3.0")
}
