package html

import (
	"math"
	"testing"
)

func TestCalculateVanDeGraaf(t *testing.T) {
	tests := []struct {
		name       string
		width      float64
		height     float64
		wantRatio  string
	}{
		{
			name:      "6x9 inch book",
			width:     6.0,
			height:    9.0,
			wantRatio: "2:3",
		},
		{
			name:      "5x8 inch book",
			width:     5.0,
			height:    8.0,
			wantRatio: "5:8",
		},
		{
			name:      "A4 page",
			width:     8.27,
			height:    11.69,
			wantRatio: "1:√2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canon := CalculateVanDeGraaf(tt.width, tt.height)

			// Validar que a mancha de texto tem a mesma proporção da página
			pageRatio := tt.height / tt.width
			textRatio := canon.TextBlockHeight / canon.TextBlockWidth
			
			if math.Abs(pageRatio-textRatio) > 0.01 {
				t.Errorf("Proporção da mancha de texto não corresponde à página. Página: %.3f, Texto: %.3f",
					pageRatio, textRatio)
			}

			// Validar que altura da mancha = largura da página (propriedade do cânone)
			if math.Abs(canon.TextBlockHeight-tt.width) > 0.01 {
				t.Errorf("Altura da mancha de texto deveria ser %.2f, obteve %.2f",
					tt.width, canon.TextBlockHeight)
			}

			// Validar que margens somam corretamente
			totalWidth := canon.MarginInner + canon.TextBlockWidth + canon.MarginOuter
			if math.Abs(totalWidth-tt.width) > 0.01 {
				t.Errorf("Margens horizontais não somam à largura da página. Esperado: %.2f, Obtido: %.2f",
					tt.width, totalWidth)
			}

			totalHeight := canon.MarginTop + canon.TextBlockHeight + canon.MarginBottom
			if math.Abs(totalHeight-tt.height) > 0.01 {
				t.Errorf("Margens verticais não somam à altura da página. Esperado: %.2f, Obtido: %.2f",
					tt.height, totalHeight)
			}

			// Validar proporção das margens para página 2:3
			if tt.width == 6.0 && tt.height == 9.0 {
				// Para 2:3, margens devem seguir 2:3:4:6
				expectedRatio := 2.0 / 3.0
				marginRatio := canon.MarginInner / canon.MarginTop
				
				if math.Abs(marginRatio-expectedRatio) > 0.05 {
					t.Errorf("Proporção inner/top deveria ser ~%.3f (2/3), obteve %.3f",
						expectedRatio, marginRatio)
				}
			}
		})
	}
}

func TestCommonPageSizes(t *testing.T) {
	tests := []string{"6x9", "5x8", "5.5x8.5", "A4", "A5"}
	
	for _, size := range tests {
		t.Run(size, func(t *testing.T) {
			width, height, ok := GetPageSize(size)
			if !ok {
				t.Fatalf("Tamanho de página '%s' não encontrado", size)
			}

			if width <= 0 || height <= 0 {
				t.Errorf("Dimensões inválidas: %.2f x %.2f", width, height)
			}

			// Calcular Van de Graaf
			canon := CalculateVanDeGraaf(width, height)

			// Gerar CSS
			css := canon.ToCSS()
			if len(css) == 0 {
				t.Error("CSS gerado está vazio")
			}

			t.Logf("Página %s: %.2f\" x %.2f\"", size, width, height)
			t.Logf("Mancha de texto: %.2f\" x %.2f\"", canon.TextBlockWidth, canon.TextBlockHeight)
			t.Logf("Margens (I/T/O/B): %.2f / %.2f / %.2f / %.2f",
				canon.MarginInner, canon.MarginTop, canon.MarginOuter, canon.MarginBottom)
		})
	}
}

func TestCanonToCSS(t *testing.T) {
	canon := CalculateVanDeGraaf(6.0, 9.0)
	css := canon.ToCSS()

	// Verificar que CSS contém elementos esperados
	expected := []string{
		"@page",
		"size:",
		"margin-top:",
		"margin-right:",
		"margin-bottom:",
		"margin-left:",
		".text-block",
		"width:",
		"height:",
	}

	for _, exp := range expected {
		if !contains(css, exp) {
			t.Errorf("CSS não contém string esperada: %s", exp)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
