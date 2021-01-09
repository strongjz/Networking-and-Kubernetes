package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

func main() {
	http.HandleFunc("/", hello)

	http.HandleFunc("/healthz", healthz)

	http.ListenAndServe("0.0.0.0:8080", nil)
}