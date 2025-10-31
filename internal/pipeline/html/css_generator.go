package html

import (
	"bytes"
	"fmt"
	"text/template"
)

// CSSConfig agrupa toda a configuração de estilo
type CSSConfig struct {
	Canon       CanonDimensions
	Grid        GridSystem
	FontFamily  string
	FontSize    float64
	LineHeight  float64
	Colors      ColorPalette
	Typography  TypographyScale
}

// ColorPalette define a paleta de cores
type ColorPalette struct {
	Primary   string // Cor primária
	Secondary string // Cor secundária
	Accent    string // Cor de destaque
	Text      string // Cor do texto
	Background string // Cor de fundo
}

// TypographyScale define a escala tipográfica
type TypographyScale struct {
	BaseSize float64 // Tamanho base em pt
	Ratio    float64 // Razão da escala (ex: 1.2 = terça menor, 1.618 = áurea)
	H1       float64 // Calculado: BaseSize * Ratio^4
	H2       float64 // Calculado: BaseSize * Ratio^3
	H3       float64 // Calculado: BaseSize * Ratio^2
	H4       float64 // Calculado: BaseSize * Ratio^1
	Body     float64 // = BaseSize
}

// NewTypographyScale cria uma escala tipográfica harmoniosa
//
// Referência: Blueprint Seção 2.3 - "Ritmo, Proporção e Harmonia"
// "Os tamanhos de fonte seguirão uma escala tipográfica baseada em
// uma proporção musical (ex: 1.2, a terça menor; ou 1.618, a Seção Áurea)"
func NewTypographyScale(baseSize, ratio float64) TypographyScale {
	return TypographyScale{
		BaseSize: baseSize,
		Ratio:    ratio,
		Body:     baseSize,
		H4:       baseSize * ratio,
		H3:       baseSize * ratio * ratio,
		H2:       baseSize * ratio * ratio * ratio,
		H1:       baseSize * ratio * ratio * ratio * ratio,
	}
}

// DefaultColorPalette retorna uma paleta neutra
func DefaultColorPalette() ColorPalette {
	return ColorPalette{
		Primary:   "#2c3e50",
		Secondary: "#34495e",
		Accent:    "#3498db",
		Text:      "#2c3e50",
		Background: "#ffffff",
	}
}

// CSSGenerator gera CSS dinâmico para o livro
type CSSGenerator struct {
	Config CSSConfig
}

// NewCSSGenerator cria um gerador de CSS
func NewCSSGenerator(config CSSConfig) *CSSGenerator {
	return &CSSGenerator{Config: config}
}

// Generate gera o CSS completo
func (g *CSSGenerator) Generate() (string, error) {
	tmpl := template.Must(template.New("css").Parse(cssTemplate))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, g.Config)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar CSS: %w", err)
	}

	return buf.String(), nil
}

// cssTemplate é o template CSS base
const cssTemplate = `
/* ========================================
   TYPECRAFT - AUTO-GENERATED CSS
   Baseado nos princípios de Van de Graaf,
   Müller-Brockmann e Bringhurst
   ======================================== */

/* === PAGE SETUP (Van de Graaf Canon) === */
{{.Canon.ToCSS}}

/* === GRID SYSTEM (Müller-Brockmann) === */
{{.Grid.ToCSS}}

/* === TYPOGRAPHY SCALE === */
:root {
  --font-family: {{.FontFamily}};
  --font-size-base: {{.Typography.Body}}pt;
  --line-height-base: {{.LineHeight}};
  
  /* Scale */
  --font-size-h1: {{.Typography.H1}}pt;
  --font-size-h2: {{.Typography.H2}}pt;
  --font-size-h3: {{.Typography.H3}}pt;
  --font-size-h4: {{.Typography.H4}}pt;
  --font-size-body: {{.Typography.Body}}pt;
  
  /* Colors */
  --color-primary: {{.Colors.Primary}};
  --color-secondary: {{.Colors.Secondary}};
  --color-accent: {{.Colors.Accent}};
  --color-text: {{.Colors.Text}};
  --color-background: {{.Colors.Background}};
}

/* === BASE STYLES === */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  line-height: var(--line-height-base);
  color: var(--color-text);
  background-color: var(--color-background);
}

/* === TYPOGRAPHY === */
h1 {
  font-size: var(--font-size-h1);
  margin-bottom: calc(var(--line-height-base) * 2em);
  color: var(--color-primary);
}

h2 {
  font-size: var(--font-size-h2);
  margin-top: calc(var(--line-height-base) * 2em);
  margin-bottom: calc(var(--line-height-base) * 1em);
  color: var(--color-primary);
}

h3 {
  font-size: var(--font-size-h3);
  margin-top: calc(var(--line-height-base) * 1.5em);
  margin-bottom: calc(var(--line-height-base) * 0.75em);
}

h4 {
  font-size: var(--font-size-h4);
  margin-top: calc(var(--line-height-base) * 1em);
  margin-bottom: calc(var(--line-height-base) * 0.5em);
}

p {
  margin-bottom: calc(var(--line-height-base) * 1em);
  text-align: justify;
  hyphens: auto;
}

/* === MICROTIPOGRAFIA === */
/* Referência: Blueprint Seção 3.1 - "O Espaçamento Fino" */
p:first-child {
  text-indent: 0;
}

p + p {
  text-indent: 1.5em;
}

/* Evitar viúvas e órfãs */
/* Referência: Blueprint Seção 3.3 */
p {
  orphans: 3;
  widows: 3;
}

h1, h2, h3, h4, h5, h6 {
  page-break-after: avoid;
  orphans: 3;
  widows: 3;
}

/* === CODE BLOCKS === */
pre {
  margin: calc(var(--line-height-base) * 1em) 0;
  padding: 1em;
  background-color: #f5f5f5;
  border-left: 4px solid var(--color-accent);
  overflow-x: auto;
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
}

code {
  font-family: 'Courier New', monospace;
  background-color: #f5f5f5;
  padding: 0.2em 0.4em;
}

/* === BLOCKQUOTES === */
blockquote {
  margin: calc(var(--line-height-base) * 1em) 2em;
  padding-left: 1em;
  border-left: 4px solid var(--color-accent);
  font-style: italic;
  color: var(--color-secondary);
}

/* === LISTS === */
ul, ol {
  margin: calc(var(--line-height-base) * 1em) 0;
  padding-left: 2em;
}

li {
  margin-bottom: calc(var(--line-height-base) * 0.5em);
}

/* === IMAGES === */
img {
  max-width: 100%;
  height: auto;
  display: block;
  margin: calc(var(--line-height-base) * 1em) auto;
}

figure {
  margin: calc(var(--line-height-base) * 1.5em) 0;
  text-align: center;
}

figcaption {
  font-size: 0.9em;
  font-style: italic;
  color: var(--color-secondary);
  margin-top: 0.5em;
}

/* === TABLES === */
table {
  width: 100%;
  margin: calc(var(--line-height-base) * 1em) 0;
  border-collapse: collapse;
}

th, td {
  padding: 0.5em;
  border: 1px solid #ddd;
  text-align: left;
}

th {
  background-color: var(--color-primary);
  color: white;
}

/* === LINKS === */
a {
  color: var(--color-accent);
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

/* === PRINT-SPECIFIC === */
@media print {
  body {
    background-color: white;
  }
  
  a {
    color: var(--color-text);
  }
}
`
