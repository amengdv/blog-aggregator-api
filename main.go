package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
    DB *database.Queries
    jwtSecret string
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("ERROR LOADING ENV FILE")
    }

    PORT := os.Getenv("PORT")
    dbURL := os.Getenv("POSTGRES")
    secret := os.Getenv("JWT_SECRET")

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
    }

    dbQueries := database.New(db)

    cfg := apiConfig{
        DB: dbQueries,
        jwtSecret: secret,
    }

    mux := http.NewServeMux()

    mux.HandleFunc("GET /v1/healthz", healthHandler)
    mux.HandleFunc("GET /v1/errors", errHandler)

    mux.HandleFunc("POST /v1/users", cfg.createUserHandler)
    mux.HandleFunc("PUT /v1/users", cfg.authMiddleware(cfg.updateUserHandler))
    mux.HandleFunc("DELETE /v1/users/{userID}", cfg.authMiddleware(cfg.deleteUserHandler))

    mux.HandleFunc("POST /v1/login", cfg.loginUserHandler)
    mux.HandleFunc("POST /v1/refresh", cfg.refreshTokenHandler)
    mux.HandleFunc("POST /v1/revoke", cfg.revokeRefreshTokenHandler)

    mux.HandleFunc("POST /v1/feeds", cfg.authMiddleware(cfg.createFeedsHandler))

    httpServer := &http.Server{
        Addr: ":" + PORT,
        Handler: mux,
    }

    fmt.Println("Starting Server and Listening on port ", PORT)
    log.Fatal(httpServer.ListenAndServe())
}
