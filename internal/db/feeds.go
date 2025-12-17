package db

import (
	"database/sql"
	"fmt"

	"rss_alert_app/internal/models"
)

func GetFeeds(sqlDB *sql.DB) ([]models.Feed, error) {
	rows, err := sqlDB.Query("SELECT id, name, url FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("query feeds: %w", err)
	}
	defer rows.Close()

	out := make([]models.Feed, 0, 8)
	for rows.Next() {
		var f models.Feed
		if err := rows.Scan(&f.FeedID, &f.Title, &f.Link); err != nil {
			return nil, fmt.Errorf("scan feed row: %w", err)
		}
		out = append(out, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate feeds: %w", err)
	}
	return out, nil
}