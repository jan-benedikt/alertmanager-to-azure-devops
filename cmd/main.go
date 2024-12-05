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

const webPort = 8080

type Config struct {
	Target   string
	Template *template.Template
}

func main() {
	urlEnv := os.Getenv("DEVOPS_URL")
	tmplEnv := os.Getenv("DEVOPS_TEMPLATE")

	var url string
	flag.StringVar(&url, "target", urlEnv, "Target URL")

	var tmplPath string
	flag.StringVar(&tmplPath, "template", tmplEnv, "Path to payload transformation template")

	flag.Parse()

	if url == "" || tmplPath == "" {
		msg := "env TARGET or -target\nenv TEMPLATE or -template"
		log.Panicf("Missing required flags or environment variables. See settings below:\n\n%s", msg)
	}

	log.Printf("Bind address: http://localhost:%d", webPort)
	log.Printf("Target address: %v", url)
	log.Printf("Template path: %v", tmplPath)

	tmpl, err := parser.New(tmplPath)
	if err != nil {
		log.Panic(err)
	}

	app := Config{
		Target:   url,
		Template: tmpl,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
