package response

import (
	_ "embed"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/cmd/api/core"
	"github.com/lbenedar/fresta/internal/models/db"
)

func GetCmdData(cmd string, app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	procFunc, ok := processFuncs[cmd]
	if ok {
		return procFunc(app, i)
	}
	return nil, nil
}

func PrepareCommandResponse(cmd string, app *core.Application, i *discordgo.InteractionCreate) (*discordgo.InteractionResponseData, error) {
	data, err := GetCmdData(cmd, app, i)
	if err != nil {
		return nil, err
	}

	return GetParserResponse(cmd, data)
}

var processFuncs = map[string]func(app *core.Application, i *discordgo.InteractionCreate) (any, error){
	"help":   procHelpCmd,
	"status": procStatusCmd,
	"worlds": procWorldsCmd,
	"link":   procLinkCmd,
}

func procHelpCmd(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	data := &struct {
		Host string
	}{
		Host: fmt.Sprintf("http://%s", app.FoundryApp.Transport.Http.Host),
	}
	return data, nil
}

func procStatusCmd(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	return app.FoundryApp.Status, nil
}

func procLinkCmd(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	data := &struct {
		Host string
	}{
		Host: fmt.Sprintf("http://%s", app.FoundryApp.Transport.Http.Host),
	}
	return data, nil
}

func procWorldsCmd(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	games := make(db.Games, 0)
	err := games.Get(app.FoundryApp.Transport.DB)
	if err != nil {
		return nil, err
	}

	type WorldData struct {
		Name      string   `json:"name"`
		UserNames []string `json:"user_names"`
	}

	worldData := make([]WorldData, len(games))

	for i, game := range games {
		names := make([]string, len(game.Users))
		for i, v := range game.Users {
			names[i] = v.Name
		}

		worldData[i] = WorldData{Name: game.World.Title, UserNames: names}
	}

	app.Slogger.Info("WorldData", "data", worldData)

	return worldData, nil
}

func procLinkUsersCmd(app *core.Application, i *discordgo.InteractionCreate) (any, error) {
	games := make(db.Games, 0)
	err := games.Get(app.FoundryApp.Transport.DB)
	if err != nil {
		return nil, err
	}

	type WorldData struct {
		Name      string   `json:"name"`
		UserNames []string `json:"user_names"`
	}

	worldData := make([]WorldData, len(games))

	for i, game := range games {
		names := make([]string, len(game.Users))
		for i, v := range game.Users {
			names[i] = v.Name
		}

		worldData[i] = WorldData{Name: game.World.Title, UserNames: names}
	}

	app.Slogger.Info("WorldData", "data", worldData)

	return worldData, nil
}
