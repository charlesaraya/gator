package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charlesaraya/gator/internal/config"
	"github.com/charlesaraya/gator/internal/database"
	"github.com/charlesaraya/gator/internal/rss"
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

func UsersHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s called with arguments", cmd.Name)
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.UserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func ResetHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s called with arguments", cmd.Name)
	}
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	log.Printf("Reset")
	return nil
}

func AggregateFeedHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("%s called no arguments", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}
	log.Printf("Aggregate Feed: collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err = scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fetchedFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	if err = s.Db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}
	log.Printf("Fetched: %s (%v items)\n", feed.Name, len(fetchedFeed.Channel.Items))
	for i, item := range fetchedFeed.Channel.Items {
		fmt.Printf("\t%d. %s\n", i, item.Title)
	}
	return nil
}

func AddFeedHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("%s called with not enough arguments", cmd.Name)
	}
	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}
	feed, err := s.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	log.Printf("Add Feed: '%s' added '%s' (%s)", user.Name, feedName, feedUrl)
	return nil
}

func DeleteFeedHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("%s called with no arguments", cmd.Name)
	}
	feedUrl := cmd.Arguments[0]
	err := s.Db.DeleteFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}
	log.Printf("Delete Feed: %s", feedUrl)
	return nil
}

func FeedsHandler(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s called with arguments", cmd.Name)
	}

	feeds, err := s.Db.GetUserFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s(%s) from %s\n", feed.Name, feed.Url, feed.UserName)
	}
	log.Printf("Feeds: %v feeds", len(feeds))
	return nil
}

func FollowFeedsHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("%s called with no arguments", cmd.Name)
	}
	feedUrl := cmd.Arguments[0]
	feed, err := s.Db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feed_follow, err := s.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	log.Printf("Follow: '%s' followed '%s' feed\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
}

func FollowedFeedsHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s called with arguments", cmd.Name)
	}
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s follows %s\n", user.Name, feed.FeedName)
	}
	log.Printf("Follows: %s follows %v feeds\n", user.Name, len(feeds))
	return nil
}

func UnFollowFeedHandler(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("%s called with not enough arguments", cmd.Name)
	}
	feedUrl := cmd.Arguments[0]
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    feedUrl,
	}
	_, err := s.Db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	log.Printf("Unfollow '%s' unfollowed '%s'\n", user.Name, feedUrl)
	return nil
}
