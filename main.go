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
		commands: make(map[string]cmdDet),
	}

	cmds.register("login", cmdDet{
		handler:     handlerLogin,
		description: "login <username>\t\tLogs in user using username"})
	cmds.register("register", cmdDet{
		handler:     handlerRegister,
		description: "register <username>\t\tRegisters a new user"})
	cmds.register("reset", cmdDet{
		handler:     handlerReset,
		description: "reset\t\t\tresets the database to an empty slate. Only for testing purposes"})
	cmds.register("users", cmdDet{
		handler:     handlerUsers,
		description: "users\t\t\tlists registered users"})
	cmds.register("agg", cmdDet{
		handler:     handlerAgg,
		description: "agg <formattedtime>\t\ttaggregates posts from feeds every time amount passed in. ex. 1s, 1m, 1h, 1d"})
	cmds.register("feeds", cmdDet{
		handler:     handlerFeeds,
		description: "feeds\t\t\tlists feeds and user that added it"})
	cmds.register("set-db", cmdDet{
		handler:     handlerSetDb,
		description: "set-db <db-conn-str>\tupdate db connection string in config"})
	cmds.register("addfeed", cmdDet{
		handler:     middleWareLoggedIn(handlerAddFeed),
		description: "addfeed <feedname> <url>\tAdds feed at specified url wiht specified name"})
	cmds.register("follow", cmdDet{
		handler:     middleWareLoggedIn(handlerFollow),
		description: "follow <url>\t\tFollows feed at the given url"})
	cmds.register("following", cmdDet{
		handler:     middleWareLoggedIn(handlerFollowing),
		description: "following\t\t\tLists feeds the logged in user is following"})
	cmds.register("unfollow", cmdDet{
		handler:     middleWareLoggedIn(handlerUnfollow),
		description: "unfollow <feed-url>\t\tunfollows feed at provided url for logged in user"})
	cmds.register("browse", cmdDet{
		handler:     middleWareLoggedIn(handlerBrowse),
		description: "browse <limit>\t\tlist feeds by most recently published. Defaults to limit of 2"})

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
	if cmd.name == "-h" || cmd.name == "--help" {
		fmt.Println("usage GATOR [command] <args>")
		for _, val := range cmds.commands {
			fmt.Printf("    %s\n", val.description)
		}
		os.Exit(0)
	}

	err = cmds.run(&ste, cmd)
	if err != nil {
		fmt.Printf("%s error %v\n", cmd.name, err)
		os.Exit(1)
	}

	os.Exit(0)
}
