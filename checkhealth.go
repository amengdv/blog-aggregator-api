package main

import (
	"net/http"

)

func healthHandler(w http.ResponseWriter, req *http.Request) {
    respondWithJson(w, http.StatusOK, struct{Status string `json:"status"`}{
        Status: "SERVER IS HEALTHY!!",
    })
}

func errHandler(w http.ResponseWriter, req *http.Request) {
    respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
