package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
		return err
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
		return err
	}

	if err := s.Cfg.SetUser(insertedUser.Name); err != nil {
		return err
	}
	fmt.Println("The user has been created")
	log.Println(insertedUser)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	if err := s.Db.ResetDatabase(ctx); err != nil {
		return err
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
	ctx := context.Background()

	rssFeed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}
