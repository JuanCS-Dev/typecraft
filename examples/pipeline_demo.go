package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JuanCS-Dev/typecraft/pkg/ai"
	"github.com/JuanCS-Dev/typecraft/pkg/pipeline"
	"github.com/JuanCS-Dev/typecraft/pkg/typography"
)

func main() {
	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Println("║       📚 TypeCraft Pipeline Demo        ║")
	fmt.Println("║   Automação Editorial com IA             ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Println()

	// 1. Configurar Style Engine
	styleEngine := typography.NewStyleEngine()
	
	// 2. Configurar AI Client (opcional)
	var aiClient *ai.Client
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		aiClient = ai.NewClient(apiKey, "gpt-4o-mini", 2000, 0.7)
		fmt.Println("✅ Cliente AI configurado")
	} else {
		fmt.Println("⚠️  Cliente AI não configurado (OPENAI_API_KEY não definida)")
	}

	// 3. Criar diretório de output
	outputDir := "./output"
	os.MkdirAll(outputDir, 0755)

	// 4. Criar pipeline
	pipe, err := pipeline.NewPipeline(styleEngine, aiClient, outputDir)
	if err != nil {
		log.Fatalf("❌ Erro ao criar pipeline: %v", err)
	}
	defer pipe.Cleanup()

	// 5. Criar arquivo de exemplo
	exampleFile := "./example_chapter.txt"
	exampleContent := `Era uma vez, em um reino distante, um jovem programador que descobriu a arte da tipografia digital.

Ele aprendeu que a beleza de um livro não está apenas nas palavras, mas em como elas são apresentadas na página.

"A tipografia é a voz visual da palavra escrita", pensava ele, enquanto ajustava cuidadosamente cada espaço e cada letra.

E assim começou sua jornada...`

	if err := os.WriteFile(exampleFile, []byte(exampleContent), 0644); err != nil {
		log.Fatalf("❌ Erro ao criar arquivo de exemplo: %v", err)
	}
	defer os.Remove(exampleFile)

	fmt.Println("📝 Arquivo de exemplo criado")
	fmt.Println()

	// 6. Configurar processamento
	config := pipeline.ProcessBookConfig{
		InputFiles: []string{exampleFile},
		Title:      "A Jornada do Tipógrafo Digital",
		Author:     "TypeCraft Demo",
		DesignPrompt: "Design elegante e clássico, inspirado em livros do século XIX. " +
			"Use cores sóbrias e tipografia serif tradicional. " +
			"Adicione ornamentos sutis.",
		PageSize: "A5",
	}

	// 7. Processar livro
	result, err := pipe.ProcessBook(config)
	if err != nil {
		log.Fatalf("❌ Erro ao processar livro: %v", err)
	}

	// 8. Exibir resultado
	fmt.Println(result.Report())

	if result.Success {
		fmt.Println("🎉 Livro processado com sucesso!")
		fmt.Println()
		fmt.Println("📂 Arquivos disponíveis em:", outputDir)
		fmt.Println("   - book.html (visualize no navegador)")
		fmt.Println("   - book.pdf (pronto para impressão)")
	}
}
