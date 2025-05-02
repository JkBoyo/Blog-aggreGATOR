package main

import (
	"GATOR/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No username entered")
	}
	userName := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			os.Exit(1)
		}
	}
	err = s.config.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No name entered.")
	}
	userName := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}
	_, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			fmt.Println("User already exists.")
			os.Exit(1)
		}
	}
	s.config.SetUser(userName)
	fmt.Printf("User %s was created.\n   ID: %v\n   Created At: %v\n   Updated At: %v\n",
		params.Name, params.ID, params.CreatedAt, params.UpdatedAt)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.config.CurrentUserName {
			user += " (current)"
		}
		fmt.Println("*", user)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("No time between requests passed")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	fmt.Println("Fetching feeds every ", cmd.args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerFeeds(s *state, cmd command) error {
	feed, err := s.db.FetchFeed(context.Background())
	if err != nil {
		return err
	}
	for _, item := range feed {
		fmt.Printf("Name: %s\n  URL: %s\n  User: %s\n",
			item.Name,
			item.Url,
			item.Name_2)
	}
	return nil
}
