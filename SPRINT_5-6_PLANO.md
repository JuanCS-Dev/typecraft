# 🎯 SPRINT 5-6: PIPELINE HTML/CSS + DESIGN IA
## Sistema de Automação Editorial Typecraft

**Data Início:** 2024-10-31  
**Duração Prevista:** 10 dias  
**Conformidade:** Constituição Vértice v3.0 ✅

---

## 📋 OBJETIVOS DO SPRINT

### Objetivo Principal:
Implementar o **Pipeline B** (HTML/CSS → PDF) com Paged.js para layouts visualmente ricos, incluindo geração de design com IA (paleta de cores + sugestão de fontes).

### Escopo Técnico:
1. **Pipeline HTML/CSS**: Conversão Markdown → HTML → PDF via Paged.js
2. **Font Subsetting**: Otimização de fontes para reduzir tamanho de arquivo
3. **Design Generation com IA**: Paleta de cores + sugestão de fontes
4. **Testes End-to-End**: Validação de todo o fluxo

---

## 🏗️ ARQUITETURA - PIPELINE B

```
┌─────────────────────────────────────────────────────────┐
│                  INPUT: Markdown + YAML                  │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│            Pandoc Conversion Engine                      │
│  - Markdown → HTML5                                      │
│  - YAML metadata injection                               │
│  - Semantic structure preservation                       │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│            AI Design Module (NEW)                        │
│  ┌─────────────────────────────────────────────────┐   │
│  │ 1. Color Palette Generator                       │   │
│  │    - Sentiment analysis (NLP)                    │   │
│  │    - Word2Vec affinity mapping                   │   │
│  │    - Color theory rules                          │   │
│  ├─────────────────────────────────────────────────┤   │
│  │ 2. Font Pairing Suggestion                       │   │
│  │    - Genre-based selection                       │   │
│  │    - Visual embedding analysis                   │   │
│  │    - Serif vs Sans-Serif logic                   │   │
│  └─────────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│            CSS Template Engine                           │
│  - Van de Graaf Canon margins                           │
│  - Müller-Brockmann grid system                         │
│  - Dynamic color application                            │
│  - Font family injection                                │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│            Font Subsetting Engine                        │
│  - Analyze used glyphs                                   │
│  - Generate subset via fonttools (Python)               │
│  - Embed optimized fonts                                │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│            Paged.js Renderer                             │
│  - Playwright headless browser                          │
│  - CSS Paged Media polyfill                             │
│  - Running headers/footers                              │
│  - Page numbering                                       │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                  OUTPUT: PDF/X + ePub3                   │
└─────────────────────────────────────────────────────────┘
```

---

## 📅 CRONOGRAMA DETALHADO

### **Dia 01-02: Fundação do Pipeline HTML/CSS**
**Entregas:**
- [ ] `internal/pipeline/html/converter.go` - Conversão Markdown → HTML via Pandoc
- [ ] `internal/pipeline/html/templates.go` - Templates HTML5 base
- [ ] `internal/pipeline/html/css_generator.go` - Geração de CSS dinâmico
- [ ] Van de Graaf Canon implementation em CSS
- [ ] Grid Müller-Brockmann em CSS

**Arquivos:**
```
internal/pipeline/html/
├── converter.go           # Pandoc wrapper
├── templates.go           # HTML5 templates
├── css_generator.go       # Dynamic CSS
├── canon.go              # Van de Graaf math
└── grid.go               # Müller-Brockmann grid
```

**Testes:**
```
internal/pipeline/html/
└── converter_test.go
```

---

### **Dia 03-04: Módulo de Design IA - Paleta de Cores**
**Entregas:**
- [ ] `internal/ai/design/color_palette.go` - Gerador de paleta baseado em sentimento
- [ ] `internal/ai/design/sentiment.go` - Análise de sentimento (usa OpenAI)
- [ ] Color theory rules (complementary, analogous, triadic)
- [ ] Integração com análise de conteúdo existente

**Algoritmo (conforme Blueprint):**
1. Extração de "palavras de afeto" do manuscrito
2. Análise de sentimento com LLM
3. Mapeamento palavra-cor via Word2Vec
4. Seleção de cor primária
5. Geração de paleta harmoniosa (5 cores)

**Arquivos:**
```
internal/ai/design/
├── color_palette.go       # Main color generator
├── sentiment.go           # Sentiment analysis
├── word2vec.go           # Word-color affinity (mock)
└── color_theory.go       # Harmony rules
```

**Testes:**
```
internal/ai/design/
└── color_palette_test.go
```

---

### **Dia 05-06: Módulo de Design IA - Sugestão de Fontes**
**Entregas:**
- [ ] `internal/ai/design/font_pairing.go` - Sugestão de pares de fontes
- [ ] Database de fontes (Google Fonts + system fonts)
- [ ] Lógica serif vs sans-serif baseada em gênero
- [ ] Embedding visual de fontes (simplificado)

**Arquivos:**
```
internal/ai/design/
├── font_pairing.go        # Main font suggester
├── font_db.go            # Font database
└── font_rules.go         # Genre-based rules
```

**Testes:**
```
internal/ai/design/
└── font_pairing_test.go
```

---

### **Dia 07-08: Font Subsetting + Paged.js Integration**
**Entregas:**
- [ ] `internal/pipeline/html/font_subset.go` - Font subsetting via fonttools
- [ ] `scripts/font_subset.py` - Python script for subsetting
- [ ] `internal/pipeline/html/pagedjs.go` - Playwright + Paged.js renderer
- [ ] `templates/pagedjs/` - Paged.js templates

