package json

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lbenedar/fresta/internal/models/db"
)

type DiscordUser discordgo.User

func (d *DiscordUser) ToDB(dest **db.DiscordUser) bool {
	if dest == nil {
		return false
	}

	discordUser := &db.DiscordUser{
		ID:            d.ID,
		Email:         d.Email,
		Username:      d.Username,
		Avatar:        d.Avatar,
		Locale:        d.Locale,
		Discriminator: d.Discriminator,
		GlobalName:    d.GlobalName,
	}

	*dest = discordUser

	return true
}
