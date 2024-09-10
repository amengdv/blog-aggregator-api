package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
    decodeJsonError string = "Failure decoding request body"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
    // Expected Request
    type Request struct {
        Email string `json:"email"`
        Name string `json:"name"`
        Password string `json:"password"`
    }

    requestBody := Request{}

    err := decodeJson(req, &requestBody)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, decodeJsonError)
        return
    }

    err = validatePassword(requestBody.Password)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Password is not strong enough")
        return
    }

    hashedPass, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failure hashing password")
        return
    }

    dbUser, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Email: requestBody.Email,
        Name: requestBody.Name,
        Password: string(hashedPass),
    })
    
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failure Creating User")
        return
    }

    user := dbUserToUser(dbUser)

    type Respond struct {
        ID uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updatedAt"`
        Email string `json:"email"`
        Name string `json:"name"`
    }

    respondWithJson(w, http.StatusCreated, Respond{
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email: user.Email,
        Name: user.Name,
    })
}


func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, req *http.Request) {
    type Request struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    reqBody := Request{}

    err := decodeJson(req, &reqBody)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, decodeJsonError)
        return
    }

    userInfo, err := cfg.DB.GetUserByEmail(req.Context(), reqBody.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "User Does Not Exist")
        return
    }

    err = authenticatePassword(reqBody.Password, userInfo.Password)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Wrong Password!")
        return
    }

    //TODO Issue JWT
    token := issueJWT(userInfo.ID.String())
    jwtString, err := token.SignedString([]byte(cfg.jwtSecret))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Create Refresh Token
    refreshToken, err := generateRefreshToken()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token, please try again")
        return
    }

    // Update Refresh Token and Token Created At
    err = cfg.DB.UpdateRefreshToken(req.Context(), database.UpdateRefreshTokenParams{
        RefreshToken: sql.NullString{
            String: refreshToken,
            Valid: true,
        },
        TknExpiresAt: sql.NullTime{
            Time: time.Now().Add(1440 * time.Hour),
            Valid: true,
        },
        ID: userInfo.ID,
    })

    // Response
    type Response struct {
        ID uuid.UUID `json:"id"`
        Email string `json:"email"`
        JWToken string `json:"jwtoken"`
        RefreshToken string `json:"refresh_token"`
    }

    respondWithJson(w, http.StatusAccepted, Response{
        ID: userInfo.ID,
        Email: reqBody.Email,
        JWToken: jwtString,
        RefreshToken: refreshToken,
    })
}

func (cfg *apiConfig) deleteUserHandler(w http.ResponseWriter, req *http.Request, user database.User) {
    idQuery := req.PathValue("userID")
    if len(idQuery) == 0 {
        respondWithError(w, http.StatusInternalServerError, "No id provided")
        return
    }
    _, err := uuid.Parse(idQuery)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "ID is not a UUID type")
        return
    }

    if idQuery != user.ID.String() {
        respondWithError(w, http.StatusUnauthorized, "Unauthorized User")
        return
    }

    err = cfg.DB.DeleteUser(req.Context(), user.ID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Fail to delete user, try again")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

