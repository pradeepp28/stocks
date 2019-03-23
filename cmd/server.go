package main

import (
	"fmt"
	"log"
	"net/http"
	"stocks"
)

const port = 8080

var (
	s stocks.Server
)

func main() {
	s = stocks.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/stock/", handleRequest)
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		get(w, r)
	default:
		w.WriteHeader(405)
		fmt.Fprint(w, "Request Method not supported")
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	s.Get(w, r)
}
