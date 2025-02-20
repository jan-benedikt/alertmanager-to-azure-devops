package parser

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

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
		"createURLTimerange": func(input string, duration int) string {
			// "2024-12-11T07:27:30Z"
			half_duration := duration / 2
			layout := "2006-01-02T15:04:05Z"
			t, _ := time.Parse(layout, input)
			t_start := t.Add((time.Duration(half_duration) * time.Hour) * -1)
			t_end := t.Add(time.Duration(half_duration) * time.Hour)
			// from=2025-02-17T23:00:00.000Z&to=2025-02-18T22:59:59.000Z
			return fmt.Sprintf("from=%s&to=%s", t_start, t_end)
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
