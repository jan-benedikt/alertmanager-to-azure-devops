package parser

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	amt "github.com/prometheus/alertmanager/template"
)

func New(path string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"getHost": func(input, separator string) string {
			sl := strings.Split(input, separator)
			if len(sl) < 3 {
				return input
			}
			return strings.Join(sl[:3], separator)
		},
		"removePort": func(input string) string {
			return strings.Split(input, ":")[0]
		},
	}
	tmpl := template.New("template").Funcs(funcMap)

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %v", err)
	}

	t, err := tmpl.Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template file: %v", err)
	}

	return t, err
}

func Render(t *template.Template, data amt.Data) (string, error) {
	var b bytes.Buffer

	err := t.Execute(&b, data)
	if err != nil {
		return "", err
	}

	return b.String(), err
}
