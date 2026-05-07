# Gator

Another project from the incredible Boot.dev backend course. This CLI tool helps to aggregate RSS feeds, construct them into posts, while implementing some interesting features to tag posts with users. 

## Prerequisites

- [Go](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/download)

## Installation

```bash
go install github.com/your-username/gator@latest
```

## Config

Create a `.gatorconfig.json` file in your home directory:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your actual Postgres connection string.

## Usage

Register a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start aggregating (e.g., every 30 seconds):

```bash
gator agg 30s
```

Browse posts:

```bash
gator browse [limit]
```

Other commands:

- `gator login <name>` - Log in as an existing user
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow an existing feed
- `gator unfollow <url>` - Unfollow a feed
- `gator reset` - Reset the database to the clean slate

## Motivation

Overall I learn a whole bunch of new things from this project. It was indeed the most demanding project that I have experienced so far, and I had to utilize Boot AI for quite a bit during the journey. However, I gained solid knowledge of building my own database, and query into it to get the right data. I still need to learn much more in terms of organizing my code, but I am really excited and pumped after this chapter. 
