package html

import (
	"fmt"
)

// CanonDimensions representa as dimensões calculadas pelo Cânone de Van de Graaf
type CanonDimensions struct {
	PageWidth       float64 // Largura da página (polegadas)
	PageHeight      float64 // Altura da página (polegadas)
	TextBlockWidth  float64 // Largura da mancha de texto
	TextBlockHeight float64 // Altura da mancha de texto
	MarginInner     float64 // Margem interna
	MarginTop       float64 // Margem superior
	MarginOuter     float64 // Margem externa
	MarginBottom    float64 // Margem inferior
}

// CalculateVanDeGraaf implementa o Cânone de Van de Graaf para definir
// proporções harmoniosas de página e margens.
//
// O Cânone de Van de Graaf é uma construção geométrica que cria proporções
// visualmente agradáveis independentemente do tamanho da página. A mancha de
// texto resultante tem a mesma proporção que a página.
//
// Para uma página 2:3, as margens seguem a proporção 2:3:4:6
// (inner:top:outer:bottom)
//
// Referência: Blueprint Seção 2.1 - "O Legado dos Manuscritos"
func CalculateVanDeGraaf(pageWidth, pageHeight float64) CanonDimensions {
	// Altura da mancha de texto = largura da página (propriedade do cânone)
	textBlockHeight := pageWidth

	// Largura da mancha de texto mantém a proporção da página
	pageRatio := pageHeight / pageWidth
	textBlockWidth := textBlockHeight / pageRatio

	// Calcular margens
	// Margem superior = (pageHeight - textBlockHeight) / (1 + pageRatio)
	marginTop := (pageHeight - textBlockHeight) / (1 + pageRatio)

	// Margem inferior = marginTop * pageRatio
	marginBottom := marginTop * pageRatio

	// Margem interna = (pageWidth - textBlockWidth) / (1 + pageRatio)
	marginInner := (pageWidth - textBlockWidth) / (1 + pageRatio)

	// Margem externa = marginInner * pageRatio
	marginOuter := marginInner * pageRatio

	return CanonDimensions{
		PageWidth:       pageWidth,
		PageHeight:      pageHeight,
		TextBlockWidth:  textBlockWidth,
		TextBlockHeight: textBlockHeight,
		MarginInner:     marginInner,
		MarginTop:       marginTop,
		MarginOuter:     marginOuter,
		MarginBottom:    marginBottom,
	}
}

// ToCSS converte as dimensões do Cânone para regras CSS
func (d CanonDimensions) ToCSS() string {
	return fmt.Sprintf(`
@page {
  size: %.2fin %.2fin;
  margin-top: %.2fin;
  margin-right: %.2fin;
  margin-bottom: %.2fin;
  margin-left: %.2fin;
}

.text-block {
  width: %.2fin;
  height: %.2fin;
}
`,
		d.PageWidth, d.PageHeight,
		d.MarginTop, d.MarginOuter, d.MarginBottom, d.MarginInner,
		d.TextBlockWidth, d.TextBlockHeight,
	)
}

// CommonPageSizes define tamanhos de página padrão
var CommonPageSizes = map[string][2]float64{
	"5x8":         {5.0, 8.0},
	"5.5x8.5":     {5.5, 8.5},
	"6x9":         {6.0, 9.0},
	"7x10":        {7.0, 10.0},
	"8.5x11":      {8.5, 11.0},
	"A4":          {8.27, 11.69},
	"A5":          {5.83, 8.27},
	"crown_quarto": {7.44, 9.69},
}

// GetPageSize retorna as dimensões de um tamanho de página padrão
func GetPageSize(name string) (width, height float64, ok bool) {
	size, exists := CommonPageSizes[name]
	if !exists {
		return 0, 0, false
	}
	return size[0], size[1], true
}
