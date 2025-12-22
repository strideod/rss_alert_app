# rss_alert_app

A command-line tool for monitoring RSS feeds (such as status pages) and tracking incidents and events in a local SQLite database.

## Features

- Add, list, update, and delete RSS feed URLs.
- Fetch and parse feeds, deduplicate events, and track incident status.
- Stores all data in a local SQLite database.
- Designed for extensibility and automation.

## Requirements

- Go 1.18 or newer
- SQLite3

## Setup

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/rss_alert_app.git
   cd rss_alert_app
   ```

2. **Install dependencies:**

   ```sh
   go mod download
   ```

3. **Configure environment:**
   - Copy `.env` if needed and set `DB_PATH` (default: `./rss_alerts.db`).

4. **Build the app:**

   ```sh
   go build -o rss_alert_app
   ```

## Usage

```sh
./rss_alert_app [add|list|update|delete|fetch] [options] [url]
```

### Commands

- `add [--name NAME] URL`  
  Add a new RSS feed URL (optionally with a name).

- `list`  
  List all feed URLs.

- `update OLD_URL NEW_URL`  
  Update an existing feed URL.

- `delete URL`  
  Delete a feed URL.

- `fetch`  
  Fetch all feeds, process new events, and print incident updates.

### Example

```sh
./rss_alert_app add --name "Example Status" https://status.example.com/history.rss
./rss_alert_app list
./rss_alert_app fetch
```

## Database

- Uses SQLite for persistence.
- Schema is created automatically on first run.

## Development

- Main entry point: [`main.go`](main.go)
- Database logic: [`internal/db/`](internal/db/)
- Feed fetching and parsing: [`internal/feeds/`](internal/feeds/)

## License

MIT License. See [LICENSE](LICENSE).
