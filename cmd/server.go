package main

import (
	"bytes"
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

	mux.HandleFunc("POST /", app.GetTemplate)

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

	err = app.CreateTicket(s)
	if err != nil {
		log.Println(err)
		return
	}
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

func (app *Config) CreateTicket(payload string) error {
	req, err := http.NewRequest("POST", app.Target, bytes.NewReader([]byte(payload)))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending POST request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return nil
}
