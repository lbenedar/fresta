package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var cmdToFuncMap = map[string]func(*discordgo.Session, *discordgo.MessageCreate, *DiscordBot){
	"test": test,
}

func test(s *discordgo.Session, m *discordgo.MessageCreate, bot *DiscordBot) {
	msg := fmt.Sprintf("%s %s#%s: %s", m.Author.ID, m.Author.Username, m.Author.Discriminator, m.Content)

	component := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Click Me",
				Style:    discordgo.PrimaryButton,
				CustomID: "button_clicked",
			},
		},
	}

	msgSend := &discordgo.MessageSend{
		Content:    msg,
		Components: []discordgo.MessageComponent{component},
	}
	msgRet, err := s.ChannelMessageSendComplex(m.ChannelID, msgSend)
	if err != nil {
		bot.app.Slogger.Error("Error", "error", err)
		return
	}

	bot.app.Slogger.Info("Msg Return", "msg", msgRet)
}
