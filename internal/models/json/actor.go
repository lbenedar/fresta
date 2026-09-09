package json

import "github.com/lbenedar/fresta/internal/models/db"

type Actor struct {
	PrototypeToken Token          `json:"prototypeToken"`
	Img            string         `json:"img"`
	Items          []*Item        `json:"items"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	Folder         string         `json:"folder"`
	Ownership      map[string]int `json:"ownership,omitempty"`
	Stats          Stats          `json:"_stats"`
	Sort           int            `json:"sort"`
	ID             string         `json:"_id"`
	// System         any            `json:"system,omitempty"`
	// Flags          any            `json:"flags,omitempty"`
	// ActorsEffects   []any          `json:"effects"`
}

func (a *Actor) ToDB(dest **db.Actor) bool {
	if dest == nil {
		return false
	}

	actor := &db.Actor{
		Img:    a.Img,
		Name:   a.Name,
		Type:   a.Type,
		Folder: a.Folder,
		Sort:   a.Sort,
		ID:     a.ID,
	}

	a.PrototypeToken.ToDB(&actor.PrototypeToken)
	a.Stats.ToDB(&actor.Stats)

	OwnershipToDB(&actor.Ownership, a.Ownership)

	CopySliceToDB(&actor.Items, a.Items)

	*dest = actor

	return true
}
