package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FontSubsetter realiza subsetting de fontes para otimização
type FontSubsetter struct {
	pyftsubsetPath string
}

// NewFontSubsetter cria um novo subsetter de fontes
func NewFontSubsetter() (*FontSubsetter, error) {
	// Verifica se pyftsubset está disponível
	path, err := exec.LookPath("pyftsubset")
	if err != nil {
		return nil, fmt.Errorf("pyftsubset não encontrado. Instale fonttools: pip install fonttools")
	}

	return &FontSubsetter{
		pyftsubsetPath: path,
	}, nil
}

// SubsetFont cria subset de uma fonte baseado no texto usado
func (f *FontSubsetter) SubsetFont(fontPath string, text string, outputPath string) error {
	// Remove duplicatas do texto
	uniqueChars := uniqueRunes(text)

	// Cria arquivo temporário com os caracteres
	tempFile, err := os.CreateTemp("", "subset-chars-*.txt")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(uniqueChars); err != nil {
		return fmt.Errorf("erro ao escrever caracteres: %w", err)
	}
	tempFile.Close()

	// Executa pyftsubset
	cmd := exec.Command(f.pyftsubsetPath,
		fontPath,
		fmt.Sprintf("--text-file=%s", tempFile.Name()),
		fmt.Sprintf("--output-file=%s", outputPath),
		"--flavor=woff2",
		"--layout-features=*",
		"--no-hinting",
		"--desubroutinize",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao executar pyftsubset: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// SubsetFontFamily cria subsets de uma família de fontes completa
func (f *FontSubsetter) SubsetFontFamily(fontDir string, text string, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de saída: %w", err)
	}

	// Lista arquivos de fonte
	entries, err := os.ReadDir(fontDir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de fontes: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".ttf" && ext != ".otf" {
			continue
		}

		inputPath := filepath.Join(fontDir, entry.Name())
		baseName := strings.TrimSuffix(entry.Name(), ext)
		outputPath := filepath.Join(outputDir, baseName+".woff2")

		if err := f.SubsetFont(inputPath, text, outputPath); err != nil {
			return fmt.Errorf("erro ao processar %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// ExtractTextFromHTML extrai texto de HTML para análise de caracteres
func (f *FontSubsetter) ExtractTextFromHTML(htmlContent string) string {
	// Remove tags HTML (simplificado)
	text := htmlContent
	text = removeHTMLTags(text)
	return text
}

// uniqueRunes retorna string com caracteres únicos
func uniqueRunes(s string) string {
	seen := make(map[rune]bool)
	var result []rune

	for _, r := range s {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}

	return string(result)
}

// removeHTMLTags remove tags HTML básicas (simplificado)
func removeHTMLTags(html string) string {
	var result strings.Builder
	inTag := false

	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// GenerateFontFaceCSS gera CSS @font-face para fontes subsetted
func (f *FontSubsetter) GenerateFontFaceCSS(fontDir string, fontFamily string) (string, error) {
	entries, err := os.ReadDir(fontDir)
	if err != nil {
		return "", fmt.Errorf("erro ao ler diretório: %w", err)
	}

	var css strings.Builder
	css.WriteString("/* Auto-generated font-face declarations */\n\n")

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".woff2") {
			continue
		}

		// Detecta weight e style pelo nome do arquivo
		weight := "400"
		style := "normal"

		name := strings.ToLower(entry.Name())
		if strings.Contains(name, "bold") {
			weight = "700"
		} else if strings.Contains(name, "light") {
			weight = "300"
		}

		if strings.Contains(name, "italic") {
			style = "italic"
		}

		css.WriteString(fmt.Sprintf(`@font-face {
    font-family: '%s';
    src: url('%s') format('woff2');
    font-weight: %s;
    font-style: %s;
    font-display: swap;
}

`, fontFamily, entry.Name(), weight, style))
	}

	return css.String(), nil
}
