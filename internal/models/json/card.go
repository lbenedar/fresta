package json

import "github.com/lbenedar/fresta/internal/models/db"

type CardDeck struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Description  string         `json:"description"`
	Img          string         `json:"img"`
	Cards        []*Card        `json:"cards"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	Rotation     int            `json:"rotation"`
	DisplayCount bool           `json:"displayCount"`
	Stats        Stats          `json:"_stats"`
	Ownership    map[string]int `json:"ownership,omitempty"`
	Folder       string         `json:"folder"`
	Sort         int            `json:"sort"`
	ID           string         `json:"_id"`
	// Flags           any             `json:"flags"`
	// Cards0System    any             `json:"system"`
}

func (c *CardDeck) ToDB(dest **db.CardDeck) bool {
	if dest == nil {
		return false
	}

	cardDeck := &db.CardDeck{
		Name:         c.Name,
		Type:         c.Type,
		Description:  c.Description,
		Img:          c.Img,
		Width:        c.Width,
		Height:       c.Height,
		Rotation:     c.Rotation,
		DisplayCount: c.DisplayCount,
		Folder:       c.Folder,
		Sort:         c.Sort,
		ID:           c.ID,
	}

	c.Stats.ToDB(&cardDeck.Stats)

	OwnershipToDB(&cardDeck.Ownership, c.Ownership)
	CopySliceToDB(&cardDeck.Cards, c.Cards)

	*dest = cardDeck

	return true
}

type Card struct {
	Name        string  `json:"name"`
	Faces       []*Face `json:"faces"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Rotation    int     `json:"rotation"`
	Type        string  `json:"type"`
	Value       int     `json:"value"`
	Suit        string  `json:"suit"`
	Description string  `json:"description"`
	Face        int     `json:"face"`
	Drawn       bool    `json:"drawn"`
	Sort        int     `json:"sort"`
	Back        Back    `json:"back"`
	Origin      string  `json:"origin"`
	ID          string  `json:"_id"`
	Stats       Stats   `json:"_stats"`
	// Flags       any    `json:"flags"`
	// System      any    `json:"system"`
}

func (c *Card) ToDB(dest **db.Card) bool {
	if dest == nil {
		return false
	}

	card := &db.Card{
		Name:        c.Name,
		Width:       c.Width,
		Height:      c.Height,
		Rotation:    c.Rotation,
		Type:        c.Type,
		Value:       c.Value,
		Suit:        c.Suit,
		Description: c.Description,
		Face:        c.Face,
		Drawn:       c.Drawn,
		Origin:      c.Origin,
		ID:          c.ID,
		Sort:        c.Sort,
	}

	c.Back.ToDB(&card.Back)
	c.Stats.ToDB(&card.Stats)

	CopySliceToDB(&card.Faces, c.Faces)

	*dest = card

	return true
}

type Face struct {
	Name string `json:"name"`
	Img  string `json:"img"`
	Text string `json:"text"`
}

func (f *Face) ToDB(dest *db.Face) bool {
	if dest == nil {
		return false
	}

	dest.Name = f.Name
	dest.Img = f.Img
	dest.Text = f.Text

	return true
}

type Back struct {
	Img  any    `json:"img"`
	Name string `json:"name"`
	Text string `json:"text"`
}

func (b *Back) ToDB(dest *db.Back) bool {
	if dest == nil {
		return false
	}

	dest.Name = b.Name
	dest.Text = b.Text

	return true
}
