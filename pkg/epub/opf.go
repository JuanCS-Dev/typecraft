package epub

import (
	"fmt"
	"strings"
	"time"
)

// OPFGenerator gera o arquivo content.opf
type OPFGenerator struct {
	epub *EPub
}

// NewOPFGenerator cria um novo gerador OPF
func NewOPFGenerator(epub *EPub) *OPFGenerator {
	return &OPFGenerator{epub: epub}
}

// Generate gera o conteúdo do OPF
func (o *OPFGenerator) Generate() string {
	var sb strings.Builder
	
	// XML declaration
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	
	// Package element
	if o.epub.Version == EPub3 {
		sb.WriteString(`<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="BookID">`)
	} else {
		sb.WriteString(`<package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="BookID">`)
	}
	sb.WriteString("\n")
	
	// Metadata
	sb.WriteString(o.generateMetadata())
	
	// Manifest
	sb.WriteString(o.generateManifest())
	
	// Spine
	sb.WriteString(o.generateSpine())
	
	// Guide (optional, EPUB 2 legacy)
	if o.epub.Version == EPub2 {
		sb.WriteString(o.generateGuide())
	}
	
	sb.WriteString("</package>\n")
	
	return sb.String()
}

// generateMetadata gera a seção de metadata
func (o *OPFGenerator) generateMetadata() string {
	var sb strings.Builder
	
	sb.WriteString("  <metadata xmlns:dc=\"http://purl.org/dc/elements/1.1/\" ")
	sb.WriteString("xmlns:opf=\"http://www.idpf.org/2007/opf\">\n")
	
	m := o.epub.Metadata
	
	// Required metadata
	sb.WriteString(fmt.Sprintf("    <dc:title>%s</dc:title>\n", escapeXML(m.Title)))
	sb.WriteString(fmt.Sprintf("    <dc:creator opf:role=\"aut\">%s</dc:creator>\n", escapeXML(m.Author)))
	sb.WriteString(fmt.Sprintf("    <dc:language>%s</dc:language>\n", m.Language))
	sb.WriteString(fmt.Sprintf("    <dc:identifier id=\"BookID\">%s</dc:identifier>\n", m.Identifier))
	
	// Optional metadata
	if m.Publisher != "" {
		sb.WriteString(fmt.Sprintf("    <dc:publisher>%s</dc:publisher>\n", escapeXML(m.Publisher)))
	}
	
	if m.Description != "" {
		sb.WriteString(fmt.Sprintf("    <dc:description>%s</dc:description>\n", escapeXML(m.Description)))
	}
	
	for _, subject := range m.Subject {
		sb.WriteString(fmt.Sprintf("    <dc:subject>%s</dc:subject>\n", escapeXML(subject)))
	}
	
	if !m.Date.IsZero() {
		sb.WriteString(fmt.Sprintf("    <dc:date>%s</dc:date>\n", m.Date.Format("2006-01-02")))
	}
	
	if m.Rights != "" {
		sb.WriteString(fmt.Sprintf("    <dc:rights>%s</dc:rights>\n", escapeXML(m.Rights)))
	}
	
	// Cover metadata (EPUB 3)
	if o.epub.Version == EPub3 && m.CoverImage != "" {
		sb.WriteString("    <meta name=\"cover\" content=\"cover-image\"/>\n")
	}
	
	// Modified date (EPUB 3)
	if o.epub.Version == EPub3 {
		sb.WriteString(fmt.Sprintf("    <meta property=\"dcterms:modified\">%s</meta>\n", 
			time.Now().UTC().Format("2006-01-02T15:04:05Z")))
	}
	
	sb.WriteString("  </metadata>\n")
	
	return sb.String()
}

// generateManifest gera a seção manifest
func (o *OPFGenerator) generateManifest() string {
	var sb strings.Builder
	
	sb.WriteString("  <manifest>\n")
	
	// NCX (EPUB 2) or Nav (EPUB 3)
	if o.epub.Version == EPub3 {
		sb.WriteString("    <item id=\"nav\" href=\"nav.xhtml\" media-type=\"application/xhtml+xml\" properties=\"nav\"/>\n")
	}
	sb.WriteString("    <item id=\"ncx\" href=\"toc.ncx\" media-type=\"application/x-dtbncx+xml\"/>\n")
	
	// CSS
	sb.WriteString("    <item id=\"css\" href=\"Styles/style.css\" media-type=\"text/css\"/>\n")
	
	// Chapters
	for _, chapter := range o.epub.Chapters {
		sb.WriteString(fmt.Sprintf("    <item id=\"%s\" href=\"Text/%s\" media-type=\"application/xhtml+xml\"/>\n",
			chapter.ID, chapter.FileName))
	}
	
	// Fonts
	for i, font := range o.epub.Fonts {
		mimeType := font.MimeType
		if mimeType == "" {
			mimeType = "application/vnd.ms-opentype"
		}
		sb.WriteString(fmt.Sprintf("    <item id=\"font%d\" href=\"Fonts/%s\" media-type=\"%s\"/>\n",
			i+1, font.Name, mimeType))
	}
	
	// Images
	for _, img := range o.epub.Images {
		mimeType := img.MimeType
		if mimeType == "" {
			mimeType = "image/jpeg"
		}
		properties := ""
		if o.epub.Metadata.CoverImage != "" && img.ID == "cover-image" {
			properties = " properties=\"cover-image\""
		}
		sb.WriteString(fmt.Sprintf("    <item id=\"%s\" href=\"Images/%s\" media-type=\"%s\"%s/>\n",
			img.ID, img.Path, mimeType, properties))
	}
	
	sb.WriteString("  </manifest>\n")
	
	return sb.String()
}

// generateSpine gera a seção spine
func (o *OPFGenerator) generateSpine() string {
	var sb strings.Builder
	
	sb.WriteString("  <spine toc=\"ncx\">\n")
	
	for _, chapter := range o.epub.Chapters {
		sb.WriteString(fmt.Sprintf("    <itemref idref=\"%s\"/>\n", chapter.ID))
	}
	
	sb.WriteString("  </spine>\n")
	
	return sb.String()
}

// generateGuide gera a seção guide (EPUB 2)
func (o *OPFGenerator) generateGuide() string {
	var sb strings.Builder
	
	sb.WriteString("  <guide>\n")
	
	if len(o.epub.Chapters) > 0 {
		firstChapter := o.epub.Chapters[0]
		sb.WriteString(fmt.Sprintf("    <reference type=\"text\" title=\"Start\" href=\"Text/%s\"/>\n",
			firstChapter.FileName))
	}
	
	sb.WriteString("  </guide>\n")
	
	return sb.String()
}

// escapeXML escapa caracteres especiais XML
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
