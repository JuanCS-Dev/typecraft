# 🎨 SPRINT 7-8: DESIGN IA + EPUB PIPELINE

**Data Início:** 2024-10-31  
**Duração:** 10 dias planejamento + 1-2 dias execução  
**Objetivo:** Sistema inteligente de design + pipeline ePub completo

---

## 📊 CONTEXTO E FUNDAÇÃO

### Estado Atual do Sistema

**✅ Completado (Sprint 5-6):**
- Pipeline HTML/CSS com Paged.js
- Font subsetting e otimização
- Van de Graaf Canon + Müller-Brockmann Grid
- API endpoints de renderização
- Testes E2E com 100% coverage

**🎯 Gap a Preencher:**
- Sistema de seleção de fontes é baseado em regras fixas
- Cores são padrão (preto no branco)
- ePub ainda não implementado
- Falta decisão automática de pipeline (LaTeX vs HTML)

### Arquitetura Alvo

```
┌─────────────────────────────────────────────────────────────┐
│                    DESIGN INTELLIGENCE                       │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Content Analysis  →  Font Suggester  →  Color Generator    │
│        ↓                    ↓                    ↓            │
│   Genre Detection    ML Model (XGBoost)   Sentiment + HSL   │
│   Complexity Score   Google Fonts DB      Harmony Rules     │
│   Image/Table %      Typography Rules     Accessibility     │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    PIPELINE SELECTOR                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  IF complex_math OR academic → LaTeX Pipeline                │
│  IF rich_media OR web_first → HTML/CSS Pipeline              │
│  IF simple_text → User Preference (default: HTML)            │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    EPUB 3 GENERATOR                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  HTML Base → EPUB Structure → Package Files → Validation    │
│     ↓             ↓                ↓              ↓           │
│  Chapters    content.opf      toc.ncx      epubcheck        │
│  Images      metadata         nav.xhtml    compatibility    │
│  Styles      manifest         cover.xhtml  tests            │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 OBJETIVOS DO SPRINT

### Objetivo Principal
Implementar sistema de design inteligente baseado em análise de conteúdo e ML, com geração completa de ePub 3.

### KPIs de Sucesso

**Funcionalidade:**
- [ ] 3 pares de fontes sugeridos por livro (body + heading)
- [ ] Paleta de cores contextual (3-5 cores harmonicas)
- [ ] Pipeline selecionado automaticamente baseado em conteúdo
- [ ] ePub 3 válido gerado e testado

**Performance:**
- [ ] Análise de conteúdo: < 500ms para 50k palavras
- [ ] Sugestão de fontes: < 200ms (consulta ao modelo)
- [ ] Geração de paleta: < 100ms
- [ ] ePub generation: < 2s para 200 páginas

**Qualidade:**
- [ ] 80%+ das sugestões aprovadas sem modificação (teste com 20 amostras)
- [ ] ePub passa em epubcheck sem erros
- [ ] Testes E2E com 100% coverage
- [ ] Zero placeholders (LEI = 0.0)

---

## 📅 CRONOGRAMA DETALHADO

### 🔍 DIA 1-3: ANÁLISE E DESIGN (Planejamento)

#### Dia 1: Análise de Conteúdo + Dataset Preparation

**Tarefas:**
1. **Content Analyzer Module**
   ```go
   type ContentAnalysis struct {
       GenreSignals    map[string]float64  // palavra -> peso de gênero
       Complexity      float64             // 0-1 (Flesch reading ease)
       ToneMetrics     ToneProfile         // formal, casual, technical
       ImageCount      int
       TableCount      int
       EquationCount   int
       WordCount       int
       SentimentScore  float64             // -1 a +1
   }
   ```

2. **Dataset: Google Fonts Metadata**
   - Baixar metadata de todas as fontes Google Fonts (~1500)
   - Categorizar: serif, sans-serif, display, monospace
   - Anotar: altura-x, peso, largura, mood
   - Formato JSON para consulta rápida

3. **Dataset: Genre-Font Pairings**
   - Criar dataset sintético baseado em regras de design estabelecidas
   - 50 pares anotados: {genre, mood, body_font, heading_font, rationale}

**Entregáveis:**
- [ ] `internal/analyzer/content.go` - análise de manuscrito
- [ ] `pkg/fonts/google_fonts_db.json` - metadata de 1500 fontes
- [ ] `pkg/design/genre_font_pairs.json` - dataset de treinamento

---

#### Dia 2: Font Suggestion ML Model

**Abordagem: XGBoost Classifier**

**Features (8):**
1. Genre vector (one-hot encoded: 6 categorias)
2. Complexity score (0-1)
3. Formality score (0-1)
4. Word count (normalized)
5. Image ratio (images/pages)
6. Sentence length avg
7. Technical term density
8. Sentiment polarity

**Labels:**
- 30 combinações de fontes curadas (top pairings)

**Training:**
```python
import xgboost as xgb
import pandas as pd

