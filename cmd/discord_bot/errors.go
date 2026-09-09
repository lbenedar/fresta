package main

import "errors"

var (
	ErrorTokenNotExist        = errors.New("Token of discord bot does not passed to application")
	ErrorGuildIdNotExist      = errors.New("Guild id for discord bot does not passed to application")
	ErrorWrongInteractionData = errors.New("Cant process data in interaction structure")
)
