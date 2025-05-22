package main

import (
	"log"
	"os"

	"github.com/charlesaraya/gator/internal/commands"
	"github.com/charlesaraya/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("reading config file failed, %s", err.Error())
	}
	state := commands.GetState()
	state.Config = &cfg
	state.Db, err = config.LoadDB(&cfg)
	if err != nil {
		log.Fatalf("loading DB failed, %s", err.Error())
	}

	cmds := commands.GetCommands()
	cmds.Register("login", commands.LoginHandler)

	var cliCommand commands.Command
	switch len(os.Args) {
	case 1, 2:
		log.Fatal("not enough arguments were provided")
	default:
		cliCommand = commands.Command{
			Name:      os.Args[1],
			Arguments: os.Args[2:],
		}
	}
	if err = cmds.Run(&state, cliCommand); err != nil {
		log.Fatalf("running command failed, %s", err.Error())
	}
}
