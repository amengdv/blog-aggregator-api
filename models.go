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

type UserNoPass struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
    Name string `json:"name"`
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

func respondWithUserSec(dbUser database.User) UserNoPass {
    return UserNoPass{
        ID: dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email: dbUser.Email,
        Name: dbUser.Name,
    }
}

type Feed struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Name      string `json:"name"`
    Url       string `json:"url"`
    UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
    return Feed{
        ID: dbFeed.ID,
        CreatedAt: dbFeed.CreatedAt,
        UpdatedAt: dbFeed.UpdatedAt,
        Name: dbFeed.Name,
        Url: dbFeed.Url,
        UserID: dbFeed.UserID,
    }
}

func dbFeedsToFeeds(dbFeeds []database.Feed) (feeds []Feed) {
    for _, dbFeed := range dbFeeds {
        feeds = append(feeds, databaseFeedToFeed(dbFeed))
    }

    return feeds;
}

type FeedFollow struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    FeedID    uuid.UUID `json:"feed_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func dbFeedFollowToFeedFollow(ff database.FeedFollow) FeedFollow {
    return FeedFollow {
        ID: ff.ID,
        UserID: ff.UserID,
        FeedID: ff.FeedID,
        CreatedAt: ff.CreatedAt,
        UpdatedAt: ff.UpdatedAt,
    }
}


func dbFeedFollowsToFeedFollows(ffs []database.FeedFollow) (feedFollows []FeedFollow) {
    for _, ff := range ffs {
        feedFollows = append(feedFollows, dbFeedFollowToFeedFollow(ff))
    }

    return feedFollows
}


