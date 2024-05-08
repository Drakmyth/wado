package server

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const DEFAULT_HOST = "localhost"
const DEFAULT_PORT = "8080"

var tmpl *template.Template = nil

func ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /example", getExample)
	var err error = nil
	tmpl, err = template.ParseGlob("./templates/*.tmpl.html")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
	}

	mux.Handle("/", http.FileServer(http.Dir("public")))

	host, host_set := os.LookupEnv("HOST")
	if !host_set {
		host = DEFAULT_HOST
	}

	port, port_set := os.LookupEnv("PORT")
	if !port_set {
		port = DEFAULT_PORT
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	slog.Info("go-template is listening,", "addr", addr)
	return http.ListenAndServe(addr, mux)
}

func getExample(w http.ResponseWriter, r *http.Request) {
	slog.Info("Getting example")

	tmpl.ExecuteTemplate(w, "example-page", "Hello")
}
