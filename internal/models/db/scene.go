package db

type Scene struct {
	ID string

	Folder              string
	Name                string
	Active              bool
	Navigation          bool
	NavOrder            int
	NavName             string
	Background          SceneBackground
	Foreground          string
	ForegroundElevation int
	Thumb               string
	Width               int
	Height              int
	Padding             float64
	Initial             SceneInitial
	BackgroundColor     string
	Grid                SceneGrid
	TokenVision         bool
	Drawings            []SceneDrawing
	Tokens              []*Token
	Lights              []SceneLight
	Notes               []Note
	Sounds              []ScenesSound
	Walls               []Wall
	Playlist            string
	PlaylistSound       string
	Journal             string
	JournalEntryPage    string
	Weather             string
	Sort                int
	Ownership           []OwnershipString
	Stats               Stats
	Fog                 SceneFog
	Environment         Environment
}

type SceneBackground struct {
	ID uint

	Src            string
	ScaleX         float64
	ScaleY         float64
	OffsetX        float64
	OffsetY        float64
	Rotation       int
	Tint           string
	AnchorX        float64
	AnchorY        float64
	Fit            string
	AlphaThreshold int
}

type SceneInitial struct {
	ID uint

	X     float64
	Y     float64
	Scale float64
}

type SceneGrid struct {
	ID uint

	Type      int
	Size      int
	Color     string
	Alpha     float64
	Distance  int
	Units     string
	Style     string
	Thickness int
}

type SceneDrawing struct {
	ID string

	Author       string
	Shape        SceneDrawingShape
	X            float64
	Y            float64
	Rotation     int
	BezierFactor int
	FillType     int
	FillColor    string
	FillAlpha    float64
	StrokeWidth  int
	StrokeColor  string
	StrokeAlpha  int
	Texture      string
	Text         string
	FontFamily   string
	FontSize     int
	TextColor    string
	TextAlpha    int
	Hidden       bool
	Locked       bool
	Interface    bool
	Elevation    int
	Sort         int
}

type SceneDrawingShape struct {
	ID uint

	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type SceneLight struct {
	ID string

	X         float64
	Y         float64
	Rotation  int
	Walls     bool
	Vision    bool
	Config    *Light
	Hidden    bool
	Elevation int
}

type ScenesSound struct {
	ID string

	Path      string
	X         float64
	Y         float64
	Radius    float64
	Easing    bool
	Walls     bool
	Volume    float64
	Darkness  ScenesSoundsDarkness
	Repeat    bool
	Hidden    bool
	Elevation float64
	Effects   ScenesSoundsEffects
}

type ScenesSoundsDarkness struct {
	ID uint

	Min int
	Max int
}

type ScenesSoundsEffects struct {
	ID uint

	Base    ScenesSoundsEffectsBase
	Muffled ScenesSoundsEffectsBase
}

type ScenesSoundsEffectsBase struct {
	ID uint

	Intensity int
}

type SceneFog struct {
	ID uint

	Exploration bool
	Reset       int64
	Overlay     string
	Colors      SceneFogColors
}

type SceneFogColors struct {
	ID uint

	Explored   string
	Unexplored string
}
