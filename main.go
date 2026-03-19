package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackngzx/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	newState := &config.State{
		Cfg: &cfg,
	}

	cmds := config.Commands{
		RegisteredCommand: make(map[string]func(*config.State, config.Command) error),
	}

	cmds.Register("login", config.HandlerLogin)
	if len(os.Args) < 2 {
		fmt.Println("Missing arguments. Exiting program...")
		os.Exit(1)
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.Run(newState, config.Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
