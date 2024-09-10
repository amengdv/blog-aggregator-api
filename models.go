package main

import (
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/google/uuid"
)

type User struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
    Name string `json:"name"`
    Password string `json:""`
}

func dbUserToUser(dbUser database.User) User {
    return User{
        ID: dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email: dbUser.Email,
        Name: dbUser.Name,
        Password: dbUser.Password,
    }
}
