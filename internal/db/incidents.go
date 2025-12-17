package db

import (
	"database/sql"
	"fmt"
	"time"

	"rss_alert_app/internal/models"
)

func toRFC3339Nullable(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func GetIncidentLastHash(sqlDB *sql.DB, feedID int, incidentKey string) (string, bool, error) {
	const q = `
SELECT last_event_hash
FROM incidents
WHERE feed_id = ? AND incident_key = ?
LIMIT 1
`
	var h sql.NullString
	err := sqlDB.QueryRow(q, feedID, incidentKey).Scan(&h)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("get incident last hash %w", err)
	}

	if !h.Valid {
		return "", true, nil
	}
	return h.String, true, nil
}

func UpsertIncident(sqlDB *sql.DB, i models.Incident) error {
	// started_at: set on insert only; on conflict keep existing.
	// resolved_at: set to value if resolved, else NULL.
	const q = `
INSERT INTO incidents (
  feed_id, incident_key,
  title, link,
  status, severity,
  started_at, updated_at, resolved_at,
  last_event_hash, last_event_guid
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(feed_id, incident_key) DO UPDATE SET
  title = excluded.title,
  link = excluded.link,
  status = excluded.status,
  severity = excluded.severity,
  updated_at = excluded.updated_at,
  resolved_at = excluded.resolved_at,
  last_event_hash = excluded.last_event_hash,
  last_event_guid = excluded.last_event_guid,
  started_at = incidents.started_at;
`
	now := i.UpdatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}
	started := i.StartedAt
	if started.IsZero() {
		// set started_at on insert to "now"; updates won’t overwrite
		started = now
	}
	_, err := sqlDB.Exec(
		q,
		i.FeedID,
		i.IncidentKey,
		i.Title,
		i.Link,
		i.Status,
		i.Severity,
		started.UTC().Format(time.RFC3339),
		now.UTC().Format(time.RFC3339),
		toRFC3339Nullable(i.ResolvedAt),
		i.LastEventHash,
		i.LastEventGUID,
	)
	if err != nil {
		return fmt.Errorf("upsert incident: %w", err)
	}
	return nil
}