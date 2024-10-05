package main

import (
	"log"

	"github.com/sumant1122/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8888")
	log.Fatal(srv.ListenAndServe())
}
