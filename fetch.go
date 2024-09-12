package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/amengdv/blog-aggregator-api/internal/database"
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

func printTitleFromItem(rss RSS) {
    items := rss.Channel.Items
    for i, item := range items {
        fmt.Printf("%v, %d: %v\n",rss.Channel.Title, i, item.Title)
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

    printTitleFromItem(rss)
}
