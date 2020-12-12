package main

import (
	"log"
	"net/http"
	"os"

	"github.com/graphql-services/id/graph"
)

const defaultPort = "8080"

func main() {
	mux := graph.GetServerMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
