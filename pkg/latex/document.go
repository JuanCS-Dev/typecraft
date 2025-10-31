package latex

import (
	"fmt"
	"strings"
)

// DocumentClass representa a classe do documento
type DocumentClass string

const (
	ClassArticle DocumentClass = "article"
	ClassBook    DocumentClass = "book"
	ClassReport  DocumentClass = "report"
	ClassBeamer  DocumentClass = "beamer"
)

// Document representa um documento LaTeX
type Document struct {
	Class    DocumentClass
	Options  []string
	Packages []Package
	Metadata DocumentMetadata
	Content  []string
}

// Package representa um pacote LaTeX
type Package struct {
	Name    string
	Options []string
}

// DocumentMetadata representa metadados do documento
type DocumentMetadata struct {
	Title   string
	Author  string
	Date    string
	Subject string
}

// NewDocument cria um novo documento
func NewDocument(class DocumentClass) *Document {
	return &Document{
		Class:    class,
		Options:  make([]string, 0),
		Packages: make([]Package, 0),
		Content:  make([]string, 0),
	}
}

// AddOption adiciona uma opção à classe do documento
func (d *Document) AddOption(option string) *Document {
	d.Options = append(d.Options, option)
	return d
}

// AddPackage adiciona um pacote
func (d *Document) AddPackage(name string, options ...string) *Document {
	d.Packages = append(d.Packages, Package{
		Name:    name,
		Options: options,
	})
	return d
}

// SetMetadata define metadados
func (d *Document) SetMetadata(metadata DocumentMetadata) *Document {
	d.Metadata = metadata
	return d
}

// SetTitle define o título
func (d *Document) SetTitle(title string) *Document {
	d.Metadata.Title = title
	return d
}

// SetAuthor define o autor
func (d *Document) SetAuthor(author string) *Document {
	d.Metadata.Author = author
	return d
}

// SetDate define a data
func (d *Document) SetDate(date string) *Document {
	d.Metadata.Date = date
	return d
}

// SetSubject define o assunto
func (d *Document) SetSubject(subject string) *Document {
	d.Metadata.Subject = subject
	return d
}

// AddContent adiciona conteúdo
func (d *Document) AddContent(content string) *Document {
	d.Content = append(d.Content, content)
	return d
}

// AddSection adiciona uma seção
func (d *Document) AddSection(title string, content string) *Document {
	d.Content = append(d.Content, fmt.Sprintf("\\section{%s}\n%s", title, content))
	return d
}

// AddSubsection adiciona uma subseção
func (d *Document) AddSubsection(title string, content string) *Document {
	d.Content = append(d.Content, fmt.Sprintf("\\subsection{%s}\n%s", title, content))
	return d
}

// AddChapter adiciona um capítulo (apenas para book/report)
func (d *Document) AddChapter(title string, content string) *Document {
	if d.Class == ClassBook || d.Class == ClassReport {
		d.Content = append(d.Content, fmt.Sprintf("\\chapter{%s}\n%s", title, content))
	}
	return d
}

// Generate gera o código LaTeX
func (d *Document) Generate() string {
	var sb strings.Builder
	
	// Document class
	if len(d.Options) > 0 {
		sb.WriteString(fmt.Sprintf("\\documentclass[%s]{%s}\n", 
			strings.Join(d.Options, ","), d.Class))
	} else {
		sb.WriteString(fmt.Sprintf("\\documentclass{%s}\n", d.Class))
	}
	
	sb.WriteString("\n")
	
	// Packages
	for _, pkg := range d.Packages {
		if len(pkg.Options) > 0 {
			sb.WriteString(fmt.Sprintf("\\usepackage[%s]{%s}\n",
				strings.Join(pkg.Options, ","), pkg.Name))
		} else {
			sb.WriteString(fmt.Sprintf("\\usepackage{%s}\n", pkg.Name))
		}
	}
	
	if len(d.Packages) > 0 {
		sb.WriteString("\n")
	}
	
	// Metadata
	if d.Metadata.Title != "" {
		sb.WriteString(fmt.Sprintf("\\title{%s}\n", d.Metadata.Title))
	}
	if d.Metadata.Author != "" {
		sb.WriteString(fmt.Sprintf("\\author{%s}\n", d.Metadata.Author))
	}
	if d.Metadata.Date != "" {
		sb.WriteString(fmt.Sprintf("\\date{%s}\n", d.Metadata.Date))
	}
	
	if d.Metadata.Title != "" || d.Metadata.Author != "" || d.Metadata.Date != "" {
		sb.WriteString("\n")
	}
	
	// Begin document
	sb.WriteString("\\begin{document}\n\n")
	
	// Maketitle if metadata exists
	if d.Metadata.Title != "" || d.Metadata.Author != "" {
		sb.WriteString("\\maketitle\n\n")
	}
	
	// Content
	for _, content := range d.Content {
		sb.WriteString(content)
		sb.WriteString("\n\n")
	}
	
	// End document
	sb.WriteString("\\end{document}\n")
	
	return sb.String()
}

// DocumentBuilder facilita construção de documentos
type DocumentBuilder struct {
	document *Document
}

// NewDocumentBuilder cria um novo builder
func NewDocumentBuilder(class DocumentClass) *DocumentBuilder {
	return &DocumentBuilder{
		document: NewDocument(class),
	}
}

// WithOptions adiciona opções
func (db *DocumentBuilder) WithOptions(options ...string) *DocumentBuilder {
	db.document.Options = append(db.document.Options, options...)
	return db
}

// WithPackage adiciona um pacote
func (db *DocumentBuilder) WithPackage(name string, options ...string) *DocumentBuilder {
	db.document.AddPackage(name, options...)
	return db
}

// WithMetadata adiciona metadados
func (db *DocumentBuilder) WithMetadata(metadata DocumentMetadata) *DocumentBuilder {
	db.document.SetMetadata(metadata)
	return db
}

// WithContent adiciona conteúdo
func (db *DocumentBuilder) WithContent(content string) *DocumentBuilder {
	db.document.AddContent(content)
	return db
}

// WithSection adiciona seção
func (db *DocumentBuilder) WithSection(title, content string) *DocumentBuilder {
	db.document.AddSection(title, content)
	return db
}

// WithChapter adiciona capítulo
func (db *DocumentBuilder) WithChapter(title, content string) *DocumentBuilder {
	db.document.AddChapter(title, content)
	return db
}

// Build constrói o documento
func (db *DocumentBuilder) Build() *Document {
	return db.document
}

// StandardPackages retorna pacotes comuns
func StandardPackages() []Package {
	return []Package{
		{Name: "inputenc", Options: []string{"utf8"}},
		{Name: "fontenc", Options: []string{"T1"}},
		{Name: "babel", Options: []string{"english"}},
		{Name: "geometry", Options: []string{"margin=1in"}},
		{Name: "graphicx"},
		{Name: "hyperref"},
		{Name: "amsmath"},
		{Name: "amsfonts"},
		{Name: "amssymb"},
	}
}

// AcademicPackages retorna pacotes para documentos acadêmicos
func AcademicPackages() []Package {
	packages := StandardPackages()
	packages = append(packages,
		Package{Name: "natbib"},
		Package{Name: "algorithm"},
		Package{Name: "algorithmic"},
		Package{Name: "listings"},
		Package{Name: "booktabs"},
		Package{Name: "caption"},
	)
	return packages
}

// BookPackages retorna pacotes para livros
func BookPackages() []Package {
	packages := StandardPackages()
	packages = append(packages,
		Package{Name: "fancyhdr"},
		Package{Name: "titlesec"},
		Package{Name: "tocloft"},
	)
	return packages
}