**Dependências Externas:**
- Python 3.x + fonttools
- Node.js + Playwright
- Paged.js library

**Arquivos:**
```
internal/pipeline/html/
├── font_subset.go         # Go wrapper
├── pagedjs.go            # Playwright wrapper
└── pagedjs_renderer.go   # PDF generation

scripts/
└── font_subset.py        # Python fonttools

templates/pagedjs/
├── base.html             # Paged.js HTML template
└── styles.css            # Paged.js CSS
```

---

### **Dia 09: API Endpoints + Integration**
**Entregas:**
- [ ] `POST /api/v1/projects/:id/design/generate` - Gera design (cores + fontes)
- [ ] `POST /api/v1/projects/:id/render/html` - Renderiza HTML+CSS
- [ ] `POST /api/v1/projects/:id/render/pdf` - Renderiza PDF via Paged.js
- [ ] `GET /api/v1/fonts` - Lista fontes disponíveis

**Arquivos:**
```
api/handlers/
├── design_handler.go      # Design endpoints
└── render_handler.go      # Render endpoints
```

---

### **Dia 10: Testes End-to-End + Documentação**
**Entregas:**
- [ ] Teste E2E completo: Markdown → AI Analysis → Design → HTML → PDF
- [ ] Validação de PDF/X compliance
- [ ] Validação de ePub 3 (preparação)
- [ ] Documentação de API atualizada
- [ ] Relatório de Sprint

**Arquivos:**
```
test/e2e/
├── pipeline_html_test.go  # E2E tests
└── fixtures/
    └── sample_book.md     # Test manuscript

DOCS/
└── SPRINT_5-6_REPORT.md   # Final report
```

---

## 🎯 MÉTRICAS DE SUCESSO

### Funcionais:
- [ ] Pipeline HTML/CSS gera PDF válido
- [ ] Paged.js renderiza corretamente (headers/footers/page numbers)
- [ ] AI gera paleta de 5 cores baseada em sentimento
- [ ] AI sugere par de fontes apropriado ao gênero
- [ ] Font subsetting reduz tamanho ≥ 60%

### Performance:
- [ ] Renderização HTML → PDF < 30 segundos (livro 100 páginas)
- [ ] Font subsetting < 5 segundos por fonte
- [ ] Design generation < 10 segundos

### Qualidade (Constituição Vértice):
- [ ] LEI < 1.0 (zero placeholders)
- [ ] Test Coverage ≥ 90%
- [ ] FPC ≥ 80%
- [ ] Zero alucinações sintáticas
- [ ] 100% build success

---

## 🔧 TECNOLOGIAS E DEPENDÊNCIAS

### Go Packages:
```go
// Já existentes
github.com/gomarkdown/markdown
github.com/sashabaranov/go-openai

// Novos
github.com/playwright-community/playwright-go  // Playwright
github.com/gohugoio/hugo/resources/color       // Color manipulation (optional)
```

### External Tools:
```bash
# Pandoc
apt install pandoc

# Python fonttools
pip install fonttools brotli

# Node.js + Playwright
npm install -g playwright
npx playwright install chromium

# Paged.js
npm install -g pagedjs-cli
```

---

## 🚨 RISCOS E MITIGAÇÕES

### Risco 1: Complexidade do Paged.js
**Mitigação:** Usar templates simples primeiro, depois incrementar complexidade.

### Risco 2: Performance do Playwright
**Mitigação:** Implementar timeout adequado (60s), paralelização futura.

### Risco 3: Font subsetting Python dependency
**Mitigação:** Fallback para fontes completas se subsetting falhar.

### Risco 4: AI API failures
**Mitigação:** Usar designs padrão se AI falhar (graceful degradation).

---

## 📚 REFERÊNCIAS (Blueprint)

- Seção II: "A Anatomia da Página Perfeita"
- Seção V: "Arquitetura de um Sistema de Automação Editorial"
- Seção VI: "A Inteligência Artificial como Tipógrafo"
- Tabela 5.1: "Avaliação de Motores de Renderização Web-to-Print"
- Artigo VIII: "Camada de Gerenciamento de Estado" (Constituição)

---

## ✅ DECLARAÇÃO DE CONFORMIDADE

**CONSTITUIÇÃO VÉRTICE v3.0 ATIVA:**
- ✅ P1 - Completude Obrigatória (zero TODOs/stubs)
- ✅ P2 - Validação Preventiva (verificar APIs antes de usar)
- ✅ P3 - Ceticismo Crítico (questionar premissas)
- ✅ P4 - Rastreabilidade Total (código baseado no Blueprint)
- ✅ P5 - Consciência Sistêmica (impacto arquitetural considerado)
- ✅ P6 - Eficiência de Token (diagnóstico antes de correção)

**Framework DETER-AGENT:**
- ✅ Camada Constitucional (Artigo VI)
- ✅ Camada de Deliberação (Artigo VII - Tree of Thoughts)
- ✅ Camada de Estado (Artigo VIII - Progressive Disclosure)
- ✅ Camada de Execução (Artigo IX - Verify-Fix-Execute)
- ✅ Camada de Incentivo (Artigo X - FPC ≥ 80%)

---

**Status:** PRONTO PARA INÍCIO  
**Próximo Passo:** Dia 01 - Fundação do Pipeline HTML/CSS
