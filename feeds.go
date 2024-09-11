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

    dbFeedFollow, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        UserID: user.ID,
        FeedID: feed.ID,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    })

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create DB Follows")
        return
    }

    type Response struct {
        Feed Feed `json:"feed"`
        FeedFollow FeedFollow `json:"feed_follow"`
    }

    respondWithJson(w, http.StatusCreated, Response{
        Feed: databaseFeedToFeed(feed),
        FeedFollow: dbFeedFollowToFeedFollow(dbFeedFollow),
    })
}

func (cfg *apiConfig) getAllFeedsHandler(w http.ResponseWriter, req *http.Request) {

    dbFeeds, err := cfg.DB.GetFeeds(req.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed To Get Feeds")
        return
    }

    respondWithJson(w, http.StatusOK, dbFeedsToFeeds(dbFeeds))
}
