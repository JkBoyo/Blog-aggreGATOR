package main

import (
	"errors"
	"fmt"
	"os"
)

type command struct {
	name string
	args []string
}

type cmdDet struct {
	handler     func(*state, command) error
	description string
}

type commands struct {
	commands map[string]cmdDet
}

func (c *commands) register(name string, cD cmdDet) {
	c.commands[name] = cD
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commands[cmd.name]
	if !ok {
		return errors.New("Command not found")
	}
	if cmd.args[0] == "-h" || cmd.args[0] == "--help" {
		fmt.Println(f.description)
		os.Exit(0)
	}
	return f.handler(s, cmd)
}
