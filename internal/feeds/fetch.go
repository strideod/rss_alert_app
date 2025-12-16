package feeds

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"database/sql"
	"strings"

	"rss_alert_app/internal/db"
	"rss_alert_app/internal/log"

	"github.com/mmcdole/gofeed"
)

func itemKey(item *gofeed.Item) string {
	if item == nil {
		return ""
	}
	if strings.TrimSpace(item.GUID) != "" {
		return strings.TrimSpace(item.GUID)
	}
	if strings.TrimSpace(item.Link) != "" {
		return strings.TrimSpace(item.Link)
	}

	// last resort: hash fields
	h := sha256.Sum256([]byte(strings.TrimSpace(item.Title) + "|" + strings.TrimSpace(item.Published) + "|" + strings.TrimSpace(item.Description)))
	return "hash:" + hex.EncodeToString(h[:])
}

func FetchFeeds(sqlDB *sql.DB) {
	// Implementation for fetching RSS feeds will go here
	feedURLs := db.GetFeedURLs(sqlDB)
	if len(feedURLs) == 0 {
		log.LogWarning("No feed URLs found in database", "URLS", feedURLs)
		return
	}
	fp := gofeed.NewParser()

	feeds := make([]struct {
		URL  string
		Name string
	}, 0, len(feedURLs))

	for _, feed := range feeds {
		parsedFeed, err := fp.ParseURL(feed.URL)
		if err != nil {
			log.LogError("Failed to parse feed", "url", feed.URL, "error", err)
			continue
		}

		for _, item := range parsedFeed.Items {
			key := itemKey(item)
			if key == "" {
				continue
			}

			seen, err := db.IsSeen(sqlDB, feed.URL, key)
			if err != nil {
				log.LogError("seen check failed", "feed_url", feed.URL, "key", key, "error", err)
				continue
			}
			if seen {
				continue
			}

			// NEW ITEM -> print and mark seen
			fmt.Printf("[%s] %s\n%s\n\n", feed.Name, item.Title, item.Link)

			if err := db.MarkSeen(sqlDB, feed.URL, key); err != nil {
				log.LogError("mark seen failed", "feed", feed.URL, "key", key, "error", err)
			}
		}
	}
}
