package db

import (
	"database/sql"
	"rss_alert_app/internal/log"

	_ "github.com/mattn/go-sqlite3"

)

func SetupDB(db *sql.DB) {
	_, err := db.Exec(`
	    CREATE TABLE IF NOT EXISTS feeds (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT UNIQUE
		);
		CREATE TABLE IF NOT EXISTS seen_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			feed_url TEXT,
			item_guid TEXT,
			UNIQUE(feed_url, item_guid)
		);
	`)
	if err != nil {
		log.LogError(err.Error())
	}
}

func OpenDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.LogError(err.Error())
		return nil
	}
	return db
}

func AddFeedURL(db *sql.DB, url string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO feeds (url) VALUES (?)", url)
	if err != nil {
		log.LogError("Failed to add URL" + err.Error())
	}
	return err
}

func GetFeedURLs(db *sql.DB) ([]string) {
	rows, err := db.Query("SELECT url FROM feeds")
	if err != nil {
		log.LogError(err.Error())
		return nil
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			urls = append(urls, url)
		}
	}
	return urls
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
