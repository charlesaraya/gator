package commands

import (
	"context"

	"github.com/charlesaraya/gator/internal/database"
)

func LoggedInMiddleware(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.UserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
