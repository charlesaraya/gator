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
	cmds.Register("register", commands.RegisterHandler)
	cmds.Register("users", commands.UsersHandler)
	cmds.Register("reset", commands.ResetHandler)
	cmds.Register("agg", commands.AggregateFeedHandler)
	cmds.Register("addfeed", commands.AddFeedHandler)

	var cliCommand commands.Command
	switch len(os.Args) {
	case 1:
		log.Fatal("not enough arguments were provided")
	case 2:
		cliCommand = commands.Command{
			Name:      os.Args[1],
			Arguments: []string{},
		}
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
