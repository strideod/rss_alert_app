package db

import (
	"database/sql"
	"fmt"
)

type Feed struct {
	ID   int
	Name string
	URL  string
}

func GetFeeds(sqlDB *sql.DB) ([]Feed, error) {
	rows, err := sqlDB.Query("SELECT id, name, url FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("query feeds: %w", err)
	}
	defer rows.Close()

	out := make([]Feed, 0, 8)
	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.ID, &f.Name, &f.URL); err != nil {
			return nil, fmt.Errorf("scan feed row: %w", err)
		}
		out = append(out, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate feeds: %w", err)
	}
	return out, nil
}