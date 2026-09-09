package json

import "github.com/lbenedar/fresta/internal/models/db"

type Grid struct {
	Type      int     `json:"type"`
	Size      int     `json:"size,omitempty"`
	Color     string  `json:"color,omitempty"`
	Alpha     float64 `json:"alpha,omitempty"`
	Distance  int     `json:"distance"`
	Units     string  `json:"units"`
	Diagonals int     `json:"diagonals,omitempty"`
	Style     string  `json:"style,omitempty"`
	Thickness int     `json:"thickness,omitempty"`
}

func (g *Grid) ToDB(dest **db.Grid) bool {
	if dest == nil {
		return false
	}

	grid := &db.Grid{
		Type:      g.Type,
		Size:      g.Size,
		Color:     g.Color,
		Alpha:     g.Alpha,
		Distance:  g.Distance,
		Units:     g.Units,
		Diagonals: g.Diagonals,
		Style:     g.Style,
		Thickness: g.Thickness,
	}

	*dest = grid

	return true
}
