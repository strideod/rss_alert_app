package db

import (
	"database/sql"
	"fmt"

	"rss_alert_app/internal/log"
	"rss_alert_app/internal/models"
)

func GetFeedURLs(sqlDB *sql.DB) ([]models.Feed, error) {
	rows, err := sqlDB.Query("SELECT id, name, url FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Feed
	for rows.Next() {
		var f models.Feed
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

func AddFeedURL(db *sql.DB, name, url string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO feeds (name, url) VALUES (?, ?)", name, url)
	if err != nil {
		log.LogError("Failed to add URL" + err.Error())
	}
	return err
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