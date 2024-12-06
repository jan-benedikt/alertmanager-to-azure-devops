package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/lukas-blaha/alertmanager-to-azure-devops/pkg/parser"
)

const (
	webPort = 8080
	baseUrl = "https://dev.azure.com"
)

type Config struct {
	Org      string
	Project  string
	WorkItem string
	Template *template.Template
	Token    string
}

func main() {
	orgEnv := os.Getenv("ORGANIZATION")
	projectEnv := os.Getenv("PROJECT")
	workItemEnv := os.Getenv("WORKITEM")
	tmplEnv := os.Getenv("TEMPLATE")
	tokenEnv := os.Getenv("TOKEN")

	var org string
	flag.StringVar(&org, "organization", orgEnv, "Azure DevOps organization name")

	var project string
	flag.StringVar(&project, "project", projectEnv, "Azure DevOps project name")

	var workItem string
	flag.StringVar(&workItem, "workitem", workItemEnv, "Azure DevOps work item name")

	var tmplPath string
	flag.StringVar(&tmplPath, "template", tmplEnv, "Path to payload transformation template")

	var token string
	flag.StringVar(&token, "token", tokenEnv, "Authorization token")

	flag.Parse()

	if org == "" || project == "" || workItem == "" || tmplPath == "" || token == "" {
		msg := "env ORG or -organization\nenv PROJECT or -project\nenv WORKITEM or -workitem\nenv TEMPLATE or -template\nenv TOKEN or -token"
		log.Panicf("Missing required flags or environment variables. See settings below:\n\n%s", msg)
	}

	tmpl, err := parser.New(tmplPath)
	if err != nil {
		log.Panic(err)
	}

	app := Config{
		Org:      org,
		Project:  project,
		WorkItem: workItem,
		Template: tmpl,
		Token:    token,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

	fmt.Println("Server is running on port", webPort)
}
