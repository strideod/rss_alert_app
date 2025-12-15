package main

import (
	"fmt"
	"os"
	"rss_alert_app/internal/db"
	"rss_alert_app/internal/feeds"
	"rss_alert_app/internal/log"
)

func main() {
	log.Init()
	sqlDB, err := db.OpenDB(dbPath)
	if err != nil {
		log.LogError("failed to open database", "error", err, "dbPath", dbPath)
		os.Exit(1)
	defer func () {
		if err := sqlDB.Close(); err != nil {
			log.LogError("failed to close database", "error", err)
		}
	}()
	db.SetupDB(sqlDB)

	if len(os.Args) < 2 {
		fmt.Println("Usage: rss_alert_app [add|list|update|delete|fetch] [url]")
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rss_alert_app add [url]")
			return
		}
		err := db.AddFeedURL(sqlDB, os.Args[2])
		if err == nil {
			fmt.Println("Feed URL added.")
		}
	case "list":
		urls := db.GetFeedURLs(sqlDB)
		if len(urls) == 0 {
			fmt.Println("No feed URLs found.")
		}
		for _, url := range urls {
			fmt.Println(url)
		}
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: rss_alert_app update [old_url] [new_url]")
			return
		}
		err := db.UpdateFeedURL(sqlDB, os.Args[2], os.Args[3])
		if err == nil {
			fmt.Println("Feed URL updated.")
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: rss_alert_app delete [url]")
			return
		}
		err := db.DeleteFeedURL(sqlDB, os.Args[2])
		if err == nil {
			fmt.Println("Feed URL deleted.")
		}
	case "fetch":
		feeds.FetchFeeds(sqlDB)
	default:
		fmt.Println("Unknown command. Usage: rss_alert_app [add|list|update|delete|fetch] [url]")
	}
}
