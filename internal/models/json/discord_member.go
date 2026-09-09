package json

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/internal/models/db"
)

type DiscordMember discordgo.Member

func (d *DiscordMember) ToDB(dest **db.DiscordMember) bool {
	if dest == nil {
		return false
	}

	discordMember := &db.DiscordMember{
		GuildId:  d.GuildID,
		JoinedAt: d.JoinedAt,
		Nick:     d.Nick,
		Deaf:     d.Deaf,
		Mute:     d.Mute,
	}

	if d.User != nil {
		(*DiscordUser)(d.User).ToDB(&discordMember.User)
	}

	*dest = discordMember

	return true
}
