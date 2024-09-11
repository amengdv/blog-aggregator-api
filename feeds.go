package main

import (
	"net/http"
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeedsHandler(w http.ResponseWriter, req *http.Request, user database.User) {
    type Request struct {
        Name string `json:"feeds_name"`
        Url string `json:"feeds_url"`
    }

    reqBody := Request{}

    err := decodeJson(req, &reqBody)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, decodeJsonError)
        return
    }

    feed, err := cfg.DB.CreateFeeds(req.Context(), database.CreateFeedsParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: reqBody.Name,
        Url: reqBody.Url,
        UserID: user.ID,
    })
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed Creating Feeds")
        return
    }

    respondWithJson(w, http.StatusCreated, databaseFeedToFeed(feed))
}
