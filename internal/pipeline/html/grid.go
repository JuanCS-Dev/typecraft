package html

import "fmt"

// GridSystem representa um sistema de grid Müller-Brockmann
type GridSystem struct {
	Columns      int     // Número de colunas
	Rows         int     // Número de linhas (opcional)
	GutterWidth  float64 // Espaço entre colunas (em)
	GutterHeight float64 // Espaço entre linhas (em)
	BaselineGrid float64 // Grid de linha de base (em)
}

// GridType define tipos pré-configurados de grid
type GridType string

const (
	GridSingleColumn GridType = "single"  // 1 coluna - prosa simples
	GridTwoColumn    GridType = "two"     // 2 colunas - manuais, revistas
	GridThreeColumn  GridType = "three"   // 3 colunas - jornais, revistas
	GridSixColumn    GridType = "six"     // 6 colunas - layouts complexos
	GridTwelveColumn GridType = "twelve"  // 12 colunas - máxima flexibilidade
)

// NewGrid cria um sistema de grid baseado no tipo
//
// Referência: Blueprint Seção 2.2 - "A Ordem Racionalista"
// Baseado nos sistemas de Josef Müller-Brockmann
func NewGrid(gridType GridType) GridSystem {
	switch gridType {
	case GridSingleColumn:
		return GridSystem{
			Columns:      1,
			GutterWidth:  0,
			BaselineGrid: 1.5, // 1.5em = ~24pt para corpo 16pt
		}

	case GridTwoColumn:
		return GridSystem{
			Columns:      2,
			GutterWidth:  1.5, // 1.5em between columns
			BaselineGrid: 1.5,
		}

	case GridThreeColumn:
		return GridSystem{
			Columns:      3,
			GutterWidth:  1.0,
			BaselineGrid: 1.5,
		}

	case GridSixColumn:
		return GridSystem{
			Columns:      6,
			GutterWidth:  0.75,
			BaselineGrid: 1.5,
		}

	case GridTwelveColumn:
		return GridSystem{
			Columns:      12,
			GutterWidth:  0.5,
			BaselineGrid: 1.5,
		}

	default:
		return NewGrid(GridSingleColumn)
	}
}

// ToCSS gera regras CSS para o sistema de grid
func (g GridSystem) ToCSS() string {
	if g.Columns == 1 {
		// Grid simples não precisa de CSS complexo
		return fmt.Sprintf(`
.content {
  line-height: %.2fem;
}
`, g.BaselineGrid)
	}

	// Grid multi-coluna usando CSS Grid
	columnTemplate := ""
	for i := 0; i < g.Columns; i++ {
		if i > 0 {
			columnTemplate += " "
		}
		columnTemplate += "1fr"
	}

	return fmt.Sprintf(`
.grid-container {
  display: grid;
  grid-template-columns: %s;
  gap: %.2fem;
  line-height: %.2fem;
}

.grid-span-2 {
  grid-column: span 2;
}

.grid-span-3 {
  grid-column: span 3;
}

.grid-span-4 {
  grid-column: span 4;
}

.grid-span-6 {
  grid-column: span 6;
}

.grid-span-12 {
  grid-column: span 12;
}
`,
		columnTemplate,
		g.GutterWidth,
		g.BaselineGrid,
	)
}

// DetermineGridType decide o tipo de grid baseado na análise de conteúdo
//
// Referência: Blueprint Seção 2.2
// "A complexidade do grid não será arbitrária. Ela será determinada
// dinamicamente pelo módulo de análise de conteúdo da IA."
func DetermineGridType(hasImages bool, hasCode bool, hasTables bool, complexity string) GridType {
	// Conteúdo simples: apenas texto
	if !hasImages && !hasCode && !hasTables {
		return GridSingleColumn
	}

	// Conteúdo técnico com código: 2 colunas para code blocks laterais
	if hasCode && complexity == "high" {
		return GridTwoColumn
	}

	// Conteúdo com imagens e tabelas: 3-6 colunas
	if hasImages && hasTables {
		return GridSixColumn
	}

	// Conteúdo com imagens: 2-3 colunas
	if hasImages {
		return GridThreeColumn
	}

	// Default para conteúdo misto
	return GridTwoColumn
}
