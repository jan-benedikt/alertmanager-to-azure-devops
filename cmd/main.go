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
	Org            string
	Project        string
	WorkItem       string
	CreateTemplate *template.Template
	CloseTemplate  string
	Token          string
}

func main() {
	orgEnv := os.Getenv("ORGANIZATION")
	projectEnv := os.Getenv("PROJECT")
	workItemEnv := os.Getenv("WORKITEM")
	createTmplEnv := os.Getenv("CREATE_TEMPLATE")
	closeTmplEnv := os.Getenv("CLOSE_TEMPLATE")
	tokenEnv := os.Getenv("TOKEN")

	var org string
	flag.StringVar(&org, "organization", orgEnv, "Azure DevOps organization name")

	var project string
	flag.StringVar(&project, "project", projectEnv, "Azure DevOps project name")

	var workItem string
	flag.StringVar(&workItem, "workitem", workItemEnv, "Azure DevOps work item name")

	var createTmplPath string
	flag.StringVar(&createTmplPath, "create-template", createTmplEnv, "Path to payload transformation template")

	var closeTmplPath string
	flag.StringVar(&closeTmplPath, "close-template", closeTmplEnv, "Path to payload transformation template")

	var token string
	flag.StringVar(&token, "token", tokenEnv, "Authorization token")

	flag.Parse()

	if org == "" || project == "" || workItem == "" || createTmplPath == "" || closeTmplPath == "" || token == "" {
		msg := "env ORG or -organization\nenv PROJECT or -project\nenv WORKITEM or -workitem\nenv CREATE_TEMPLATE or -create-template\nenv CLOSE_TEMPLATE or -close-template\nenv TOKEN or -token"
		log.Panicf("Missing required flags or environment variables. See settings below:\n\n%s", msg)
	}

	createTmpl, err := parser.New(createTmplPath)
	if err != nil {
		log.Panic(err)
	}

	closeTmpl, err := os.ReadFile(closeTmplPath)
	if err != nil {
		log.Panic(err)
	}

	app := Config{
		Org:            org,
		Project:        project,
		WorkItem:       workItem,
		CreateTemplate: createTmpl,
		CloseTemplate:  string(closeTmpl),
		Token:          token,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	fmt.Println("Server is running on port", webPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Cannot start server:", err)
	}
}
