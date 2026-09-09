package json

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/internal/models/db"
)

type DiscordInteraction discordgo.InteractionCreate

func (d *DiscordInteraction) ToDB(dest *db.DiscordInteraction) bool {
	if dest == nil {
		return false
	}

	dest.ID = d.ID
	dest.GuildId = d.GuildID
	dest.AppID = d.AppID

	if d.Message != nil {
		(*DiscordMessage)(d.Message).ToDB(&dest.Message)
	}
	switch d.Type {
	case discordgo.InteractionMessageComponent:
		dest.Message.Command = d.MessageComponentData().CustomID
	case discordgo.InteractionApplicationCommand:
		dest.Message.Command = d.ApplicationCommandData().Name
	}

	if d.Member != nil {
		(*DiscordMember)(d.Member).ToDB(&dest.Member)
	}

	if d.User != nil {
		(*DiscordUser)(d.User).ToDB(&dest.User)
	}

	return true
}
