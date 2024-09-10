package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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
