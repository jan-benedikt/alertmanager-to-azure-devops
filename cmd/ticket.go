package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Result struct {
	QueryType string   `json:"queryType"`
	AsOf      string   `json:"asOf"`
	Columns   []Column `json:"columns"`
	WorkItems []Ticket `json:"workItems"`
}

type Column struct {
	Name          string `json:"name"`
	ReferenceName string `json:"referenceName"`
	Url           string `json:"url"`
}

type Ticket struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

func replaceBlanks(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}

func (app *Config) CreateUrl(op string) string {
	var (
		urlPath string
		params  string
		wrktItm string
	)
	if op == "create" {
		urlPath = "_apis/wit/workitems/$"
		params = "?api-version=7.1"
		wrktItm = replaceBlanks(app.WorkItem)
	} else if op == "get" {
		urlPath = "_apis/wit/wiql"
		params = "?api-version=6.0"
		wrktItm = ""
	} else if op == "update" {
		urlPath = "_apis/wit/workitems/"
		params = "?api-version=7.1"
		wrktItm = replaceBlanks(app.WorkItem)
	}

	url := fmt.Sprintf(
		"%s/%s/%s/%s%s%s",
		baseUrl,
		replaceBlanks(app.Org),
		replaceBlanks(app.Project),
		urlPath,
		wrktItm,
		params,
	)

	return url
}

func (app *Config) MakeRequest(method, url, payload, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (app *Config) CreateTicket(payload string) error {
	// "https://dev.azure.com/<organization>/<project>/_apis/wit/workitems/<workitem>?api-version=7.1"
	url := app.CreateUrl("create")

	resp, err := app.MakeRequest("POST", url, payload, "application/json-patch+json")
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return nil
}

func (app *Config) GetTicket(id string) (Ticket, error) {
	var result Result

	// "https://dev.azure.com/<organization>/<project>/_apis/wit/wiql?api-version=6.0"
	url := app.CreateUrl("get")

	query := fmt.Sprintf(`{"query": "SELECT [System.Id] FROM WorkItems WHERE [PEScrum.WeblistName] CONTAINS \"%s\""}`, id)

	resp, err := app.MakeRequest("POST", url, query, "application/json")
	if err != nil {
		return Ticket{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)
	if err != nil {
		return Ticket{}, err
	}

	if resp.StatusCode != 200 {
		return Ticket{}, err
	}

	if len(result.WorkItems) == 0 {
		return Ticket{}, nil
	}
	return result.WorkItems[0], nil
}

func (app *Config) CloseTicket(ticket Ticket) error {
	// "https://dev.azure.com/<organization>/<project>/_apis/wit/workitems/<ticket.id>?api-version=7.1"
	url := fmt.Sprintf("%s%s", ticket.Url, "?api-version=7.1")

	resp, err := app.MakeRequest("PATCH", url, app.CloseTemplate, "application/json-patch+json")
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return nil
}
