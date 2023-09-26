package v1handlers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AtinAgnihotri/glog-aggregator/internal/database"
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
	fmt.Println("Fetched RSS", rss)

	updatedFeed, err := v1.DB.MarkFeedFetched(context.TODO(), feed.ID)
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
