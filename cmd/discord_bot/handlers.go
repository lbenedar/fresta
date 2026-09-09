package main

import (
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/cmd/api/core"
	"github.com/lbenedar/fresta/cmd/discord_bot/response"
	"github.com/lbenedar/fresta/internal/models/db"
	"github.com/lbenedar/fresta/internal/models/json"
)

func setupHandlers(b *DiscordBot) {
	b.dgSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.app.Slogger.Info("Logged in", "user", s.State.User.Username, "descriminator", s.State.User.Discriminator)
	})
	b.dgSession.AddHandler(b.interactionCreateHandler)
}

func InsertInteractionToDB(i *discordgo.InteractionCreate, app *core.Application) error {
	discordInteraction := &db.DiscordInteraction{}
	ok := (*json.DiscordInteraction)(i).ToDB(discordInteraction)

	if !ok {
		return ErrorWrongInteractionData
	}

	err := discordInteraction.Insert(app.FoundryApp.Transport.DB)
	if err != nil {
		return err
	}
	return nil
}

func AddBackButton(resp *discordgo.InteractionResponseData, msgID string, btnId string, b *DiscordBot) error {
	firstCommand, err := db.GetFirstCommand(b.app.FoundryApp.Transport.DB, msgID)
	if err != nil {
		return err
	}

	if btnId != firstCommand.Command {
		prevCommand, err := db.GetLastCommand(b.app.FoundryApp.Transport.DB, msgID)
		if err != nil {
			return err
		}

		type ActionRowData struct {
			Row   *discordgo.ActionsRow
			Index int
		}

		actionRows := make([]ActionRowData, 0)

		for i, comp := range resp.Components {
			if comp.Type() == discordgo.ActionsRowComponent {
				actionRows = append(actionRows, ActionRowData{Row: comp.(*discordgo.ActionsRow), Index: i})
			}
		}

		backBtn := discordgo.Button{
			Label:    "Back",
			Style:    2,
			Emoji:    &discordgo.ComponentEmoji{Name: "⬅"},
			CustomID: prevCommand.Command,
		}

		homeBtn := discordgo.Button{
			Label:    "Home",
			Style:    2,
			Emoji:    &discordgo.ComponentEmoji{Name: "🏠"},
			CustomID: firstCommand.Command,
		}

		var moveComponents []discordgo.MessageComponent
		if prevCommand.Command == firstCommand.Command {
			moveComponents = make([]discordgo.MessageComponent, 1)
			moveComponents[0] = backBtn
		} else {
			moveComponents = make([]discordgo.MessageComponent, 2)
			moveComponents[0] = backBtn
			moveComponents[1] = homeBtn
		}

		for _, actionRow := range actionRows {
			resp.Components = slices.Delete(resp.Components, actionRow.Index, actionRow.Index+1)

			if moveComponents == nil {
				break
			}

			comps := &actionRow.Row.Components

			newComps := make([]discordgo.MessageComponent, len(*comps)+len(moveComponents))
			copy(newComps, moveComponents)

			moveLen := len(moveComponents)
			copy(newComps[moveLen:], *comps)

			if len(newComps) > 5 {
				moveComponents = newComps[5:]
				newComps = newComps[:5]
			} else {
				moveComponents = nil
			}
			*comps = newComps
		}

		if len(actionRows) == 0 {
			actionRows = append(actionRows, ActionRowData{Row: &discordgo.ActionsRow{
				Components: moveComponents,
			}})
		}

		for _, actionRow := range actionRows {
			resp.Components = append(resp.Components, actionRow.Row)
		}
	}
	return nil
}

func (b *DiscordBot) interactionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		msgComponentData := i.MessageComponentData()

		resp, err := response.PrepareButtonResponse(msgComponentData.CustomID, b.app, i)
		if err != nil {
			b.app.Slogger.Error("Error", "error", err)
			return
		}

		err = AddBackButton(resp, i.Message.ID, msgComponentData.CustomID, b)
		if err != nil {
			b.app.Slogger.Error("Error", "error", err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: resp,
		})
		if err != nil {
			b.app.Slogger.Error("Error", "error", err)
			return
		}
	case discordgo.InteractionApplicationCommand:
		appCommandData := i.ApplicationCommandData()
		command := appCommandData.Name

		b.app.Slogger.Info("TargetID", "id", appCommandData.Name)

		resp, err := response.PrepareCommandResponse(command, b.app, i)
		if err != nil {
			b.app.Slogger.Error("Error", "error", err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: resp,
		})
		if err != nil {
			b.app.Slogger.Error("Error", "error", err)
			return
		}

		msg, err := s.InteractionResponse(i.Interaction)
		if err != nil {
			b.app.Slogger.Error("Error", "error", err)
			return
		}

		i.Message = msg
	}

	err := InsertInteractionToDB(i, b.app)
	if err != nil {
		b.app.Slogger.Error("Error", "error", err)
		return
	}
}
