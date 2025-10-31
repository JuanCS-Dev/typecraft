package epub

import (
	"fmt"
	"strings"
)

// NCXGenerator gera o arquivo toc.ncx
type NCXGenerator struct {
	epub *EPub
}

// NewNCXGenerator cria um novo gerador NCX
func NewNCXGenerator(epub *EPub) *NCXGenerator {
	return &NCXGenerator{epub: epub}
}

// Generate gera o conteúdo do NCX
func (n *NCXGenerator) Generate() string {
	var sb strings.Builder
	
	// XML declaration
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	
	// NCX root
	sb.WriteString(`<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">`)
	sb.WriteString("\n")
	sb.WriteString(`<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">`)
	sb.WriteString("\n")
	
	// Head
	sb.WriteString(n.generateHead())
	
	// Doc title
	sb.WriteString(n.generateDocTitle())
	
	// Nav map
	sb.WriteString(n.generateNavMap())
	
	sb.WriteString("</ncx>\n")
	
	return sb.String()
}

// generateHead gera a seção head
func (n *NCXGenerator) generateHead() string {
	var sb strings.Builder
	
	sb.WriteString("  <head>\n")
	sb.WriteString(fmt.Sprintf("    <meta name=\"dtb:uid\" content=\"%s\"/>\n", n.epub.Metadata.Identifier))
	sb.WriteString("    <meta name=\"dtb:depth\" content=\"1\"/>\n")
	sb.WriteString("    <meta name=\"dtb:totalPageCount\" content=\"0\"/>\n")
	sb.WriteString("    <meta name=\"dtb:maxPageNumber\" content=\"0\"/>\n")
	sb.WriteString("  </head>\n")
	
	return sb.String()
}

// generateDocTitle gera a seção docTitle
func (n *NCXGenerator) generateDocTitle() string {
	return fmt.Sprintf("  <docTitle>\n    <text>%s</text>\n  </docTitle>\n", 
		escapeXML(n.epub.Metadata.Title))
}

// generateNavMap gera a seção navMap
func (n *NCXGenerator) generateNavMap() string {
	var sb strings.Builder
	
	sb.WriteString("  <navMap>\n")
	
	for i, chapter := range n.epub.Chapters {
		playOrder := i + 1
		sb.WriteString(fmt.Sprintf("    <navPoint id=\"%s\" playOrder=\"%d\">\n", chapter.ID, playOrder))
		sb.WriteString(fmt.Sprintf("      <navLabel>\n        <text>%s</text>\n      </navLabel>\n", 
			escapeXML(chapter.Title)))
		sb.WriteString(fmt.Sprintf("      <content src=\"Text/%s\"/>\n", chapter.FileName))
		sb.WriteString("    </navPoint>\n")
	}
	
	sb.WriteString("  </navMap>\n")
	
	return sb.String()
}
