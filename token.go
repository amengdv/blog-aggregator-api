package main

import (
	"database/sql"
	"net/http"
	"time"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, req *http.Request) {

    refToken := getAuthorizationBearer(req)
    if refToken == "" {
        respondWithError(w, http.StatusUnauthorized, "Unauthorized User: Refresh Issue")
        return
    }

    user, err := cfg.DB.GetTokenInfo(req.Context(), sql.NullString{
        String: refToken,
        Valid: true,
    })
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Token Does Not Exist")
        return
    }

    // Check for expired
    if time.Since(user.TknExpiresAt.Time) >= days_60 {
        respondWithError(w, http.StatusInternalServerError, "Token has expired")
        return
    }

    token := issueJWT(user.ID.String())
    jwtString, err := token.SignedString([]byte(cfg.jwtSecret))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed Signing Token")
        return
    }

    type Response struct {
        Token string `json:"token"`
    }

    respondWithJson(w, http.StatusOK, Response{
        Token: jwtString,
    })
}

func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, req *http.Request) {
    refToken := getAuthorizationBearer(req)

    err := cfg.DB.RevokeToken(req.Context(), sql.NullString{
        String: refToken,
        Valid: true,
    })

    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Failed To Revoke Token: " + err.Error())
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
