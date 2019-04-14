package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	log.Printf("Server running at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
