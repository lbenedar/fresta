package db

type Note struct {
	ID string

	EntryID    string
	PageID     string
	Text       string
	X          float64
	Y          float64
	Global     bool
	IconSize   int
	Texture    NoteTexture
	FontFamily string
	FontSize   int
	TextColor  string
	TextAnchor int
	Elevation  int
	Sort       int
}

type NoteTexture struct {
	ID uint

	Tint           string
	Src            string
	ScaleX         int
	ScaleY         int
	OffsetX        float64
	OffsetY        float64
	Rotation       int
	AnchorX        float64
	AnchorY        float64
	Fit            string
	AlphaThreshold int
}
