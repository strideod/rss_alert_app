package feeds

import (
	"fmt"
	"time"
	"database/sql"

	"rss_alert_app/internal/db"
	"rss_alert_app/internal/log"

	"github.com/mmcdole/gofeed"
)

func FetchFeeds(sqlDB *sql.DB) {
	// Implementation for fetching RSS feeds will go here
	feedURLs := db.GetFeedURLs(sqlDB)
	if len(feedURLs) == 0 {
		log.LogWarning("No feed URLs found in database", "URLS", feedURLs)
		//log.Warning("No feed URLs found in database.")
		return
	}

	fp := gofeed.NewParser()
	for _, url := range feedURLs {
		log.LogInfo("Processing feed", "url", url)
		feed, err:= fp.ParseURL(url)
		if err != nil {
			log.LogError("RSS parse error", "url", url, "error", err)
			// log.Error(err.Error())
			continue
		}

		fmt.Println("Feed Title:", feed.Title)
		now := time.Now().UTC()
		for _, item := range feed.Items {
			if item.PublishedParsed != nil {
				itemDate := item.PublishedParsed.UTC()
				if itemDate.Year() == now.Year() && itemDate.Month() == now.Month() && itemDate.Day() == now.Day() {
					fmt.Printf("Item: %s\nLink: %s\n\n", item.Title, item.Link)
				}
			}
		}
	}
}
