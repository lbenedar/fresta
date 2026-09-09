package json

import "github.com/lbenedar/fresta/internal/models/db"

type User struct {
	Name       string         `json:"name"`
	Role       int            `json:"role"`
	ID         string         `json:"_id"`
	Avatar     string         `json:"avatar"`
	Character  string         `json:"character"`
	Color      string         `json:"color"`
	Pronouns   string         `json:"pronouns"`
	Hotbar     map[int]string `json:"hotbar,omitempty"`
	UsersStats Stats          `json:"_stats"`
	// Permissions any    `json:"permissions"`
	// UsersFlags   any   `json:"flags,omitempty"`
}

func (u *User) ToDB(dest **db.User) bool {
	if dest == nil {
		return false
	}

	user := &db.User{
		Name:      u.Name,
		Role:      u.Role,
		ID:        u.ID,
		Avatar:    u.Avatar,
		Character: u.Character,
		Color:     u.Color,
		Pronouns:  u.Pronouns,
	}

	u.UsersStats.ToDB(&user.Stats)

	HotbarToDB(&user.Hotbar, u.Hotbar)

	*dest = user

	return true
}
