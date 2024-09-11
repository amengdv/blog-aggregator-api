package main

import (
	"net/http"
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) followFeedsHandler(w http.ResponseWriter, req *http.Request, user database.User) {
    type Request struct {
        FeedID uuid.UUID `json:"feed_id"`
    }

    reqBody := Request{}

    if err := decodeJson(req, &reqBody); err != nil {
        respondWithError(w, http.StatusInternalServerError, decodeJsonError)
        return
    }

    dbFeedFollow, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        UserID: user.ID,
        FeedID: reqBody.FeedID,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    })

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create DB Follows")
        return
    }

    respondWithJson(w, http.StatusCreated, dbFeedFollowToFeedFollow(dbFeedFollow))

}

func (cfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, req *http.Request, user database.User) {
    feedFollows, err := cfg.DB.GetFeedFollowByID(req.Context(), user.ID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed To Get Followed Feeds")
        return
    }

    respondWithJson(w, http.StatusOK, dbFeedFollowsToFeedFollows(feedFollows))
}

func (cfg *apiConfig) deleteFeedFollowsHandler(w http.ResponseWriter, req *http.Request, user database.User) {
    feedFollowID := req.PathValue("feedFollowID")
    if len(feedFollowID) == 0 {
        respondWithError(w, http.StatusBadRequest, "Followed Feed ID not provided")
        return
    }

    feedFollowParsed, err := uuid.Parse(feedFollowID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Query Provided is not a type UUID")
        return
    }
    
    err = cfg.DB.DeleteFeedFollow(req.Context(), database.DeleteFeedFollowParams{
        ID: feedFollowParsed,
        UserID: user.ID,
    })
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Cannot Delete Feed Follow")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
