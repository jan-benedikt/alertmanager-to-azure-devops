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

	mux.HandleFunc("POST /", app.GetTemplate)

	return mux
}

func (app *Config) GetTemplate(w http.ResponseWriter, r *http.Request) {
	data, err := Decode(r.Body)
	if err != nil {
		log.Println("Cannot decode request:", err)
		return
	}

	s, err := parser.Render(app.CreateTemplate, data)
	if err != nil {
		log.Println("Cannot render template:", err)
		return
	}

	ticket, err := app.GetTicket(data.Alerts[0].Fingerprint)
	if err != nil {
		log.Println("Could not get ticket:", err)
		return
	}

	if data.Alerts[0].Status == "firing" {
		if ticket == (Ticket{}) {
			fmt.Println("Creating ticket for grafana alert:", data.Alerts[0].Fingerprint)
			err = app.CreateTicket(s)
			if err != nil {
				log.Println("Cannot create ticket:", err)
				return
			}
		}
	} else if data.Alerts[0].Status == "resolved" {
		if ticket != (Ticket{}) {
			fmt.Println("Closing ticket for grafana alert:", data.Alerts[0].Fingerprint)
			err = app.CloseTicket(ticket)
			if err != nil {
				log.Println("Cannot close ticket:", err)
				return
			}
		}
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