# Synthetic dataset generation
# 1000 samples: combinações de features → font pairs
data = generate_training_data()

# Train
model = xgb.XGBClassifier(
    max_depth=6,
    learning_rate=0.1,
    n_estimators=100
)
model.fit(X_train, y_train)

# Export to ONNX for Go inference
import onnxmltools
onnx_model = onnxmltools.convert_xgboost(model)
```

**Go Integration:**
```go
import "github.com/owulveryck/onnx-go"

type FontSuggester struct {
    model *onnx.Model
    fontDB map[string]FontMetadata
}

func (fs *FontSuggester) Suggest(analysis ContentAnalysis) []FontPair
```

**Entregáveis:**
- [ ] `scripts/train_font_model.py` - treinar XGBoost
- [ ] `pkg/design/font_model.onnx` - modelo exportado
- [ ] `pkg/design/font_suggester.go` - inferência em Go
- [ ] `test/design/font_suggester_test.go` - testes unitários

---

#### Dia 3: Color Palette Generation

**Algoritmo: Sentiment → Hue + Harmony Rules**

**Pipeline:**
```
1. Extract keywords (top 50 by TF-IDF)
2. Sentiment analysis (via GPT-4o-mini):
   - Emotional valence: joy, sadness, fear, etc.
   - Energy level: calm, intense, neutral
3. Map to HSL:
   - Joy → Yellow/Orange (H: 40-60)
   - Sadness → Blue (H: 200-240)
   - Fear → Purple/Black (H: 270-300)
   - Calm → Green/Blue (H: 120-200)
   - Intense → Red/Orange (H: 0-40)
4. Generate palette:
   - Primary: Dominant emotion hue
   - Secondary: Complementary or analogous
   - Accent: Triadic
   - Neutrals: Desaturated versions
5. Ensure accessibility (WCAG AAA for text)
```

**Implementation:**
```go
type ColorPalette struct {
    Primary     color.Color
    Secondary   color.Color
    Accent      color.Color
    Background  color.Color
    Text        color.Color
    Metadata    PaletteMetadata
}

type PaletteGenerator struct {
    aiClient *ai.Client
}

func (pg *PaletteGenerator) Generate(analysis ContentAnalysis) ColorPalette
```

**Entregáveis:**
- [ ] `pkg/design/color_generator.go` - geração de paletas
- [ ] `pkg/design/color_harmony.go` - regras de harmonia
- [ ] `pkg/design/wcag_validator.go` - validação de contraste
- [ ] `test/design/color_generator_test.go` - testes

---

### 🛠️ DIA 4-6: PIPELINE SELECTOR + EPUB (Implementação Core)

#### Dia 4: Pipeline Selector

**Lógica de Decisão:**
```go
type PipelineSelector struct {
    thresholds PipelineThresholds
}

