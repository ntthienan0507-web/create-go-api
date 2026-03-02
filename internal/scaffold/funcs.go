package scaffold

import (
	"strings"
	"text/template"
	"unicode"
)

// funcMap provides helper functions available in all templates.
var funcMap = template.FuncMap{
	"snakeCase": toSnakeCase,
	"kebabCase": toKebabCase,
	"lower":     strings.ToLower,
	"upper":     strings.ToUpper,
	"contains":  strings.Contains,
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else if r == '-' || r == ' ' {
			result = append(result, '_')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func toKebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '-')
			}
			result = append(result, unicode.ToLower(r))
		} else if r == '_' || r == ' ' {
			result = append(result, '-')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
