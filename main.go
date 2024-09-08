package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("ERROR LOADING ENV FILE")
    }

    PORT := os.Getenv("PORT")

    mux := http.NewServeMux()

    mux.HandleFunc("GET /v1/healthz", healthHandler)
    mux.HandleFunc("GET /v1/errors", errHandler)

    httpServer := &http.Server{
        Addr: ":" + PORT,
        Handler: mux,
    }

    fmt.Println("Starting Server and Listening on port ", PORT)
    log.Fatal(httpServer.ListenAndServe())
}
