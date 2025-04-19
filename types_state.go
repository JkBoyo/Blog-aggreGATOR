package main

import (
	"GATOR/internal/config"
	"GATOR/internal/database"
)

type state struct {
	db     *database.Queries
	config config.Config
}
