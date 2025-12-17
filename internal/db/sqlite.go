package db

import (
	"database/sql"
	"fmt"
	"rss_alert_app/internal/log"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB(db *sql.DB) {
	_, err := db.Exec(`
	    CREATE TABLE IF NOT EXISTS feeds (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS seen_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			feed_name TEXT,
			feed_url TEXT,
			item_guid TEXT,
			UNIQUE(feed_url, item_guid)
		);
		CREATE TABLE IF NOT EXISTS incidents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			
			feed_id INTEGER NOT NULL,
			incident_key TEXT NOT NULL,   -- stable per incident (e.g., status page slug)
			
			title TEXT NOT NULL,
			link TEXT,
			
			status TEXT NOT NULL CHECK (status IN ('open', 'resolved', 'unknown')),
			severity TEXT NOT NULL CHECK (severity IN ('info', 'degraded', 'outage', 'maintenance')),

			started_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			resolved_at TEXT,

			last_event_hash TEXT,
			last_event_guid TEXT,

			UNIQUE(feed_id, incident_key),
			FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			
			feed_id INTEGER NOT NULL,
			incident_key TEXT NOT NULL,
			
			event_guid TEXT NOT NULL UNIQUE,
			content_hash TEXT NOT NULL,
			
			title TEXT,
			link TEXT,
			
			derived_status TEXT NOT NULL,
			derived_severity TEXT NOT NULL,
			
			published_at TEXT,
			ingested_at TEXT NOT NULL,
			
			raw_content TEXT,
			
			FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
			);
			CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents(status);
			CREATE INDEX IF NOT EXISTS idx_events_incident ON events(feed_id, incident_key);
			CREATE INDEX IF NOT EXISTS idx_events_ingested ON events(ingested_at);
	`)
	if err != nil {
		log.LogError(err.Error())
	}
}

func OpenDB(path string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	// Validate the connection early
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping sqlite db: %w", err)
	}

	return sqlDB, nil
}


