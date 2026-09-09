package json

import "github.com/lbenedar/fresta/internal/models/db"

type Ownership struct {
	Player    string `json:"PLAYER,omitempty"`
	Trusted   string `json:"TRUSTED,omitempty"`
	Assistant string `json:"ASSISTANT,omitempty"`
}

func (o *Ownership) ToDB(dest *db.Ownership) {
	dest.Player = o.Player
	dest.Trusted = o.Trusted
	dest.Assistant = o.Assistant
}
