package json

import "github.com/lbenedar/fresta/internal/models/db"

type Macro struct {
	Command   string         `json:"command"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Img       string         `json:"img"`
	ID        string         `json:"_id"`
	Author    string         `json:"author"`
	Scope     string         `json:"scope"`
	Folder    string         `json:"folder"`
	Sort      int            `json:"sort"`
	Ownership map[string]int `json:"ownership,omitempty"`
	Stats     Stats          `json:"_stats"`
	// Flags any `json:"flags,omitempty"`
}

func (m *Macro) ToDB(dest **db.Macro) bool {
	if dest == nil {
		return false
	}

	macro := &db.Macro{
		Command: m.Command,
		Name:    m.Name,
		Type:    m.Type,
		Img:     m.Img,
		ID:      m.ID,
		Author:  m.Author,
		Scope:   m.Scope,
		Folder:  m.Folder,
		Sort:    m.Sort,
	}

	m.Stats.ToDB(&macro.Stats)

	OwnershipToDB(&macro.Ownership, m.Ownership)

	*dest = macro

	return true
}
