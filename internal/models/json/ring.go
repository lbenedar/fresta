package json

import "github.com/lbenedar/fresta/internal/models/db"

type Ring struct {
	Enabled    bool       `json:"enabled"`
	RingColors RingColors `json:"colors"`
	Effects    int        `json:"effects"`
	Subject    Subject    `json:"subject"`
}

func (r *Ring) ToDB(dest **db.Ring) bool {
	if dest == nil {
		return false
	}

	ring := &db.Ring{
		Enabled: r.Enabled,
		Effects: r.Effects,
	}

	r.RingColors.ToDB(&ring.RingColors)
	r.Subject.ToDB(&ring.Subject)

	*dest = ring

	return true
}

type RingColors struct {
	Ring       string `json:"ring"`
	Background string `json:"background"`
}

func (r *RingColors) ToDB(dest *db.RingColors) bool {
	if dest == nil {
		return false
	}

	dest.Ring = r.Ring
	dest.Background = r.Background

	return true
}

type Subject struct {
	Scale   int    `json:"scale"`
	Texture string `json:"texture"`
}

func (s *Subject) ToDB(dest *db.Subject) bool {
	if dest == nil {
		return false
	}

	dest.Scale = s.Scale
	dest.Texture = s.Texture

	return true
}
