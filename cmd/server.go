package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lukas-blaha/alertmanager-to-azure-devops/pkg/parser"
	amt "github.com/prometheus/alertmanager/template"
)

func (app *Config) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.GetTemplate)

	return mux
}

func (app *Config) GetTemplate(w http.ResponseWriter, r *http.Request) {
	data, err := Decode(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	s, err := parser.Render(app.Template, data)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Print(s)
}

func Decode(body io.Reader) (amt.Data, error) {
	var data amt.Data

	if body == nil {
		return amt.Data{}, nil
	}

	decoder := json.NewDecoder(body)

	err := decoder.Decode(&data)

	return data, err
}
