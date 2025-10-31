package main

import (
"fmt"
"os"

"github.com/JuanCS-Dev/typecraft/pkg/converter"
"github.com/JuanCS-Dev/typecraft/pkg/renderer"
)

func main() {
fmt.Println("🧪 TESTE DE CONVERSORES E RENDERIZADORES")
fmt.Println()

// Teste 1: Verificar Pandoc
fmt.Println("1️⃣  Testando Pandoc Converter...")
pandoc, err := converter.NewPandocConverter()
if err != nil {
fmt.Printf("   ❌ Erro: %v\n", err)
os.Exit(1)
}

version, _ := pandoc.GetVersion()
fmt.Printf("   ✅ Pandoc encontrado: %s\n", version)

// Teste 2: Verificar LaTeX
fmt.Println()
fmt.Println("2️⃣  Testando LaTeX Renderer...")
latex, err := renderer.NewLatexRenderer()
if err != nil {
fmt.Printf("   ❌ Erro: %v\n", err)
os.Exit(1)
}

latexVersion, _ := latex.GetVersion()
fmt.Printf("   ✅ LaTeX encontrado: %s\n", latexVersion)

// Teste 3: Conversão Markdown → PDF
fmt.Println()
fmt.Println("3️⃣  Testando conversão Markdown → PDF...")

testMd := "test/poc/test_manuscript.md"
outputPdf := "test/poc/converter_test.pdf"

err = pandoc.MarkdownToPDF(testMd, outputPdf, []string{})
if err != nil {
fmt.Printf("   ❌ Erro na conversão: %v\n", err)
os.Exit(1)
}

fmt.Printf("   ✅ PDF gerado: %s\n", outputPdf)

fmt.Println()
fmt.Println("🎉 TODOS OS TESTES PASSARAM!")
}
