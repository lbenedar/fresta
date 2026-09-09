package json

import "github.com/lbenedar/fresta/internal/models/db"

type Note struct {
	EntryID    string      `json:"entryId"`
	PageID     string      `json:"pageId"`
	Text       string      `json:"text"`
	X          float64     `json:"x"`
	Y          float64     `json:"y"`
	Global     bool        `json:"global"`
	IconSize   int         `json:"iconSize"`
	Texture    NoteTexture `json:"texture"`
	FontFamily string      `json:"fontFamily"`
	FontSize   int         `json:"fontSize"`
	TextColor  string      `json:"textColor"`
	TextAnchor int         `json:"textAnchor"`
	ID         string      `json:"_id"`
	Elevation  int         `json:"elevation"`
	Sort       int         `json:"sort"`
	// NotesFlags any `json:"flags"`
}

func (n *Note) ToDB(dest *db.Note) bool {
	if dest == nil {
		return false
	}

	dest.EntryID = n.EntryID
	dest.PageID = n.PageID
	dest.Text = n.Text
	dest.X = n.X
	dest.Y = n.Y
	dest.Global = n.Global
	dest.IconSize = n.IconSize
	dest.FontFamily = n.FontFamily
	dest.FontSize = n.FontSize
	dest.TextColor = n.TextColor
	dest.TextAnchor = n.TextAnchor
	dest.ID = n.ID
	dest.Elevation = n.Elevation
	dest.Sort = n.Sort

	n.Texture.ToDB(&dest.Texture)

	return true
}

type NoteTexture struct {
	Tint           string  `json:"tint"`
	Src            string  `json:"src"`
	ScaleX         int     `json:"scaleX"`
	ScaleY         int     `json:"scaleY"`
	OffsetX        float64 `json:"offsetX"`
	OffsetY        float64 `json:"offsetY"`
	Rotation       int     `json:"rotation"`
	AnchorX        float64 `json:"anchorX"`
	AnchorY        float64 `json:"anchorY"`
	Fit            string  `json:"fit"`
	AlphaThreshold int     `json:"alphaThreshold"`
}

func (n *NoteTexture) ToDB(dest *db.NoteTexture) bool {
	if dest == nil {
		return false
	}

	dest.Tint = n.Tint
	dest.Src = n.Src
	dest.ScaleX = n.ScaleX
	dest.ScaleY = n.ScaleY
	dest.OffsetX = n.OffsetX
	dest.OffsetY = n.OffsetY
	dest.Rotation = n.Rotation
	dest.AnchorX = n.AnchorX
	dest.AnchorY = n.AnchorY
	dest.Fit = n.Fit
	dest.AlphaThreshold = n.AlphaThreshold

	return true
}