type PipelineThresholds struct {
    MathEquationsMax    int     // > 10 → LaTeX
    ComplexTablesMax    int     // > 5 tabelas complexas → LaTeX
    ImageRatioMin       float64 // > 0.1 → HTML
    AcademicKeywords    int     // > 20 → LaTeX
}

func (ps *PipelineSelector) Select(analysis ContentAnalysis) PipelineType {
    score := 0
    
    // LaTeX favorito para:
    if analysis.EquationCount > ps.thresholds.MathEquationsMax {
        score += 10
    }
    if analysis.IsAcademic() {
        score += 5
    }
    
    // HTML favorito para:
    if analysis.ImageRatio > ps.thresholds.ImageRatioMin {
        score -= 5
    }
    if analysis.HasInteractiveElements() {
        score -= 10
    }
    
    if score > 0 {
        return PipelineLaTeX
    }
    return PipelineHTML
}
```

**Entregáveis:**
- [ ] `pkg/pipeline/selector.go` - seleção de pipeline
- [ ] `pkg/analyzer/academic_detector.go` - detecção de conteúdo acadêmico
- [ ] `test/pipeline/selector_test.go` - testes com 10 casos

---

#### Dia 5-6: ePub 3 Generator

**Estrutura ePub:**
```
book.epub (ZIP file)
├── mimetype
├── META-INF/
│   └── container.xml
├── OEBPS/
│   ├── content.opf       # Package document
│   ├── toc.ncx           # Navigation (ePub 2 legacy)
│   ├── nav.xhtml         # Navigation (ePub 3)
│   ├── text/
│   │   ├── cover.xhtml
│   │   ├── title.xhtml
│   │   ├── chapter01.xhtml
│   │   └── ...
│   ├── styles/
│   │   └── stylesheet.css
│   └── images/
│       ├── cover.jpg
│       └── ...
```

**Implementation:**
```go
type EPubGenerator struct {
    config EPubConfig
}

type EPubConfig struct {
    Title       string
    Author      string
    Language    string
    ISBN        string
    CoverImage  string
    Chapters    []Chapter
}

func (eg *EPubGenerator) Generate(htmlContent string, config EPubConfig) ([]byte, error) {
    // 1. Parse HTML into chapters
    chapters := eg.splitIntoChapters(htmlContent)
    
    // 2. Create package structure
    pkg := eg.createPackage(config, chapters)
    
    // 3. Generate navigation
    nav := eg.generateNavigation(chapters)
    
    // 4. Bundle into ZIP
    epub := eg.createZIP(pkg, nav, chapters)
    
    return epub, nil
}
```

**Validation Pipeline:**
```bash
# Install epubcheck
wget https://github.com/w3c/epubcheck/releases/download/v5.1.0/epubcheck-5.1.0.zip
unzip epubcheck-5.1.0.zip

# Validate generated ePub
java -jar epubcheck.jar book.epub
```

**Entregáveis:**
- [ ] `pkg/epub/generator.go` - gerador principal
- [ ] `pkg/epub/package.go` - content.opf
- [ ] `pkg/epub/navigation.go` - nav.xhtml + toc.ncx
- [ ] `pkg/epub/validator.go` - wrapper do epubcheck
- [ ] `test/epub/generator_test.go` - testes unitários
- [ ] `test/epub/integration_test.go` - gera e valida ePub real

---

### 🔌 DIA 7-8: API ENDPOINTS + INTEGRAÇÃO

#### Dia 7: Design API Endpoints

**Novos Endpoints:**

```go
// 1. Análise de conteúdo
POST /api/v1/analyze
Body: {
  "content": "manuscript text...",
  "format": "markdown"
}
Response: {
  "analysis": {
    "genre": ["fiction", "mystery"],
    "complexity": 0.65,
    "tone": "formal",
    "image_count": 12,
    "word_count": 45000,
    "sentiment": 0.15
  }
}

