package main

import (
	"GATOR/internal/config"
	"GATOR/internal/database"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Read error", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(errors.New("Database opening error"))
		os.Exit(1)
	}

	dbQueries := database.New(db)

	ste := state{
		db:     dbQueries,
		config: cfg,
	}

	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	if len(os.Args) < 2 {
		fmt.Println(errors.New("Too few args to run"))
		os.Exit(1)
	}

	cmdName := os.Args[1]
	var cmdArgs []string
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	} else {
		cmdArgs = []string{}
	}

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(&ste, cmd)
	if err != nil {
		fmt.Printf("%s error %v\n", cmd.name, err)
		os.Exit(1)
	}

	os.Exit(0)
}
