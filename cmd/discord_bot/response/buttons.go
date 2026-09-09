package response

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/cmd/api/core"
	"github.com/lbenedar/fresta/internal/models/db"
)

func GetBtnData(cmd string, app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	btnFunc, ok := btnToFuncMap[cmd]
	if ok {
		return btnFunc(app, i)
	}
	return nil, nil
}

func PrepareButtonResponse(cmd string, app *core.Application, i *discordgo.InteractionCreate) (*discordgo.InteractionResponseData, error) {
	data, err := GetBtnData(cmd, app, i)
	if err != nil {
		return nil, err
	}

	return GetParserResponse(cmd, data)
}

var btnToFuncMap = map[string]func(app *core.Application, i *discordgo.InteractionCreate) (any, error){
	"help":   procBtnHelp,
	"status": procBtnStatus,
	"worlds": procBtnWorlds,
}

func procBtnStatus(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	return app.FoundryApp.Status, nil
}

func procBtnHelp(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	data := &struct {
		Host string
	}{
		Host: fmt.Sprintf("http://%s", app.FoundryApp.Transport.Http.Host),
	}
	return data, nil
}

func procBtnWorlds(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	games := make(db.Games, 0)
	err := games.Get(app.FoundryApp.Transport.DB)
	if err != nil {
		return nil, err
	}

	type WorldData struct {
		Name      string
		UserNames []string
	}

	worldData := make([]WorldData, len(games))

	for i, game := range games {
		names := make([]string, len(game.Users))
		for i, v := range game.Users {
			names[i] = v.Name
		}

		worldData[i] = WorldData{Name: game.World.Title, UserNames: names}
	}

	return worldData, nil
}
