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
	if len(os.Args) < 2 {
		fmt.Println("Missing arguments. Exiting program...")
		os.Exit(1)
	}
	cmds.Register("register", HandlerRegister)
	if len(os.Args) < 2 {
		fmt.Println("Missing arguments. Exiting program...")
		os.Exit(1)
	}
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerGetUsers)
	cmds.Register("agg", agg)

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.Run(newState, Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
