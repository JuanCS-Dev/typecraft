package epub

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// EPubVersion representa a versão do ePub
type EPubVersion string

const (
	EPub2 EPubVersion = "2.0"
	EPub3 EPubVersion = "3.0"
)

// Metadata representa os metadados do livro
type Metadata struct {
	Title       string
	Author      string
	Language    string
	Identifier  string // UUID or ISBN
	Publisher   string
	Description string
	Subject     []string
	Date        time.Time
	Rights      string
	CoverImage  string // path to cover image
}

// Chapter representa um capítulo do livro
type Chapter struct {
	ID       string
	Title    string
	Content  string // HTML content
	FileName string // e.g., "chapter1.xhtml"
}

// EPub representa um livro ePub
type EPub struct {
	Version  EPubVersion
	Metadata Metadata
	Chapters []Chapter
	CSS      string // Global CSS
	Fonts    []FontFile
	Images   []ImageFile
	
	// Internal
	tempDir string
}

// FontFile representa um arquivo de fonte
type FontFile struct {
	Name     string
	Path     string
	MimeType string
}

// ImageFile representa um arquivo de imagem
type ImageFile struct {
	ID       string
	Path     string
	MimeType string
}

// NewEPub cria um novo ePub
func NewEPub(version EPubVersion) *EPub {
	return &EPub{
		Version:  version,
		Chapters: make([]Chapter, 0),
		Fonts:    make([]FontFile, 0),
		Images:   make([]ImageFile, 0),
	}
}

// AddChapter adiciona um capítulo
func (e *EPub) AddChapter(chapter Chapter) {
	if chapter.ID == "" {
		chapter.ID = fmt.Sprintf("chapter%d", len(e.Chapters)+1)
	}
	if chapter.FileName == "" {
		chapter.FileName = fmt.Sprintf("chapter%d.xhtml", len(e.Chapters)+1)
	}
	e.Chapters = append(e.Chapters, chapter)
}

// AddFont adiciona uma fonte
func (e *EPub) AddFont(font FontFile) {
	e.Fonts = append(e.Fonts, font)
}

// AddImage adiciona uma imagem
func (e *EPub) AddImage(img ImageFile) {
	e.Images = append(e.Images, img)
}

// Write gera o arquivo ePub
func (e *EPub) Write(outputPath string) error {
	// 1. Criar diretório temporário
	tempDir, err := os.MkdirTemp("", "epub-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	e.tempDir = tempDir
	defer os.RemoveAll(tempDir)
	
	// 2. Criar estrutura de diretórios
	if err := e.createStructure(); err != nil {
		return err
	}
	
	// 3. Escrever mimetype
	if err := e.writeMimetype(); err != nil {
		return err
	}
	
	// 4. Escrever container.xml
	if err := e.writeContainer(); err != nil {
		return err
	}
	
	// 5. Escrever content.opf
	if err := e.writeContentOPF(); err != nil {
		return err
	}
	
	// 6. Escrever toc.ncx
	if err := e.writeTocNCX(); err != nil {
		return err
	}
	
	// 7. Escrever nav.xhtml (EPUB 3)
	if e.Version == EPub3 {
		if err := e.writeNavXHTML(); err != nil {
			return err
		}
	}
	
	// 8. Escrever capítulos
	if err := e.writeChapters(); err != nil {
		return err
	}
	
	// 9. Escrever CSS
	if err := e.writeCSS(); err != nil {
		return err
	}
	
	// 10. Copiar fontes
	if err := e.copyFonts(); err != nil {
		return err
	}
	
	// 11. Copiar imagens
	if err := e.copyImages(); err != nil {
		return err
	}
	
	// 12. Criar arquivo ZIP
	if err := e.createZip(outputPath); err != nil {
		return err
	}
	
	return nil
}

// createStructure cria a estrutura de diretórios
func (e *EPub) createStructure() error {
	dirs := []string{
		filepath.Join(e.tempDir, "META-INF"),
		filepath.Join(e.tempDir, "OEBPS"),
		filepath.Join(e.tempDir, "OEBPS", "Text"),
		filepath.Join(e.tempDir, "OEBPS", "Styles"),
		filepath.Join(e.tempDir, "OEBPS", "Fonts"),
		filepath.Join(e.tempDir, "OEBPS", "Images"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create dir %s: %w", dir, err)
		}
	}
	
	return nil
}

// writeMimetype escreve o arquivo mimetype
func (e *EPub) writeMimetype() error {
	path := filepath.Join(e.tempDir, "mimetype")
	return os.WriteFile(path, []byte("application/epub+zip"), 0644)
}

// writeContainer escreve META-INF/container.xml
func (e *EPub) writeContainer() error {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`
	
	path := filepath.Join(e.tempDir, "META-INF", "container.xml")
	return os.WriteFile(path, []byte(content), 0644)
}

// writeContentOPF escreve OEBPS/content.opf
func (e *EPub) writeContentOPF() error {
	opf := NewOPFGenerator(e)
	content := opf.Generate()
	
	path := filepath.Join(e.tempDir, "OEBPS", "content.opf")
	return os.WriteFile(path, []byte(content), 0644)
}

// writeTocNCX escreve OEBPS/toc.ncx
func (e *EPub) writeTocNCX() error {
	ncx := NewNCXGenerator(e)
	content := ncx.Generate()
	
	path := filepath.Join(e.tempDir, "OEBPS", "toc.ncx")
	return os.WriteFile(path, []byte(content), 0644)
}

// writeNavXHTML escreve OEBPS/nav.xhtml (EPUB 3)
func (e *EPub) writeNavXHTML() error {
	nav := NewNavGenerator(e)
	content := nav.Generate()
	
	path := filepath.Join(e.tempDir, "OEBPS", "nav.xhtml")
	return os.WriteFile(path, []byte(content), 0644)
}

// writeChapters escreve os capítulos
func (e *EPub) writeChapters() error {
	for _, chapter := range e.Chapters {
		content := e.wrapChapterHTML(chapter)
		path := filepath.Join(e.tempDir, "OEBPS", "Text", chapter.FileName)
		
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write chapter %s: %w", chapter.FileName, err)
		}
	}
	return nil
}

// wrapChapterHTML envolve o conteúdo do capítulo em HTML válido
func (e *EPub) wrapChapterHTML(chapter Chapter) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
  <meta charset="UTF-8"/>
  <title>%s</title>
  <link rel="stylesheet" type="text/css" href="../Styles/style.css"/>
</head>
<body>
  <section epub:type="chapter">
    <h1>%s</h1>
    %s
  </section>
</body>
</html>`, chapter.Title, chapter.Title, chapter.Content)
}

