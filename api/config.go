package main

import (
	"os"
)

type Config struct {
	ApiLink string
	Host    string
	Port    string
	UI      struct {
		TmplDir string
		CSSDir  string
	}
}

func newConfig() *Config {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "0.0.0.0"
	}

	apiLink := os.Getenv("API_LINK")
	if len(apiLink) == 0 {
		apiLink = "http://localhost:8000"
	}

	return &Config{
		ApiLink: apiLink,
		Host:    host,
		Port:    port,
		UI: struct {
			TmplDir string
			CSSDir  string
		}{TmplDir: "./ui/html/", CSSDir: "./ui/static/"},
	}
}
