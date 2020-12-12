package main

import (
	"github.com/akrylysov/algnhsa"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/graphql-services/id/graph"
)

func main() {
	xray.Configure(xray.Config{
		LogLevel:       "info", // default
		ServiceVersion: "1.2.3",
	})

	mux := graph.GetServerMux()

	// handler := gen.GetHTTPServeMux(src.New(db, &eventController), db, src.GetMigrations(db))
	algnhsa.ListenAndServe(mux, &algnhsa.Options{
		UseProxyPath: true,
	})
}
