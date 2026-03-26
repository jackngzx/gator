package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackngzx/gator/internal/config"
	"github.com/jackngzx/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCommand map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command is empty")
	}
	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if _, err := s.Db.GetUser(ctx, cmd.Args[0]); err != nil {
		return fmt.Errorf("User does not exist in the database")
	}

	fmt.Println("User has been set")
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	val, ok := c.RegisteredCommand[cmd.Name]
	if ok {
		return val(s, cmd)
	} else {
		return fmt.Errorf("Command does not exist")
	}
}

func (c *Commands) Register(name string, f func(s *State, cmd Command) error) {
	c.RegisteredCommand[name] = f
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Name is empty")
	}
	ctx := context.Background()
	queries := s.Db
	_, err := queries.GetUser(ctx, cmd.Args[0])
	if err == nil {
		fmt.Println("User already exists")
		os.Exit(1)
	}
	insertedUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Cfg.SetUser(insertedUser.Name); err != nil {
		log.Fatal(err)
	}
	fmt.Println("The user has been created")
	log.Println(insertedUser)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	if err := s.Db.ResetDatabase(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database has been successfully reset")
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	ctx := context.Background()
	queries := s.Db
	users, err := queries.GetUsers(ctx)
	if err != nil {
		return nil
	}
	for _, user := range users {
		if user == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func agg(s *State, cmd Command) error {
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func addfeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Missing arguments")
	}
	ctx := context.Background()
	queries := s.Db

	feed, err := queries.AddFeed(ctx, database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedName:  cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Feed is created. Details below:")
	fmt.Println(feed)

	feedFollows, err := queries.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Feed follow is created")
	fmt.Printf("Name of the feed: %s\n", feedFollows.FeedName)
	fmt.Printf("Name of the current user: %s\n", feedFollows.UserName)
	return nil
}

func feeds(s *State, cmd Command) error {
	ctx := context.Background()
	queries := s.Db
	feeds, err := queries.GetFeeds(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.FeedName)
		fmt.Printf("Feed url: %s\n", feed.Url)
	}
	return nil
}

func follow(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	queries := s.Db

	feed, err := queries.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	feedFollows, err := queries.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Feed follow is created")
	fmt.Printf("Name of the feed: %s\n", feedFollows.FeedName)
	fmt.Printf("Name of the current user: %s\n", feedFollows.UserName)
	return nil
}

func following(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Unexpected arg. Can only be used for the current user")
	}
	ctx := context.Background()
	queries := s.Db

	fmt.Printf("Name of the current user: %s\n", user.Name)
	feedsFollowedByUser, err := queries.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Here are the feeds followed by the user:")
	for _, feed := range feedsFollowedByUser {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		ctx := context.Background()
		user, err := s.Db.GetUser(ctx, s.Cfg.CurrentUserName)
		if err != nil {
			log.Fatal(err)
		}
		if err := handler(s, cmd, user); err != nil {
			log.Fatal(err)
		}
		return nil
	}
}

func unfollow(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	queries := s.Db
	feed, err := queries.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	if err := queries.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		log.Fatal(err)
	}
	return nil
}

func scrapeFeeds(s *State) error {
	ctx := context.Background()
	queries := s.Db
	feed, err := queries.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := queries.MarkFeedFetched(ctx, feed.ID); err != nil {
		log.Fatal(err)
	}

	data, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Fatal(err)
	}
	for _, post := range data.Channel.Item {
		timeFormats := []string{
			time.RFC3339,
			time.RFC1123,
			time.RFC1123Z,
		}
		var publishedTime time.Time
		for _, timeFormat := range timeFormats {
			publishedTime, err := time.Parse(timeFormat, post.PubDate)
			if err != nil {
				log.Printf("issue with parsing time: %v: %v", publishedTime, err)
				continue
			} else {
				break
			}
		}
		_, err = queries.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Link,
			Description: sql.NullString{String: post.Description, Valid: true},
			PublishedAt: publishedTime,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("issue: %s\n", err)
			continue
		}
	}
	return nil
}

func browse(s *State, cmd Command, user database.User) error {
	var limit int
	var err error
	if len(cmd.Args) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("limit must be an integer: %w", err)
		}
	}
	ctx := context.Background()
	queries := s.Db
	posts, err := queries.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("Url: %s\n", post.Url)
	}
	return nil
}
