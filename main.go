package main

import (
	"flag"
	"fmt"
	"os"
	"rss_alert_app/internal/db"
	"rss_alert_app/internal/feeds"
	"rss_alert_app/internal/log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	log.Init()

	var dbPathFlag string

	flag.StringVar(&dbPathFlag, "db-path", "", "Path to SQLite database file")
	flag.Parse()

	dbPath, err := resolveDBPath(dbPathFlag)
	if err != nil {
		log.LogError("database path resolution failed", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.OpenDB(dbPath)
	if err != nil {
		log.LogError("failed to open database", "error", err, "dbPath", dbPath)
		os.Exit(1)
	}
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
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		nameFlag := addCmd.String("name", "", "Name of the feed")
		addCmd.Parse(os.Args[2:])
		args := addCmd.Args()
		if len(args) < 1 {
			fmt.Println("Usage: rss_alert_app add [--name NAME] [url]")
			return
		}
		url := args[0]
		name := *nameFlag
		err := db.AddFeedURL(sqlDB, name, url)
		if err == nil {
			fmt.Println("Feed URL added.")
		}
	case "list":
		urls, err := db.GetFeedURLs(sqlDB)
		if err != nil {
			fmt.Println("Failed to get feed URLs:", err)
			return
		}
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

func resolveDBPath(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	if env := os.Getenv("DB_PATH"); env != "" {
		return env, nil
	}

	return "", fmt.Errorf("DB_PATH not set  (use --db-path or environment variable)")
}
