package db

import (
	"database/sql"
	"fmt"
	"time"

	"rss_alert_app/internal/models"
)

func toRFC3339Ptr(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func InsertEventIgnore(sqlDB *sql.DB, e models.Event) error {
	const q = `
INSERT OR IGNORE INTO events (
feed_id, incident_key,
event_guid, content_hash,
title, link,
derived_status, derived_severity,
published_at, ingested_at,
raw_content
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`
	_, err := sqlDB.Exec(
		q,
		e.FeedID,
		e.IncidentKey,
		e.EventGUID,
		e.ContentHash,
		e.Title,
		e.Link,
		e.DerivedStatus,
		e.DerivedSeverity,
		toRFC3339Ptr(e.PublishedAt),
		e.IngestedAt.UTC().Format(time.RFC3339),
		e.RawContent,
	)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}