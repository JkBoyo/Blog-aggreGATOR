package main

import (
	"GATOR/internal/database"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func middleWareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Not enough arguments")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	newFeed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Println(newFeed)
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	url := cmd.args[0]
	feed, err := s.db.GetFeed(context.Background(), url)
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
	feed_follow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Println("feed: ", feed_follow.FeedName)
	fmt.Println("user: ", feed_follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	fmt.Println(user.Name, "'s feeds:")
	for _, feed := range userFeeds {
		fmt.Println("  ", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	unfollowParams := database.RemoveFeedFollowParams{
		Name: user.Name,
		Url:  cmd.args[0],
	}
	err := s.db.RemoveFeedFollow(context.Background(), unfollowParams)
	if err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	} else {
		limit = 2
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	postsToBrowse, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	for _, post := range postsToBrowse {
		dateStr := post.PublishedAt.Time.Format("Jan, 02, 2006")
		fmt.Println("Title: ", post.Title)
		fmt.Println("    Description: ", post.Description.String)
		fmt.Println("    Published: ", dateStr)
		fmt.Println("    Link: ", post.Url)
		fmt.Println()
	}
	return nil
}
