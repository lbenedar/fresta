package json

import "github.com/lbenedar/fresta/internal/models/db"

type Combat struct {
	Id         string       `json:"_id"`
	Type       string       `json:"type"`
	Scene      string       `json:"scene"`
	Groups     []string     `json:"groups"`
	Combatants []*Combatant `json:"combatants"`
	Active     bool         `json:"active"`
	Round      int          `json:"round"`
	Turn       int          `json:"turn"`
	Sort       int          `json:"sort"`
	Stats      Stats        `json:"stats"`
	// System     any		`json:"system"`
	// Flags      any		`json:"flags"`
}

func (c *Combat) ToDB(dest **db.Combat) bool {
	if dest == nil {
		return false
	}

	combat := &db.Combat{
		ID:     c.Id,
		Type:   c.Type,
		Scene:  c.Scene,
		Active: c.Active,
		Round:  c.Round,
		Turn:   c.Turn,
		Sort:   c.Sort,
	}

	c.Stats.ToDB(&combat.Stats)

	combat.Groups = make([]string, len(c.Groups))
	copy(combat.Groups, c.Groups)
	CopySliceToDB(&combat.Combatants, c.Combatants)

	*dest = combat

	return true
}

type Combatant struct {
	TokenId    string `json:"tokenId"`
	SceneId    string `json:"sceneId"`
	ActorId    string `json:"actorId"`
	Hidden     bool   `json:"hidden"`
	Id         string `json:"_id"`
	Type       string `json:"type"`
	Img        string `json:"img"`
	Initiative int    `json:"initiative"`
	Defeated   bool   `json:"defeated"`
	Group      string `json:"group"`
	Stats      Stats  `json:"stats"`
	// System     any		`json:"system"`
	// Flags      any		`json:"flags"`
}

func (c *Combatant) ToDB(dest **db.Combatant) bool {
	if dest == nil {
		return false
	}

	combatant := &db.Combatant{
		TokenId:    c.TokenId,
		SceneId:    c.SceneId,
		ActorId:    c.ActorId,
		Hidden:     c.Hidden,
		ID:         c.Id,
		Type:       c.Type,
		Img:        c.Img,
		Initiative: c.Initiative,
		Defeated:   c.Defeated,
		Group:      c.Group,
	}

	c.Stats.ToDB(&combatant.Stats)

	*dest = combatant

	return true
}