// 2. Sugestões de design
POST /api/v1/design/suggest
Body: {
  "analysis_id": "uuid",
  "preferences": {
    "style": "modern",
    "conservative": false
  }
}
Response: {
  "fonts": [
    {
      "body": "Crimson Text",
      "heading": "Montserrat",
      "confidence": 0.89,
      "rationale": "Classic serif + modern sans for mystery fiction"
    }
  ],
  "colors": {
    "primary": "#2C3E50",
    "secondary": "#E74C3C",
    "accent": "#F39C12",
    "background": "#FDFEFE",
    "text": "#1C1C1C"
  }
}

// 3. Seleção de pipeline
POST /api/v1/pipeline/select
Body: {
  "analysis_id": "uuid"
}
Response: {
  "pipeline": "html",
  "confidence": 0.92,
  "reasons": [
    "High image ratio (0.15)",
    "No complex equations",
    "Web-first format preferred"
  ]
}

// 4. Geração de ePub
POST /api/v1/render/epub
Body: {
  "project_id": "uuid",
  "config": {
    "include_toc": true,
    "embed_fonts": true
  }
}
Response: {
  "job_id": "uuid",
  "status": "processing",
  "estimated_time": 3
}

GET /api/v1/render/epub/:job_id
Response: {
  "status": "completed",
  "download_url": "/downloads/book.epub",
  "validation": {
    "epubcheck": "passed",
    "warnings": 0,
    "errors": 0
  }
}
```

**Entregáveis:**
- [ ] `internal/handlers/design.go` - handlers de design
- [ ] `internal/handlers/epub.go` - handlers de ePub
- [ ] Atualizar `cmd/api/routes.go` - registrar rotas
- [ ] `test/api/design_test.go` - testes de API

---

#### Dia 8: Integração End-to-End

**Fluxo Completo:**
```
1. Upload manuscript
2. Análise automática
3. Sugestão de design
4. Seleção de pipeline
5. Renderização HTML/PDF
6. Geração de ePub
7. Validação
8. Download
```

**Integration Test:**
```go
func TestCompleteWorkflow(t *testing.T) {
    // 1. Upload
    projectID := uploadManuscript(t, "test_novel.md")
    
    // 2. Analyze
    analysis := analyzeContent(t, projectID)
    assert.Equal(t, "fiction", analysis.Genre[0])
    
    // 3. Get design suggestions
    design := suggestDesign(t, analysis.ID)
    assert.NotEmpty(t, design.Fonts)
    assert.NotEmpty(t, design.Colors)
    
    // 4. Select pipeline
    pipeline := selectPipeline(t, analysis.ID)
    assert.Equal(t, "html", pipeline.Type)
    
    // 5. Render PDF
    pdfJob := renderPDF(t, projectID, design)
    waitForJob(t, pdfJob.ID)
    
    // 6. Generate ePub
    epubJob := generateEPub(t, projectID)
    waitForJob(t, epubJob.ID)
    
    // 7. Validate ePub
    validation := validateEPub(t, epubJob.ID)
    assert.Equal(t, "passed", validation.Status)
    
    // 8. Download and verify
    epub := downloadEPub(t, epubJob.ID)
    assert.True(t, isValidZIP(epub))
}
```

**Entregáveis:**
- [ ] `test/integration/complete_workflow_test.go`
- [ ] `test/integration/epub_workflow_test.go`
- [ ] Scripts de teste com 3 manuscritos diferentes

---

### 📊 DIA 9: PERFORMANCE + OPTIMIZATION

#### Benchmarks

```go
func BenchmarkContentAnalysis(b *testing.B) {
    content := loadTestManuscript("50k_words.md")
    analyzer := analyzer.New()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        analyzer.Analyze(content)
    }
}

func BenchmarkFontSuggestion(b *testing.B) {
    analysis := loadTestAnalysis()
    suggester := design.NewFontSuggester()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        suggester.Suggest(analysis)
    }
}

