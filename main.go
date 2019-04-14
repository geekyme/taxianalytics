package main

import (
	"log"
	"net/http"
	"os"
	"taxianalytics/internal/taxidata"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	subscription := taxidata.TaxiSubscription
	go taxidata.Subscribe(subscription)
	http.HandleFunc("/", taxidata.Handler)

	log.Printf("Server running at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
