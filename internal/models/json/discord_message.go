package json

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/internal/models/db"
)

type DiscordMessage discordgo.Message

func (d *DiscordMessage) ToDB(dest **db.DiscordMessage) bool {
	if dest == nil {
		return false
	}

	discordMessage := &db.DiscordMessage{
		ID:              d.ID,
		GuildId:         d.GuildID,
		Content:         d.Content,
		Timestamp:       d.Timestamp,
		EditedTimestamp: d.EditedTimestamp,
	}

	*dest = discordMessage

	return true
}
