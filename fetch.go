package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
	"github.com/google/uuid"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []Item    `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

func fetchXMLFromUrl(link string) (RSS, error) {
    res, err := http.Get(link)
    if err != nil {
        return RSS{}, err
    }

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return RSS{}, err
    }

    rss := RSS{}
    if err = xml.Unmarshal(body, &rss); err != nil {
        return RSS{}, err
    }

    return rss, nil
}

func createAndStorePost(db *database.Queries, rss RSS, feed database.Feed) {
    items := rss.Channel.Items
    for _, item := range items {
        publishedAt := sql.NullTime{}
        if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
            publishedAt = sql.NullTime{
                Time: t,
                Valid: true,
            }
        }

        _, err := db.CreatePost(context.Background(), database.CreatePostParams{
            ID: uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Title: item.Title,
            Url: item.Link,
            Description: sql.NullString {
                String: item.Description,
            },
            PublishedAt: publishedAt,
            FeedID: feed.ID,
        })

        if err != nil {
            if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
                continue
            }
            log.Printf("Couldn't create post: %v\n", err)
            continue
        }
    }
}

func fetchWorkers(db *database.Queries, limit int32, interval time.Duration) {
    ticker := time.NewTicker(interval)

    for range ticker.C {
        feeds, err := db.GetNextFeedToFetch(context.Background(), limit)
        if err != nil {
            log.Printf("Error Getting Next Feed To Fetch: %v\n", err)
            return
        }

        wg := &sync.WaitGroup{}
        for _, feed := range feeds {
            fmt.Printf("Fetching feed: %v\n", feed.Name)
            wg.Add(1)
            go processFeeds(db, wg, feed)
        }
        wg.Wait()
    }
}

func processFeeds(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
    defer wg.Done()

    rss, err := fetchXMLFromUrl(feed.Url)
    if err != nil {
        return
    }

    _, err = db.MarkFeedFetched(context.Background(), feed.ID)
    if err != nil {
        return
    }

    createAndStorePost(db, rss, feed)
}
