package ai

import (
	"fmt"
	"strings"
)

const systemPrompt = `Você é um especialista em análise literária e tipografia editorial de nível mundial.

Sua tarefa é analisar manuscritos e fornecer análise detalhada para orientar o pipeline de diagramação.

Analise:
1. GÊNERO LITERÁRIO - Classifique precisamente (fiction, non-fiction, technical, academic, poetry, etc)
2. TOM E ESTILO - Identifique o tom predominante (formal, informal, conversational, technical, narrative)
3. ELEMENTOS ESPECIAIS - Detecte e quantifique:
   - Equações matemáticas (% do texto)
   - Código-fonte (% do texto)
   - Tabelas (contagem)
   - Imagens/figuras (contagem)
4. PIPELINE RECOMENDADO - Baseado na análise:
   - "simple": Texto simples, Markdown suficiente
   - "standard": Livro padrão, pode usar Typst
   - "latex": Requer LaTeX (muitas equações/tabelas/complexidade)
5. FONTES RECOMENDADAS - Sugira fontes apropriadas para corpo, títulos e mono

Responda SEMPRE em formato JSON válido:
{
  "genre": "string",
  "sub_genre": "string (opcional)",
  "confidence": 0.0-1.0,
  "tone": "string",
  "tone_score": 0.0-1.0,
  "special_elements": {
    "equation_percentage": 0.0-1.0,
    "code_percentage": 0.0-1.0,
    "table_count": int,
    "image_count": int
  },
  "recommendations": {
    "pipeline": "simple|standard|latex",
    "pipeline_reason": "string explicando o porquê",
    "confidence": 0.0-1.0,
    "body_font": "string",
    "title_font": "string",
    "mono_font": "string",
    "font_rationale": "string explicando as escolhas"
  }
}

Seja PRECISO. Baseie percentagens em análise real do texto.`

func buildAnalysisPrompt(text string) string {
	sample := truncateText(text, 5000)
	
	return fmt.Sprintf(`Analise o seguinte manuscrito e retorne análise completa em JSON:

MANUSCRITO:
---
%s
---

TAREFAS:
1. Identifique gênero e sub-gênero com nível de confiança
2. Determine tom e estilo de escrita
3. Detecte e quantifique elementos especiais:
   - Equações matemáticas: estime %% do texto que contém fórmulas
   - Código-fonte: estime %% do texto que é código
   - Tabelas: conte quantas tabelas aparecem
   - Imagens/figuras: conte quantas referências a imagens
4. Recomende pipeline ideal:
   - "simple": Para textos sem formatação especial
   - "standard": Para livros normais com alguma formatação
   - "latex": Para conteúdo técnico/acadêmico complexo
5. Sugira fontes apropriadas para o gênero

Retorne APENAS JSON válido, sem texto adicional.`, sample)
}

func truncateText(text string, maxChars int) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxChars {
		return text
	}
	
	// Tenta truncar em quebra de linha ou ponto
	truncated := text[:maxChars]
	if lastNewline := strings.LastIndex(truncated, "\n"); lastNewline > maxChars/2 {
		return text[:lastNewline]
	}
	if lastPeriod := strings.LastIndex(truncated, "."); lastPeriod > maxChars/2 {
		return text[:lastPeriod+1]
	}
	
	return truncated + "..."
}