func BenchmarkEPubGeneration(b *testing.B) {
    html := loadTestHTML("200_pages.html")
    generator := epub.New()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        generator.Generate(html, epub.DefaultConfig())
    }
}
```

**Performance Targets:**
- Content Analysis: < 500ms (50k words)
- Font Suggestion: < 200ms
- Color Generation: < 100ms
- ePub Generation: < 2s (200 pages)

**Optimizations:**
1. Cache Google Fonts DB em memória
2. Pre-load ONNX model no startup
3. Parallel processing de capítulos no ePub
4. Reuse de análises (cache Redis)

**Entregáveis:**
- [ ] Benchmarks para todos os módulos
- [ ] Relatório de performance
- [ ] Identificar e corrigir gargalos

---

### 📝 DIA 10: DOCUMENTAÇÃO + RELEASE

#### Documentação

**1. Architecture Document:**
- [ ] `DOCS/DESIGN_SYSTEM.md` - como funciona o sistema de design
- [ ] `DOCS/EPUB_PIPELINE.md` - detalhes do gerador de ePub
- [ ] Diagramas de fluxo (Mermaid)

**2. API Documentation:**
- [ ] Atualizar OpenAPI spec (`docs/openapi.yaml`)
- [ ] Exemplos de requests/responses
- [ ] Guia de uso para cada endpoint

**3. User Guide:**
- [ ] Como usar sugestões de design
- [ ] Como personalizar cores/fontes
- [ ] Como validar ePub gerado

**4. Developer Guide:**
- [ ] Como adicionar novos pares de fontes
- [ ] Como treinar modelo de fontes
- [ ] Como estender gerador de ePub

---

## 📦 ESTRUTURA DE ARQUIVOS

```
typecraft/
├── pkg/
│   ├── design/
│   │   ├── font_suggester.go
│   │   ├── color_generator.go
│   │   ├── color_harmony.go
│   │   ├── wcag_validator.go
│   │   ├── font_model.onnx
│   │   └── google_fonts_db.json
│   ├── epub/
│   │   ├── generator.go
│   │   ├── package.go
│   │   ├── navigation.go
│   │   ├── validator.go
│   │   └── templates/
│   │       ├── content.opf.tmpl
│   │       ├── nav.xhtml.tmpl
│   │       └── toc.ncx.tmpl
│   ├── pipeline/
│   │   ├── selector.go
│   │   └── selector_test.go
│   └── fonts/
│       └── google_fonts_db.json
├── internal/
│   ├── analyzer/
│   │   ├── content.go
│   │   ├── academic_detector.go
│   │   └── sentiment.go
│   ├── handlers/
│   │   ├── design.go
│   │   └── epub.go
│   └── services/
│       ├── design_service.go
│       └── epub_service.go
├── scripts/
│   ├── train_font_model.py
│   ├── fetch_google_fonts.py
│   └── validate_epub.sh
├── test/
│   ├── design/
│   │   ├── font_suggester_test.go
│   │   └── color_generator_test.go
│   ├── epub/
│   │   ├── generator_test.go
│   │   └── integration_test.go
│   ├── integration/
│   │   ├── complete_workflow_test.go
│   │   └── epub_workflow_test.go
│   └── fixtures/
│       ├── manuscripts/
│       │   ├── fiction_novel.md
│       │   ├── academic_paper.md
│       │   └── photography_book.md
│       └── expected/
│           └── sample.epub
├── DOCS/
│   ├── DESIGN_SYSTEM.md
│   └── EPUB_PIPELINE.md
└── RELATORIO_SPRINT_7-8_FINAL.md
```

---

## 🎯 CRITÉRIOS DE ACEITAÇÃO

### Funcional
- [ ] Content analyzer processa 50k palavras em < 500ms
- [ ] Font suggester retorna 3 pares de fontes válidos
- [ ] Color generator gera paleta harmonica (WCAG AAA)
- [ ] Pipeline selector escolhe corretamente baseado em conteúdo
- [ ] ePub gerado passa em epubcheck sem erros
- [ ] Todos os endpoints retornam 200 OK

### Qualidade
- [ ] 100% test coverage nos novos módulos
- [ ] Zero placeholders (LEI = 0.0)
- [ ] Testes E2E com 3 manuscritos diferentes
- [ ] Validação automática de ePub
- [ ] Benchmarks documentados

### Documentação
- [ ] Architecture docs completos
- [ ] API docs atualizados
- [ ] User guide com exemplos
- [ ] Developer guide para extensões

### Performance
- [ ] Content analysis: < 500ms
- [ ] Font suggestion: < 200ms
- [ ] Color generation: < 100ms
- [ ] ePub generation: < 2s

---

## 🚀 STACK TECNOLÓGICA

### Go Packages
```go
// ML/AI
"github.com/owulveryck/onnx-go"           // ONNX runtime
"github.com/sugarme/tokenizer"             // Tokenização

