package v1handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
	"github.com/google/uuid"
)

const (
	FETCH_N        = 10
	FETCH_INTERVAL = 60
)

func fetchRssAndUpdateDB(feed database.Feed, wg *sync.WaitGroup, v1 *V1Handlers) {
	defer wg.Done()
	rss, err := FetchRss(feed.Url)
	if err != nil {
		log.Printf("Unable to fetch RSS from %v: %v", feed.Url, err)
		return
	}
	// fmt.Println("Fetched RSS", rss)

	for _, rssItem := range rss.Channel.Item {

		description := sql.NullString{}
		if rssItem.Description != nil {
			description.Valid = true
			description.String = *rssItem.Description
		}
		publishedAt := sql.NullTime{}
		if rssItem.PubDate != nil {
			pubTime, err := time.Parse(time.RFC1123Z, *rssItem.PubDate)
			if err == nil {
				publishedAt.Valid = true
				publishedAt.Time = pubTime
			}
		}
		post, err := v1.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			log.Printf("Unable to create post: %v\n", err)
			continue
		}

		fmt.Println("Post Created:", post)
	}

	updatedFeed, err := v1.DB.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Unable to mark update in DB for %v: %v", feed.Url, err)
		return
	}
	fmt.Println("Updated Feed", updatedFeed)
}

func RssWorker(v1 *V1Handlers) {
	ticker := time.NewTicker(time.Duration(FETCH_INTERVAL) * time.Second)
	for range ticker.C {
		log.Println("Starting another fetch loop")
		feedsToFetch, err := v1.DB.GetNextFeedsToFetch(context.Background(), FETCH_N)
		if err != nil {
			log.Printf("Unable to retrieve feeds to be fetched: %v\n", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feedToFetch := range feedsToFetch {
			wg.Add(1)
			go fetchRssAndUpdateDB(feedToFetch, wg, v1)
		}
		wg.Wait()
	}
}
