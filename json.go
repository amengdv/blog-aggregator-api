package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
    if code >= 500 {
        log.Println("RESPONDED WITH 5XX")
    }

    type errorResponse struct {
        Error string `json:"error"`
    }

    respondWithJson(w, code, errorResponse{
        Error: msg,
    })
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("ERROR! Could not marshal to JSON: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(code)
    w.Write(data)
}

