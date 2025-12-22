package feeds

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html"
	"strings"
	"time"

	"rss_alert_app/internal/db"
	"rss_alert_app/internal/log"
	"rss_alert_app/internal/models"

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
	feeds, err := db.GetFeedURLs(sqlDB)
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
		parsedFeed, err := fp.ParseURL(feed.URL)
		if err != nil {
			log.LogError("Failed to parse feed", "name", feed.Name, "url", feed.URL, "error", err)
		continue
		}

		for _, item := range parsedFeed.Items {
			eventGUID := itemKey(item)
			if eventGUID == "" {
				continue
			}

			eventKey, ok := EventKeyFromLink(item.Link)
			if !ok {
				// fallback grouping if not statuspage-like; keep it simple for now
				eventKey = strings.ToLower(strings.TrimSpace(item.Title))
			}

			derivedStatus := LatestStatusFromDescription(item.Description)
			contentHash := ContentHash(item.Description)

			// Trigger A: event GUID not seen
			seen, err := db.IsSeen(sqlDB, feed.URL, eventGUID)
			if err != nil {
				log.LogError("seen check failed", "feed_url", feed.URL, "event_guid", eventGUID, "error", err)
				continue
			}

			// Trigger B: content changed for same incident (Statuspage edits)
			lastHash, exists, err := db.GetIncidentLastHash(sqlDB, feed.ID, eventKey)
			if err != nil {
				log.LogError("failed to get incident last hash", "feed_id", feed.ID, "incident_key", eventKey, "error", err)
				continue
			}
			changed := !exists || (lastHash != "" && lastHash != contentHash)

			// If it's seen AND nothing changed, skip.
			if seen && !changed {
				continue
			}

			// Persist event (audit log)
			ev := models.Event {
				FeedID:          feed.ID,
				EventKey:        eventKey,
				EventGUID:       eventGUID,
				ContentHash:     contentHash,
				Title:           item.Title,
				Link:            item.Link,
				DerivedStatus:   derivedStatus,
				DerivedSeverity: DeriveSeverity(item.Title, item.Description), // will implement
				PublishedAt:     item.PublishedParsed,
				IngestedAt:      time.Now().UTC(),
				RawContent:      html.UnescapeString(item.Description),
			}

			if err := db.InsertEventIgnore(sqlDB, ev); err != nil {
				log.LogError("failed to insert event", "feed_id", feed.ID, "event_key", eventKey, "event_guid", eventGUID, "error", err)
				continue
			}

			// Upsert incident (board state)
			now := time.Now().UTC()
			inc := models.Incident{
				FeedID:        feed.ID,
				EventKey:   eventKey,
				Title:         item.Title,
				Link:          item.Link,
				Status:        derivedStatus,
				Severity:      ev.DerivedSeverity,
				UpdatedAt:     now,
				LastEventHash: contentHash,
				LastEventGUID: eventGUID,
			}

			// started_at should be set only on first insert; do that in Upsert SQL
			// resolved_at should be set/cleared based on status:
			if derivedStatus == "resolved" {
				inc.ResolvedAt = &now
			} else {
				inc.ResolvedAt = nil
			}

			if err := db.UpsertIncident(sqlDB, inc); err != nil {
				log.LogError("failed to upsert incident", "feed_id", feed.ID, "event_key", eventKey, "error", err)
				continue
			}

			// Mark GUID seen (transport dedupe)
			if !seen {
				if err := db.MarkSeen(sqlDB, feed.URL, eventGUID); err != nil {
					log.LogError("failed to mark seen", "feed_url", feed.URL, "event_guid", eventGUID, "error", err)
					continue
				}
			}

			// Emit to console (CLI UX)
			fmt.Printf("[%s] (%s) %s\n%s\n\n", feed.Name, derivedStatus, item.Title, item.Link)
		}
	}
}
