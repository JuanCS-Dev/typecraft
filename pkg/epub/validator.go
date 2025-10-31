package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ValidationLevel representa o nível de severidade
type ValidationLevel string

const (
	LevelError   ValidationLevel = "ERROR"
	LevelWarning ValidationLevel = "WARNING"
	LevelInfo    ValidationLevel = "INFO"
)

// ValidationIssue representa um problema encontrado
type ValidationIssue struct {
	Level   ValidationLevel
	Code    string
	Message string
	File    string
	Line    int
}

// ValidationResult representa o resultado da validação
type ValidationResult struct {
	Valid   bool
	Issues  []ValidationIssue
	Version EPubVersion
}

// Validator valida arquivos ePub
type Validator struct {
	strictMode bool
}

// NewValidator cria um novo validador
func NewValidator() *Validator {
	return &Validator{
		strictMode: false,
	}
}

// NewStrictValidator cria um validador em modo strict
func NewStrictValidator() *Validator {
	return &Validator{
		strictMode: true,
	}
}

// ValidateFile valida um arquivo ePub
func (v *Validator) ValidateFile(path string) (*ValidationResult, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open epub: %w", err)
	}
	defer reader.Close()
	
	result := &ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
	}
	
	// 1. Validate structure
	v.validateStructure(reader, result)
	
	// 2. Validate mimetype
	v.validateMimetype(reader, result)
	
	// 3. Validate container.xml
	containerXML := v.findFile(reader, "META-INF/container.xml")
	if containerXML != nil {
		v.validateContainer(containerXML, result)
	}
	
	// 4. Parse OPF path from container
	opfPath := v.extractOPFPath(containerXML)
	if opfPath == "" {
		opfPath = "OEBPS/content.opf" // fallback
	}
	
	// 5. Validate OPF
	opfFile := v.findFile(reader, opfPath)
	if opfFile != nil {
		v.validateOPF(opfFile, result)
	}
	
	// 6. Validate NCX
	ncxFile := v.findFile(reader, "OEBPS/toc.ncx")
	if ncxFile != nil {
		v.validateNCX(ncxFile, result)
	}
	
	// 7. Check for nav.xhtml (EPUB 3)
	navFile := v.findFile(reader, "OEBPS/nav.xhtml")
	if navFile != nil {
		result.Version = EPub3
		v.validateNav(navFile, result)
	}
	
	// Determine if valid
	for _, issue := range result.Issues {
		if issue.Level == LevelError {
			result.Valid = false
			break
		}
	}
	
	return result, nil
}

// validateStructure valida a estrutura de diretórios
func (v *Validator) validateStructure(reader *zip.ReadCloser, result *ValidationResult) {
	requiredDirs := map[string]bool{
		"META-INF/": false,
		"OEBPS/":    false,
	}
	
	for _, file := range reader.File {
		for dir := range requiredDirs {
			if strings.HasPrefix(file.Name, dir) {
				requiredDirs[dir] = true
			}
		}
	}
	
	for dir, found := range requiredDirs {
		if !found {
			result.Issues = append(result.Issues, ValidationIssue{
				Level:   LevelError,
				Code:    "STRUCTURE_001",
				Message: fmt.Sprintf("Missing required directory: %s", dir),
			})
		}
	}
}

// validateMimetype valida o arquivo mimetype
func (v *Validator) validateMimetype(reader *zip.ReadCloser, result *ValidationResult) {
	// Must be first file
	if len(reader.File) == 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "MIMETYPE_001",
			Message: "No files in ePub",
		})
		return
	}
	
	first := reader.File[0]
	if first.Name != "mimetype" {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "MIMETYPE_002",
			Message: "mimetype must be first file in ZIP",
			File:    first.Name,
		})
	}
	
	// Must be uncompressed
	if first.Method != zip.Store {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "MIMETYPE_003",
			Message: "mimetype must be stored uncompressed",
			File:    "mimetype",
		})
	}
	
	// Validate content
	f, err := first.Open()
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "MIMETYPE_004",
			Message: "Failed to read mimetype",
		})
		return
	}
	defer f.Close()
	
	content, err := io.ReadAll(f)
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "MIMETYPE_005",
			Message: "Failed to read mimetype content",
		})
		return
	}
	
	if string(content) != "application/epub+zip" {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "MIMETYPE_006",
			Message: fmt.Sprintf("Invalid mimetype: %s", string(content)),
		})
	}
}

// validateContainer valida container.xml
func (v *Validator) validateContainer(file *zip.File, result *ValidationResult) {
	f, err := file.Open()
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "CONTAINER_001",
			Message: "Failed to open container.xml",
		})
		return
	}
	defer f.Close()
	
	content, err := io.ReadAll(f)
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "CONTAINER_002",
			Message: "Failed to read container.xml",
		})
		return
	}
	
	// Check for required elements
	requiredElements := []string{
		"<container",
		"<rootfiles>",
		"<rootfile",
		"full-path=",
		"media-type=\"application/oebps-package+xml\"",
	}
	
	contentStr := string(content)
	for _, elem := range requiredElements {
		if !strings.Contains(contentStr, elem) {
			result.Issues = append(result.Issues, ValidationIssue{
				Level:   LevelError,
				Code:    "CONTAINER_003",
				Message: fmt.Sprintf("Missing required element: %s", elem),
				File:    "META-INF/container.xml",
			})
		}
	}
}

