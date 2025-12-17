package feeds

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
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
	feeds, err := db.GetFeeds(sqlDB)
	if err != nil {
		log.LogError("Failed to get feeds from database", "error", err)
		return
	}

	if len(feeds) == 0 {
		log.LogWarning("No feeds found in database", "URLS", feeds)
		return
	}

	fp := gofeed.NewParser()

	for _, feed := range feeds {
		parsedFeed, err := fp.ParseURL(feed.Link)
		if err != nil {
			log.LogError("Failed to parse feed", "name", feed.Title, "url", feed.Link, "error", err)
			continue
		}

		for _, item := range parsedFeed.Items {
			key := itemKey(item)
			if key == "" {
				continue
			}

			seen, err := db.IsSeen(sqlDB, feed.Link, key)
			if err != nil {
				log.LogError("seen check failed", "feed_url", feed.Link, "key", key, "error", err)
				continue
			}
			if seen {
				continue
			}

			// NEW ITEM -> print and mark seen
			fmt.Printf("[%s] %s\n%s\n\n", feed.Title, item.Title, item.Link)

			if err := db.MarkSeen(sqlDB, feed.Link, key); err != nil {
				log.LogError("mark seen failed", "feed_name", feed.Title, "feed_url", feed.Link, "key", key, "error", err)
			}
		}
	}
}
