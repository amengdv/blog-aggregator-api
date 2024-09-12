package main

import (
	"net/http"
	"strconv"

	"github.com/amengdv/blog-aggregator-api/internal/database"
)

func (cfg *apiConfig) getPostsByUserHandler(w http.ResponseWriter, req *http.Request, user database.User) {
    limit := 10
    if queryLimit:= req.URL.Query().Get("limit"); queryLimit != "" {
        l, err := strconv.Atoi(queryLimit)
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "Limit query must be a number")
            return
        }
        limit = l
    }

    posts, err := cfg.DB.GetPostsByUser(req.Context(), database.GetPostsByUserParams{
        UserID: user.ID,
        Limit: int32(limit),
    })

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed To Get Posts")
        return
    }

    respondWithJson(w, http.StatusOK, dbPostsToPosts(posts))

}
