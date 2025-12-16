package db

import (
	"database/sql"
	"fmt"
)

func IsSeen(sqlDB *sql.DB, feedURL, itemGUID string) (bool, error) {
	const q = `SELECT 1 FROM seen_items WHERE feed_url = ? AND item_guid = ? LIMIT 1`
	var one int
	err := sqlDB.QueryRow(q, feedURL, itemGUID).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("is seen query: %w", err)
	}
	return true, nil
}

func MarkSeen(sqlDB *sql.DB, feedURL, itemGUID string) error {
	const q = `INSERT OR IGNORE INTO seen_items (feed_url, item_guid) VALUES (?, ?)`
	if _, err := sqlDB.Exec(q, feedURL, itemGUID); err != nil {
		return fmt.Errorf("mark seen exec: %w", err)
	}
	return nil
}
