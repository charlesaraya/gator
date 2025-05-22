package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charlesaraya/gator/internal/config"
	"github.com/charlesaraya/gator/internal/database"
	"github.com/google/uuid"
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

	_, err := s.Db.GetUser(context.Background(), userName)
	if err != nil {
		return err
	}

	if err := s.Config.SetUser(userName); err != nil {
		return err
	}
	log.Printf("Login: %s", userName)
	return nil
}

func RegisterHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("%s called with no arguments", cmd.Name)
	}
	userName := cmd.Arguments[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}
	user, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	if err = s.Config.SetUser(userName); err != nil {
		return err
	}
	log.Printf("Register: %s(%s)", user.Name, user.ID.String())
	return nil
}