// writeCSS escreve o CSS
func (e *EPub) writeCSS() error {
	css := e.CSS
	if css == "" {
		css = defaultCSS()
	}
	
	path := filepath.Join(e.tempDir, "OEBPS", "Styles", "style.css")
	return os.WriteFile(path, []byte(css), 0644)
}

// copyFonts copia as fontes
func (e *EPub) copyFonts() error {
	for _, font := range e.Fonts {
		if _, err := os.Stat(font.Path); os.IsNotExist(err) {
			continue // Skip if font doesn't exist
		}
		
		dest := filepath.Join(e.tempDir, "OEBPS", "Fonts", font.Name)
		if err := copyFile(font.Path, dest); err != nil {
			return fmt.Errorf("failed to copy font %s: %w", font.Name, err)
		}
	}
	return nil
}

// copyImages copia as imagens
func (e *EPub) copyImages() error {
	for _, img := range e.Images {
		if _, err := os.Stat(img.Path); os.IsNotExist(err) {
			continue // Skip if image doesn't exist
		}
		
		dest := filepath.Join(e.tempDir, "OEBPS", "Images", filepath.Base(img.Path))
		if err := copyFile(img.Path, dest); err != nil {
			return fmt.Errorf("failed to copy image %s: %w", img.ID, err)
		}
	}
	return nil
}

// createZip cria o arquivo ZIP
func (e *EPub) createZip(outputPath string) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create zip: %w", err)
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// mimetype MUST be first and uncompressed
	if err := e.addMimetypeToZip(zipWriter); err != nil {
		return err
	}
	
	// Add all other files
	return filepath.Walk(e.tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || filepath.Base(path) == "mimetype" {
			return nil
		}
		
		relPath, err := filepath.Rel(e.tempDir, path)
		if err != nil {
			return err
		}
		
		return e.addFileToZip(zipWriter, path, relPath)
	})
}

// addMimetypeToZip adiciona mimetype sem compressão
func (e *EPub) addMimetypeToZip(zw *zip.Writer) error {
	header := &zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store, // No compression
	}
	
	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	
	_, err = w.Write([]byte("application/epub+zip"))
	return err
}

// addFileToZip adiciona um arquivo ao ZIP
func (e *EPub) addFileToZip(zw *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	w, err := zw.Create(zipPath)
	if err != nil {
		return err
	}
	
	_, err = io.Copy(w, file)
	return err
}

// copyFile copia um arquivo
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	
	_, err = io.Copy(destination, source)
	return err
}

// defaultCSS retorna CSS padrão
func defaultCSS() string {
	return `@charset "UTF-8";

body {
  font-family: serif;
  font-size: 1em;
  line-height: 1.6;
  margin: 1em;
  text-align: justify;
}

h1, h2, h3, h4, h5, h6 {
  font-weight: bold;
  line-height: 1.2;
  margin-top: 1em;
  margin-bottom: 0.5em;
  text-align: left;
}

h1 {
  font-size: 2em;
}

h2 {
  font-size: 1.5em;
}

h3 {
  font-size: 1.17em;
}

p {
  margin: 0;
  text-indent: 1em;
}

p:first-child,
h1 + p,
h2 + p,
h3 + p {
  text-indent: 0;
}

img {
  max-width: 100%;
  height: auto;
}

.center {
  text-align: center;
}
`
}
