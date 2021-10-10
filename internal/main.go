package main

import (
	"fmt"
	"net/http"

	v1 "github.com/crab21/middleware/internal/jaeger/v1"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	v1.InitClients()
	fmt.Println("starting")
	http.HandleFunc("/good", HelloHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