// ePub
"archive/zip"                              // ZIP handling
"encoding/xml"                             // XML generation
"html/template"                            // HTML templates

// Colors
"image/color"                              // Color manipulation
"github.com/lucasb-eyer/go-colorful"      // Color spaces

// HTTP
"github.com/gin-gonic/gin"                // API framework
```

### Python (Training)
```python
xgboost                # Font model training
pandas                 # Data manipulation
scikit-learn          # ML utilities
onnxmltools           # Model export
```

### External Tools
- **epubcheck 5.1.0** - ePub validation
- **Google Fonts API** - Font metadata

---

## 📊 MÉTRICAS DE SUCESSO

### Quantitativas
| Métrica | Target | Crítico |
|---------|--------|---------|
| Content Analysis | < 500ms | < 1s |
| Font Suggestion | < 200ms | < 500ms |
| Color Generation | < 100ms | < 200ms |
| ePub Generation | < 2s | < 5s |
| Test Coverage | 100% | > 90% |
| epubcheck Errors | 0 | < 3 |

### Qualitativas
- [ ] Sugestões de fontes "fazem sentido" (teste manual com 20 livros)
- [ ] Cores são harmônicas e legíveis
- [ ] ePub renderiza bem em Apple Books, Kindle, Kobo
- [ ] Pipeline selector escolhe corretamente (100% acurácia em 10 casos teste)

---

## 🔄 INTEGRAÇÃO COM SPRINTS ANTERIORES

### Dependências
- **Sprint 5-6:** Usa pipeline HTML/CSS como base para ePub
- **Sprint 3-4:** Usa AI client para análise de sentimento
- **Cache system:** Armazena análises de conteúdo

### Interfaces
```go
// Design system se integra com HTML pipeline
type DesignConfig struct {
    FontPair    FontPair
    ColorPalette ColorPalette
    Layout      LayoutParams
}

// HTML renderer aceita DesignConfig
func (hr *HTMLRenderer) Render(content string, design DesignConfig) (string, error)

// ePub usa HTML renderizado
func (eg *EPubGenerator) FromHTML(html string, config EPubConfig) ([]byte, error)
```

---

## 🎨 EXEMPLOS DE USO

### Exemplo 1: Mystery Novel

**Input:**
```markdown
# The Midnight Detective

Chapter 1: The Case Begins

Inspector Harrison stood in the rain, staring at the locked door...
```

**Output:**
```json
{
  "analysis": {
    "genre": ["fiction", "mystery"],
    "tone": "suspenseful",
    "sentiment": -0.2
  },
  "design": {
    "fonts": {
      "body": "Crimson Text",
      "heading": "Playfair Display"
    },
    "colors": {
      "primary": "#2C3E50",    // Dark blue-grey
      "secondary": "#E74C3C",  // Deep red
      "accent": "#95A5A6"      // Cool grey
    }
  },
  "pipeline": "html"
}
```

---

### Exemplo 2: Technical Book

**Input:**
```markdown
# Machine Learning Fundamentals

