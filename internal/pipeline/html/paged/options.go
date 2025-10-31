package paged

// PageOptions opções de paginação e formatação
// Conformidade: CONSTITUIÇÃO VÉRTICE v3.0
type PageOptions struct {
	// Format tamanho da página (A4, A5, Letter, etc)
	Format string
	
	// Landscape orientação paisagem
	Landscape bool
	
	// Timeout timeout em milisegundos para renderização
	Timeout int
	
	// CustomCSS CSS customizado para injetar
	CustomCSS string
	
	// Margins margens da página
	Margins PageMargins
	
	// RunningHeaders cabeçalhos dinâmicos
	RunningHeaders bool
	
	// PageNumbers numeração de páginas
	PageNumbers PageNumberOptions
}

// PageMargins configuração de margens
type PageMargins struct {
	Top    string // e.g., "2cm"
	Right  string
	Bottom string
	Left   string
	Inside string // Para páginas duplas
	Outside string
}

// PageNumberOptions configuração de numeração
type PageNumberOptions struct {
	Enabled   bool
	Format    string // "decimal", "roman", "alpha"
	Position  string // "top-right", "bottom-center", etc
	StartFrom int
}

// DefaultPageOptions retorna opções padrão
func DefaultPageOptions() PageOptions {
	return PageOptions{
		Format:    "A4",
		Landscape: false,
		Timeout:   30000, // 30 segundos
		Margins: PageMargins{
			Top:    "2.5cm",
			Right:  "2cm",
			Bottom: "2.5cm",
			Left:   "2cm",
		},
		RunningHeaders: true,
		PageNumbers: PageNumberOptions{
			Enabled:   true,
			Format:    "decimal",
			Position:  "bottom-center",
			StartFrom: 1,
		},
	}
}

// GeneratePagedCSS gera CSS para Paged Media
func (opts PageOptions) GeneratePagedCSS() string {
	css := "@page {\n"
	
	// Margens
	if opts.Margins.Top != "" {
		css += "  margin-top: " + opts.Margins.Top + ";\n"
	}
	if opts.Margins.Right != "" {
		css += "  margin-right: " + opts.Margins.Right + ";\n"
	}
	if opts.Margins.Bottom != "" {
		css += "  margin-bottom: " + opts.Margins.Bottom + ";\n"
	}
	if opts.Margins.Left != "" {
		css += "  margin-left: " + opts.Margins.Left + ";\n"
	}
	
	// Tamanho da página
	if opts.Format != "" {
		css += "  size: " + opts.Format
		if opts.Landscape {
			css += " landscape"
		}
		css += ";\n"
	}
	
	css += "}\n\n"
	
	// Numeração de páginas
	if opts.PageNumbers.Enabled {
		css += opts.generatePageNumberCSS()
	}
	
	// Running headers
	if opts.RunningHeaders {
		css += opts.generateRunningHeadersCSS()
	}
	
	return css
}

func (opts PageOptions) generatePageNumberCSS() string {
	position := opts.PageNumbers.Position
	if position == "" {
		position = "bottom-center"
	}
	
	format := opts.PageNumbers.Format
	if format == "" {
		format = "decimal"
	}
	
	css := "/* Page numbers */\n"
	css += "@page {\n"
	
	switch position {
	case "bottom-center":
		css += "  @bottom-center {\n"
		css += "    content: counter(page, " + format + ");\n"
		css += "  }\n"
	case "top-right":
		css += "  @top-right {\n"
		css += "    content: counter(page, " + format + ");\n"
		css += "  }\n"
	case "bottom-right":
		css += "  @bottom-right {\n"
		css += "    content: counter(page, " + format + ");\n"
		css += "  }\n"
	}
	
	css += "}\n\n"
	return css
}

func (opts PageOptions) generateRunningHeadersCSS() string {
	css := "/* Running headers */\n"
	css += "@page :left {\n"
	css += "  @top-left {\n"
	css += "    content: string(chapter);\n"
	css += "    font-style: italic;\n"
	css += "  }\n"
	css += "}\n\n"
	css += "@page :right {\n"
	css += "  @top-right {\n"
	css += "    content: string(section);\n"
	css += "    font-style: italic;\n"
	css += "  }\n"
	css += "}\n\n"
	css += "h1 { string-set: chapter content(); }\n"
	css += "h2 { string-set: section content(); }\n\n"
	return css
}
