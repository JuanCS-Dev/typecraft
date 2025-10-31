package epub

import (
	"fmt"
	"strings"
)

// NavGenerator gera o arquivo nav.xhtml (EPUB 3)
type NavGenerator struct {
	epub *EPub
}

// NewNavGenerator cria um novo gerador Nav
func NewNavGenerator(epub *EPub) *NavGenerator {
	return &NavGenerator{epub: epub}
}

// Generate gera o conteúdo do nav.xhtml
func (n *NavGenerator) Generate() string {
	var sb strings.Builder
	
	// XML declaration
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	
	// HTML root
	sb.WriteString(`<!DOCTYPE html>`)
	sb.WriteString("\n")
	sb.WriteString(`<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">`)
	sb.WriteString("\n")
	
	// Head
	sb.WriteString(n.generateHead())
	
	// Body
	sb.WriteString(n.generateBody())
	
	sb.WriteString("</html>\n")
	
	return sb.String()
}

// generateHead gera a seção head
func (n *NavGenerator) generateHead() string {
	var sb strings.Builder
	
	sb.WriteString("  <head>\n")
	sb.WriteString("    <meta charset=\"UTF-8\"/>\n")
	sb.WriteString(fmt.Sprintf("    <title>%s</title>\n", escapeXML(n.epub.Metadata.Title)))
	sb.WriteString("    <link rel=\"stylesheet\" type=\"text/css\" href=\"Styles/style.css\"/>\n")
	sb.WriteString("  </head>\n")
	
	return sb.String()
}

// generateBody gera a seção body
func (n *NavGenerator) generateBody() string {
	var sb strings.Builder
	
	sb.WriteString("  <body>\n")
	
	// Table of contents
	sb.WriteString("    <nav epub:type=\"toc\" id=\"toc\">\n")
	sb.WriteString("      <h1>Table of Contents</h1>\n")
	sb.WriteString("      <ol>\n")
	
	for _, chapter := range n.epub.Chapters {
		sb.WriteString(fmt.Sprintf("        <li><a href=\"Text/%s\">%s</a></li>\n",
			chapter.FileName, escapeXML(chapter.Title)))
	}
	
	sb.WriteString("      </ol>\n")
	sb.WriteString("    </nav>\n")
	
	// Landmarks (optional)
	sb.WriteString("    <nav epub:type=\"landmarks\" id=\"landmarks\" hidden=\"\">\n")
	sb.WriteString("      <h2>Landmarks</h2>\n")
	sb.WriteString("      <ol>\n")
	
	if len(n.epub.Chapters) > 0 {
		sb.WriteString(fmt.Sprintf("        <li><a epub:type=\"bodymatter\" href=\"Text/%s\">Start of Content</a></li>\n",
			n.epub.Chapters[0].FileName))
	}
	
	sb.WriteString("      </ol>\n")
	sb.WriteString("    </nav>\n")
	
	sb.WriteString("  </body>\n")
	
	return sb.String()
}
