package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/altar-12/rssagg/internal/database"
	"github.com/google/uuid"
)

/*
concurrency: number of rss feeds to fetch at the same time
timeBetweenFetches: time interval between triggering rss feed fetches
*/
func scraper(db *database.Queries, concurrency int, timeBetweenFetches time.Duration) {
	fmt.Printf("Starting job to fetch %v feeds in time intervals of %s\n", concurrency, timeBetweenFetches)
	ticker := time.NewTicker(timeBetweenFetches)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			fmt.Println("Error getting the feeds to fetch from the database", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeRSSFeed(db, feed, wg)
		}
		wg.Wait()
	}
}

func scrapeRSSFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("Error marking the feed with ID %v as fetched: %v\n", feed.ID, err)
		return
	}
	rssFeed, err := urlToRSSFeed(feed.Url)
	if err != nil {
		fmt.Printf("Error fetching the RSS feed from the URL %s: %v\n", feed.Url, err)
		return
	}
	for _, item := range rssFeed.Channel.Items {
		currentTime := time.Now().UTC()
		publishedAt := sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
		publishDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			publishedAt.Time = publishDate.UTC()
			publishedAt.Valid = true
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error in saving post to database: %v", err)
		}
	}
	log.Printf("Feed fetched %s, found total %v posts\n", feed.Name, len(rssFeed.Channel.Items))
}
