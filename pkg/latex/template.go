package latex

import (
	"fmt"
	"regexp"
	"strings"
)

// Template representa um template LaTeX
type Template struct {
	content   string
	variables map[string]interface{}
	functions map[string]TemplateFunc
}

// TemplateFunc é uma função de template
type TemplateFunc func(args ...interface{}) (string, error)

// NewTemplate cria um novo template
func NewTemplate(content string) *Template {
	return &Template{
		content:   content,
		variables: make(map[string]interface{}),
		functions: make(map[string]TemplateFunc),
	}
}

// SetVariable define uma variável
func (t *Template) SetVariable(name string, value interface{}) {
	t.variables[name] = value
}

// SetVariables define múltiplas variáveis
func (t *Template) SetVariables(vars map[string]interface{}) {
	for k, v := range vars {
		t.variables[k] = v
	}
}

// GetVariable obtém uma variável
func (t *Template) GetVariable(name string) (interface{}, bool) {
	val, ok := t.variables[name]
	return val, ok
}

// RegisterFunction registra uma função customizada
func (t *Template) RegisterFunction(name string, fn TemplateFunc) {
	t.functions[name] = fn
}

// Render renderiza o template
func (t *Template) Render() (string, error) {
	result := t.content
	
	// 1. Process conditionals
	var err error
	result, err = t.processConditionals(result)
	if err != nil {
		return "", err
	}
	
	// 2. Process loops
	result, err = t.processLoops(result)
	if err != nil {
		return "", err
	}
	
	// 3. Process function calls
	result, err = t.processFunctions(result)
	if err != nil {
		return "", err
	}
	
	// 4. Process variable substitution
	result = t.processVariables(result)
	
	return result, nil
}

// processVariables substitui variáveis {{var}}
func (t *Template) processVariables(content string) string {
	re := regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)
	
	return re.ReplaceAllStringFunc(content, func(match string) string {
		varName := strings.Trim(match, "{}")
		
		if val, ok := t.variables[varName]; ok {
			return fmt.Sprintf("%v", val)
		}
		
		return match // Keep original if not found
	})
}

// processConditionals processa {{if var}}...{{end}}
func (t *Template) processConditionals(content string) (string, error) {
	result := content
	
	// Pattern: {{if varname}}...{{else}}...{{end}} (process FIRST!)
	reElse := regexp.MustCompile(`(?s)\{\{if\s+([a-zA-Z0-9_]+)\}\}(.*?)\{\{else\}\}(.*?)\{\{end\}\}`)
	matches := reElse.FindAllStringSubmatch(result, -1)
	
	for _, match := range matches {
		fullMatch := match[0]
		varName := match[1]
		trueContent := match[2]
		falseContent := match[3]
		
		if val, ok := t.variables[varName]; ok && t.isTruthy(val) {
			result = strings.Replace(result, fullMatch, trueContent, 1)
		} else {
			result = strings.Replace(result, fullMatch, falseContent, 1)
		}
	}
	
	// Pattern: {{if varname}}...{{end}} (process AFTER!)
	re := regexp.MustCompile(`(?s)\{\{if\s+([a-zA-Z0-9_]+)\}\}(.*?)\{\{end\}\}`)
	matches = re.FindAllStringSubmatch(result, -1)
	
	for _, match := range matches {
		fullMatch := match[0]
		varName := match[1]
		innerContent := match[2]
		
		// Check if variable exists and is truthy
		if val, ok := t.variables[varName]; ok && t.isTruthy(val) {
			result = strings.Replace(result, fullMatch, innerContent, 1)
		} else {
			result = strings.Replace(result, fullMatch, "", 1)
		}
	}
	
	return result, nil
}

// processLoops processa {{range items}}...{{end}}
func (t *Template) processLoops(content string) (string, error) {
	// Pattern: {{range varname}}...{{.}}...{{end}}
	re := regexp.MustCompile(`(?s)\{\{range\s+([a-zA-Z0-9_]+)\}\}(.*?)\{\{end\}\}`)
	
	result := content
	matches := re.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		fullMatch := match[0]
		varName := match[1]
		template := match[2]
		
		val, ok := t.variables[varName]
		if !ok {
			result = strings.Replace(result, fullMatch, "", 1)
			continue
		}
		
		// Convert to slice
		var items []interface{}
		switch v := val.(type) {
		case []interface{}:
			items = v
		case []string:
			for _, s := range v {
				items = append(items, s)
			}
		case []int:
			for _, i := range v {
				items = append(items, i)
			}
		default:
			return "", fmt.Errorf("range variable must be a slice, got %T", val)
		}
		
		// Render each item
		var rendered strings.Builder
		for _, item := range items {
			itemStr := strings.ReplaceAll(template, "{{.}}", fmt.Sprintf("%v", item))
			rendered.WriteString(itemStr)
		}
		
		result = strings.Replace(result, fullMatch, rendered.String(), 1)
	}
	
	return result, nil
}

