package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
    mux.HandleFunc("GET /v1/feeds", cfg.getAllFeedsHandler)

    mux.HandleFunc("POST /v1/feed_follows", cfg.authMiddleware(cfg.followFeedsHandler))
    mux.HandleFunc("GET /v1/feed_follows", cfg.authMiddleware(cfg.getFeedFollowsHandler))
    mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.authMiddleware(cfg.deleteFeedFollowsHandler))

    mux.HandleFunc("GET /v1/posts/", cfg.authMiddleware(cfg.getPostsByUserHandler))

    httpServer := &http.Server{
        Addr: ":" + PORT,
        Handler: mux,
    }

    const limit = 10
    const interval = time.Second * 10
    go fetchWorkers(dbQueries, limit, interval)

    fmt.Println("Starting Server and Listening on port ", PORT)
    log.Fatal(httpServer.ListenAndServe())
}