// validateOPF valida content.opf
func (v *Validator) validateOPF(file *zip.File, result *ValidationResult) {
	f, err := file.Open()
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "OPF_001",
			Message: "Failed to open content.opf",
		})
		return
	}
	defer f.Close()
	
	content, err := io.ReadAll(f)
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "OPF_002",
			Message: "Failed to read content.opf",
		})
		return
	}
	
	// Parse XML to detect version
	type Package struct {
		Version string `xml:"version,attr"`
	}
	
	var pkg Package
	if err := xml.Unmarshal(content, &pkg); err == nil {
		if pkg.Version == "3.0" {
			result.Version = EPub3
		} else if pkg.Version == "2.0" {
			result.Version = EPub2
		}
	}
	
	// Check required elements
	requiredElements := []string{
		"<metadata",
		"<manifest",
		"<spine",
		"<dc:title>",
		"<dc:creator",
		"<dc:language>",
		"<dc:identifier",
	}
	
	contentStr := string(content)
	for _, elem := range requiredElements {
		if !strings.Contains(contentStr, elem) {
			result.Issues = append(result.Issues, ValidationIssue{
				Level:   LevelError,
				Code:    "OPF_003",
				Message: fmt.Sprintf("Missing required element: %s", elem),
				File:    file.Name,
			})
		}
	}
	
	// EPUB 3 specific validations
	if result.Version == EPub3 {
		if !strings.Contains(contentStr, "dcterms:modified") {
			if v.strictMode {
				result.Issues = append(result.Issues, ValidationIssue{
					Level:   LevelError,
					Code:    "OPF_EPUB3_001",
					Message: "EPUB 3 requires dcterms:modified metadata",
					File:    file.Name,
				})
			} else {
				result.Issues = append(result.Issues, ValidationIssue{
					Level:   LevelWarning,
					Code:    "OPF_EPUB3_001",
					Message: "EPUB 3 should include dcterms:modified metadata",
					File:    file.Name,
				})
			}
		}
	}
}

// validateNCX valida toc.ncx
func (v *Validator) validateNCX(file *zip.File, result *ValidationResult) {
	f, err := file.Open()
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelWarning,
			Code:    "NCX_001",
			Message: "Failed to open toc.ncx",
		})
		return
	}
	defer f.Close()
	
	content, err := io.ReadAll(f)
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelWarning,
			Code:    "NCX_002",
			Message: "Failed to read toc.ncx",
		})
		return
	}
	
	requiredElements := []string{
		"<ncx",
		"<head>",
		"<docTitle>",
		"<navMap>",
	}
	
	contentStr := string(content)
	for _, elem := range requiredElements {
		if !strings.Contains(contentStr, elem) {
			result.Issues = append(result.Issues, ValidationIssue{
				Level:   LevelWarning,
				Code:    "NCX_003",
				Message: fmt.Sprintf("Missing element: %s", elem),
				File:    "toc.ncx",
			})
		}
	}
}

// validateNav valida nav.xhtml
func (v *Validator) validateNav(file *zip.File, result *ValidationResult) {
	f, err := file.Open()
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "NAV_001",
			Message: "Failed to open nav.xhtml",
		})
		return
	}
	defer f.Close()
	
	content, err := io.ReadAll(f)
	if err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Code:    "NAV_002",
			Message: "Failed to read nav.xhtml",
		})
		return
	}
	
	requiredElements := []string{
		"<nav",
		"epub:type=\"toc\"",
	}
	
	contentStr := string(content)
	for _, elem := range requiredElements {
		if !strings.Contains(contentStr, elem) {
			result.Issues = append(result.Issues, ValidationIssue{
				Level:   LevelError,
				Code:    "NAV_003",
				Message: fmt.Sprintf("Missing required element: %s", elem),
				File:    "nav.xhtml",
			})
		}
	}
}

// findFile encontra um arquivo no ZIP
func (v *Validator) findFile(reader *zip.ReadCloser, name string) *zip.File {
	for _, file := range reader.File {
		if file.Name == name {
			return file
		}
	}
	return nil
}

// extractOPFPath extrai o caminho do OPF do container.xml
func (v *Validator) extractOPFPath(file *zip.File) string {
	if file == nil {
		return ""
	}
	
	f, err := file.Open()
	if err != nil {
		return ""
	}
	defer f.Close()
	
	content, err := io.ReadAll(f)
	if err != nil {
		return ""
	}
	
	// Simple extraction - look for full-path attribute
	contentStr := string(content)
	start := strings.Index(contentStr, "full-path=\"")
	if start == -1 {
		return ""
	}
	start += len("full-path=\"")
	
	end := strings.Index(contentStr[start:], "\"")
	if end == -1 {
		return ""
	}
	
	return contentStr[start : start+end]
}

// GetIssuesByLevel filtra issues por nível
func (r *ValidationResult) GetIssuesByLevel(level ValidationLevel) []ValidationIssue {
	filtered := make([]ValidationIssue, 0)
	for _, issue := range r.Issues {
		if issue.Level == level {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// HasErrors retorna true se houver erros
func (r *ValidationResult) HasErrors() bool {
	return len(r.GetIssuesByLevel(LevelError)) > 0
}

// HasWarnings retorna true se houver warnings
func (r *ValidationResult) HasWarnings() bool {
	return len(r.GetIssuesByLevel(LevelWarning)) > 0
}

// Summary retorna um resumo da validação
func (r *ValidationResult) Summary() string {
	errors := len(r.GetIssuesByLevel(LevelError))
	warnings := len(r.GetIssuesByLevel(LevelWarning))
	infos := len(r.GetIssuesByLevel(LevelInfo))
	
	return fmt.Sprintf("Errors: %d, Warnings: %d, Info: %d", errors, warnings, infos)
}
