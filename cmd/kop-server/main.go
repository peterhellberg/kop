package main

import (
	"log"
	"net/http"
	"os"

	"github.com/peterhellberg/kop/list"
	"github.com/peterhellberg/kop/rpc"
)

const defaultPort = "12432"

func main() {
	server := rpc.NewServer()

	rpc.RegisterList(server, list.New())

	http.Handle(server.Basepath, server)

	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

func port() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	return defaultPort
}
