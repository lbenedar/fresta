package types

import (
	"errors"
	"strings"
)

var (
	ErrNoAuthWorldData = errors.New("You shold pass world data for initialization")
	ErrWrongFormat     = errors.New("You have passed data in wrong format. Please check it once again")
)

type WorldData struct {
	Name     string
	Username string
	UserId   string
	UserPass string
}

func CreateWorldDataSlice(worldFlag string, userFlag string, passFlag string) ([]WorldData, error) {
	worlds := strings.Split(worldFlag, ",")
	users := strings.Split(userFlag, ",")
	passwords := strings.Split(passFlag, ",")

	if (len(worlds) == 1 && worlds[0] == "") || (len(users) == 1 && users[0] == "") {
		return nil, ErrNoAuthWorldData
	}

	worldLen := len(worlds)
	userLen := len(users)
	passLen := len(passwords)

	worldData := make([]WorldData, worldLen)
	if worldLen == userLen && userLen == passLen {
		for i := range worlds {
			worldData[i].Name = worlds[i]
			worldData[i].Username = users[i]
			worldData[i].UserPass = passwords[i]
		}
	} else if worldLen != userLen && userLen == passLen && passLen == 1 {
		for i := range worlds {
			worldData[i].Name = worlds[i]
			worldData[i].Username = users[0]
			worldData[i].UserPass = passwords[0]
		}
	} else {
		return make([]WorldData, 0), ErrWrongFormat
	}
	return worldData, nil
}
