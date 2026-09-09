package json

import "github.com/lbenedar/fresta/internal/models/db"

type JournalPage struct {
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Text      PageText       `json:"text,omitempty"`
	ID        string         `json:"_id"`
	Title     PageTitle      `json:"title"`
	Video     PageVideo      `json:"video"`
	Src       string         `json:"src"`
	Sort      int            `json:"sort"`
	Ownership map[string]int `json:"ownership"`
	Stats     Stats          `json:"_stats,omitempty"`
	// System    PagesSystem    `json:"system"`
	// Image     any            `json:"image"`
	// Flags     any     		`json:"flags"`
	// Category  any            `json:"category"`
}

func (j *JournalPage) ToDB(dest **db.JournalPage) bool {
	if dest == nil {
		return false
	}

	journalPage := &db.JournalPage{
		Name: j.Name,
		Type: j.Type,
		ID:   j.ID,
		Src:  j.Src,
		Sort: j.Sort,
	}

	j.Text.ToDB(&journalPage.Text)
	j.Title.ToDB(&journalPage.Title)
	j.Video.ToDB(&journalPage.Video)
	j.Stats.ToDB(&journalPage.Stats)

	OwnershipToDB(&journalPage.Ownership, j.Ownership)

	*dest = journalPage

	return true
}

type PageText struct {
	Content  string `json:"content"`
	Format   int    `json:"format"`
	Markdown string `json:"markdown,omitempty"`
}

func (p *PageText) ToDB(dest *db.PageText) bool {
	if dest == nil {
		return false
	}

	dest.Content = p.Content
	dest.Format = p.Format
	dest.Markdown = p.Markdown

	return true
}

type PageTitle struct {
	Show  bool `json:"show"`
	Level int  `json:"level"`
}

func (p *PageTitle) ToDB(dest *db.PageTitle) bool {
	if dest == nil {
		return false
	}

	dest.Show = p.Show
	dest.Level = p.Level

	return true
}

type PageVideo struct {
	Controls bool    `json:"controls"`
	Volume   float64 `json:"volume"`
}

func (p *PageVideo) ToDB(dest *db.PageVideo) bool {
	if dest == nil {
		return false
	}

	dest.Controls = p.Controls
	dest.Volume = p.Volume

	return true
}
