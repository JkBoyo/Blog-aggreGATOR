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

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

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
