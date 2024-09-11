package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func validatePassword(pass string) error {
    overeight := len(pass) >= 8
    if !overeight {
        return errors.New("Password is not strong enough")
    }
    return nil
}

func authenticatePassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func issueJWT(id string) *jwt.Token {
    expiresAt := time.Now().Add(time.Hour)
    claims := jwt.RegisteredClaims{
        Issuer: "amengdv",
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(expiresAt),
        Subject: id,
    }

    return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func generateRefreshToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    rt := hex.EncodeToString(b)
    return rt, nil
}

func getAuthorizationBearer(req *http.Request) string {
    authHeader := req.Header.Get("Authorization")
    if len(authHeader) == 0 {
        return ""
    }

    splitAuth := strings.Split(authHeader, " ")
    if splitAuth[0] != "Bearer" {
        return ""
    }

    return splitAuth[1]

}

func (cfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authToken := r.Header.Get("Authorization")
        splitAuthToken := strings.Split(authToken, " ")
        if len(splitAuthToken) == 0 {
            respondWithError(w, http.StatusUnauthorized, "User Unauthorized: No token provided")
            return
        }
        if splitAuthToken[0] != "Bearer" {
            respondWithError(w, http.StatusUnauthorized, "User Unauthorized: Wrong key type")
            return
        }
        token := splitAuthToken[1]
        claims := jwt.RegisteredClaims{}
        jwToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(cfg.jwtSecret), nil
        })

        if err != nil {
            respondWithError(w, http.StatusUnauthorized, "User Unauthorized: Invalid token")
            return
        }

        idString, err := jwToken.Claims.GetSubject()
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed To Find User")
            return
        }

        id, err := uuid.Parse(idString)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed to Parse ID")
            return
        }

        user, err := cfg.DB.GetUser(r.Context(), id)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get user data: %v\n", err))
            return
        }

        handler(w, r, user)
    }
}
