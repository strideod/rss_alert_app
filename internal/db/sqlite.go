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

type Feed struct {
	ID   int
	Name string
	URL  string
}

func AddFeedURL(db *sql.DB, name, url string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO feeds (name, url) VALUES (?, ?)", name, url)
	if err != nil {
		log.LogError("Failed to add URL" + err.Error())
	}
	return err
}

func GetFeedURLs(sqlDB *sql.DB) ([]Feed, error){
	rows, err := sqlDB.Query("SELECT id, name, url FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Feed
	for rows.Next() {
		var f Feed
		var name sql.NullString
		if err := rows.Scan(&f.ID, &name, &f.URL); err != nil {
			return nil, fmt.Errorf("scan feed row %w",err)
		}
		if name.Valid {
			f.Name = name.String
		} else {
			f.Name = ""
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func UpdateFeedURL(db *sql.DB, oldURL, newURL string) error {
	_, err := db.Exec("UPDATE feeds SET url = ? WHERE url = ?", newURL, oldURL)
	if err != nil {
		log.LogError("Failed to update feed URL: " + err.Error())
	}
	return err
}

func DeleteFeedURL(db *sql.DB, url string) error {
    _, err := db.Exec("DELETE FROM feeds WHERE url = ?", url)
    if err != nil {
        log.LogError("Failed to delete feed URL: " + err.Error())
    }
    return err
}
