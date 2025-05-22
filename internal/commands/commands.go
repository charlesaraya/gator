package commands

import (
	"fmt"
	"log"

	"github.com/charlesaraya/gator/internal/config"
	"github.com/charlesaraya/gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	CommandRegistry map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if err := c.CommandRegistry[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	c.CommandRegistry[name] = f
	return nil
}

func GetState() State {
	return State{}
}

func GetCommands() Commands {
	return Commands{
		CommandRegistry: make(map[string]func(*State, Command) error),
	}
}

func LoginHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("%s called with no arguments", cmd.Name)
	}
	userName := cmd.Arguments[0]

	if err := s.Config.SetUser(userName); err != nil {
		return err
	}
	log.Printf("Login: %s", userName)
	return nil
}
