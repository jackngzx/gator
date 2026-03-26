package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackngzx/gator/internal/config"
	"github.com/jackngzx/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)
	newState := &State{
		Db:  dbQueries,
		Cfg: &cfg,
	}
	cmds := Commands{
		RegisteredCommand: make(map[string]func(*State, Command) error),
	}

	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerGetUsers)
	cmds.Register("agg", agg)
	cmds.Register("addfeed", middlewareLoggedIn(addfeed))
	cmds.Register("feeds", feeds)
	cmds.Register("follow", middlewareLoggedIn(follow))
	cmds.Register("following", middlewareLoggedIn(following))
	cmds.Register("unfollow", middlewareLoggedIn(unfollow))
	cmds.Register("browse", middlewareLoggedIn(browse))
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		os.Exit(1)
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.Run(newState, Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