// processFunctions processa {{func arg1 arg2}}
func (t *Template) processFunctions(content string) (string, error) {
	// Pattern: {{funcname arg1 arg2}}
	re := regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\s+([^}]+)\}\}`)
	
	result := content
	matches := re.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		fullMatch := match[0]
		funcName := match[1]
		argsStr := match[2]
		
		// Check if it's a registered function
		fn, ok := t.functions[funcName]
		if !ok {
			continue // Not a function, might be a variable
		}
		
		// Parse arguments
		args := strings.Fields(argsStr)
		var processedArgs []interface{}
		for _, arg := range args {
			// Try to resolve as variable
			if val, ok := t.variables[arg]; ok {
				processedArgs = append(processedArgs, val)
			} else {
				processedArgs = append(processedArgs, arg)
			}
		}
		
		// Call function
		output, err := fn(processedArgs...)
		if err != nil {
			return "", fmt.Errorf("function %s error: %w", funcName, err)
		}
		
		result = strings.Replace(result, fullMatch, output, 1)
	}
	
	return result, nil
}

// isTruthy verifica se um valor é truthy
func (t *Template) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case int, int8, int16, int32, int64:
		return v != 0
	case float32, float64:
		return v != 0.0
	case []interface{}:
		return len(v) > 0
	case []string:
		return len(v) > 0
	default:
		return true
	}
}

// Escape escapa caracteres especiais LaTeX
func Escape(text string) string {
	// Order matters! Backslash must be first
	text = strings.ReplaceAll(text, "\\", "\\textbackslash{}")
	text = strings.ReplaceAll(text, "&", "\\&")
	text = strings.ReplaceAll(text, "%", "\\%")
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "~", "\\textasciitilde{}")
	text = strings.ReplaceAll(text, "^", "\\textasciicircum{}")
	
	return text
}

// TemplateBuilder facilita construção de templates
type TemplateBuilder struct {
	sections []string
}

// NewTemplateBuilder cria um novo builder
func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{
		sections: make([]string, 0),
	}
}

// AddSection adiciona uma seção
func (tb *TemplateBuilder) AddSection(content string) *TemplateBuilder {
	tb.sections = append(tb.sections, content)
	return tb
}

// AddVariable adiciona uma variável
func (tb *TemplateBuilder) AddVariable(name string) *TemplateBuilder {
	tb.sections = append(tb.sections, fmt.Sprintf("{{%s}}", name))
	return tb
}

// AddConditional adiciona um if
func (tb *TemplateBuilder) AddConditional(varName, content string) *TemplateBuilder {
	tb.sections = append(tb.sections, fmt.Sprintf("{{if %s}}%s{{end}}", varName, content))
	return tb
}

// AddLoop adiciona um range
func (tb *TemplateBuilder) AddLoop(varName, itemTemplate string) *TemplateBuilder {
	tb.sections = append(tb.sections, fmt.Sprintf("{{range %s}}%s{{end}}", varName, itemTemplate))
	return tb
}

// Build constrói o template
func (tb *TemplateBuilder) Build() string {
	return strings.Join(tb.sections, "\n")
}

// StandardFunctions retorna funções padrão
func StandardFunctions() map[string]TemplateFunc {
	return map[string]TemplateFunc{
		"upper": func(args ...interface{}) (string, error) {
			if len(args) != 1 {
				return "", fmt.Errorf("upper expects 1 argument")
			}
			return strings.ToUpper(fmt.Sprintf("%v", args[0])), nil
		},
		"lower": func(args ...interface{}) (string, error) {
			if len(args) != 1 {
				return "", fmt.Errorf("lower expects 1 argument")
			}
			return strings.ToLower(fmt.Sprintf("%v", args[0])), nil
		},
		"title": func(args ...interface{}) (string, error) {
			if len(args) != 1 {
				return "", fmt.Errorf("title expects 1 argument")
			}
			return strings.Title(strings.ToLower(fmt.Sprintf("%v", args[0]))), nil
		},
		"escape": func(args ...interface{}) (string, error) {
			if len(args) != 1 {
				return "", fmt.Errorf("escape expects 1 argument")
			}
			return Escape(fmt.Sprintf("%v", args[0])), nil
		},
		"repeat": func(args ...interface{}) (string, error) {
			if len(args) != 2 {
				return "", fmt.Errorf("repeat expects 2 arguments")
			}
			
			str := fmt.Sprintf("%v", args[0])
			count, ok := args[1].(int)
			if !ok {
				return "", fmt.Errorf("repeat count must be int")
			}
			
			return strings.Repeat(str, count), nil
		},
	}
}
