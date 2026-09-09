package db

import (
	"database/sql"
	"time"
)

type FoundryDataType string

const (
	SetupType = FoundryDataType("setup")
	GameType  = FoundryDataType("game")
)

type FoundryData struct {
	ID uint

	Type      FoundryDataType
	Setup     *Setup
	Game      *Game
	CreatedAt time.Time
}

type FoundryDataModel struct {
	DB *sql.DB
}
