# Gator CLI  
A simple RSS feed aggregator built as a command-line tool in Go! Gator allows users to collect, follow, and browse content from various RSS feeds—all from the terminal.

## Motivation  
Gator was designed to offer a lightweight, terminal-based experience for managing and consuming RSS feeds. Inspired by the need to stay up to date with multiple content sources—like blogs, news outlets, and podcasts—Gator makes it easy to track and organize feeds without a graphical UI.

This project is also a practical exercise in building CLI apps in Go, integrating PostgreSQL for persistent storage, and using tools like `sqlc` and `goose` for schema management and query generation.

## Getting Started

### Prerequisites  
- Go 1.20+
- PostgreSQL
- `sqlc` (for type-safe DB queries)  
- `goose` (for database migrations)

### Installing  
Clone the repository and build the binary:

```bash
git clone https://github.com/yourusername/gator.git
cd gator
go build -o gator ./cmd/gator
```

### Running
Run the CLI by executing:
```bash
./gator
```
Make sure your database is running and migrations are applied using goose.

## Features

- **Add Feeds**: Store RSS feeds in the PostgreSQL database.
- **Follow Feeds**: Follow any feed added by other users.
- **Aggregate Posts**: Continuously fetch updates and store new posts.
- **Browse Posts**: View summaries of posts in the terminal.
- **User Management**: Register, log in, and list users.

## Commands Reference

| Command                        | Description                                                                 |
|-------------------------------|-----------------------------------------------------------------------------|
| `login <userName>`            | Log in as a user. Sets the currently logged-in user in config.              |
| `register <userName>`         | Register a new user in the database.                                        |
| `users`                       | List all registered users, with `(current)` next to the active user.        |
| `addfeed <feedName> <feedUrl>`| Add a new RSS feed. Automatically follows it.                               |
| `delfeed <feedUrl>`           | Remove a feed from the database.                                            |
| `feeds`                       | List all feeds stored in the database.                                      |
| `follow <feedUrl>`            | Follow an existing feed.                                                    |
| `unfollow <feedUrl>`          | Unfollow a feed.                                                            |
| `following`                   | List all feeds currently followed by the user.                              |
| `agg <timeBetweenRequests>`   | Start background service that fetches RSS posts periodically.               |
| `browse [limit]`              | Browse recent posts across followed feeds, showing summaries and links.     |
| `reset`                       | Reset the database (useful for testing).                                    |


## Improvement Ideas
- Add tagging and filtering for feeds and posts.
- Enable exporting of posts or feeds to a file.
- Add authentication tokens for multi-user support across sessions.