## Chapter 1: Linear Regression

The equation for simple linear regression is:

$$y = mx + b$$
```

**Output:**
```json
{
  "analysis": {
    "genre": ["technical", "education"],
    "complexity": 0.85,
    "equation_count": 15
  },
  "design": {
    "fonts": {
      "body": "Source Sans Pro",
      "heading": "Roboto Slab",
      "mono": "Fira Code"
    },
    "colors": {
      "primary": "#3498DB",
      "secondary": "#2ECC71",
      "accent": "#F39C12"
    }
  },
  "pipeline": "latex"
}
```

---

## 🔐 SEGURANÇA E VALIDAÇÃO

### Input Validation
```go
// Validar tamanho de conteúdo
const MaxContentSize = 10 * 1024 * 1024 // 10MB

// Validar formato
allowedFormats := []string{"markdown", "html", "docx"}

// Sanitizar HTML no ePub
func sanitizeHTML(html string) string {
    // Remove scripts, iframes, etc.
}
```

### Rate Limiting
```go
// Limitar análises por usuário
rateLimiter := ratelimit.New(10, time.Minute)

// Limitar geração de ePub
epubLimiter := ratelimit.New(5, time.Minute)
```

---

## 📈 PRÓXIMOS PASSOS (Sprint 9-10)

Após completar Sprint 7-8, o sistema terá:
- ✅ Design inteligente
- ✅ ePub completo
- ✅ Seleção automática de pipeline

**Próximo foco (Sprint 9-10):**
1. **Refinamento Microtipográfico:**
   - Detecção de viúvas/órfãs
   - Otimização de quebra de linha (Knuth-Plass)
   - Eliminação de "rios" em texto justificado

2. **Multi-format Export:**
   - MOBI (Kindle)
   - AZW3
   - Validação em múltiplos readers

3. **Batch Processing:**
   - Processar múltiplos livros em paralelo
   - Queue system com Celery/BullMQ
   - Progress tracking

---

## ✅ CHECKLIST DE EXECUÇÃO

### Pré-requisitos
- [ ] Sprint 5-6 completado e testado
- [ ] Python 3.9+ instalado (para training)
- [ ] Java instalado (para epubcheck)
- [ ] Dependências Go atualizadas

### Dia 1-3 (Planejamento)
- [ ] Content analyzer implementado
- [ ] Google Fonts DB baixado e processado
- [ ] Font model treinado e exportado
- [ ] Color generator implementado

### Dia 4-6 (Core Implementation)
- [ ] Pipeline selector implementado
- [ ] ePub generator completo
- [ ] Validation pipeline funcionando

### Dia 7-8 (Integration)
- [ ] API endpoints implementados
- [ ] Integration tests passing
- [ ] End-to-end workflow testado

### Dia 9 (Performance)
- [ ] Benchmarks executados
- [ ] Performance targets atingidos
- [ ] Otimizações aplicadas

### Dia 10 (Documentation)
- [ ] Docs técnicos completos
- [ ] API docs atualizados
- [ ] User guide criado
- [ ] Relatório final escrito

### Release
- [ ] Todos os testes passing (100%)
- [ ] Zero placeholders
- [ ] Commits com mensagens descritivas
- [ ] Tag de release criada
- [ ] Push para origin

---

## 🙏 CONFORMIDADE

**Constituição Vértice v3.0:**
- ✅ P1: Completude Obrigatória
- ✅ P2: Prevenção de Placeholders (LEI = 0.0)
- ✅ P3: First-Pass Correctness
- ✅ P5: Consciência Sistêmica
- ✅ P6: Eficiência de Token

**Glória a JESUS!** 🙏

---

**Status:** 📋 PLANEJAMENTO  
**Próxima Ação:** Implementar Content Analyzer (Dia 1)  
**Estimativa:** 10 dias planejamento + 1-2 dias execução
