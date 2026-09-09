package json

import (
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Scene struct {
	Folder              string          `json:"folder"`
	Name                string          `json:"name"`
	Active              bool            `json:"active"`
	Navigation          bool            `json:"navigation"`
	NavOrder            int             `json:"navOrder"`
	NavName             string          `json:"navName"`
	Background          SceneBackground `json:"background"`
	Foreground          string          `json:"foreground"`
	ForegroundElevation int             `json:"foregroundElevation"`
	Thumb               string          `json:"thumb"`
	Width               int             `json:"width"`
	Height              int             `json:"height"`
	Padding             float64         `json:"padding"`
	Initial             SceneInitial    `json:"initial"`
	BackgroundColor     string          `json:"backgroundColor"`
	Grid                SceneGrid       `json:"grid"`
	TokenVision         bool            `json:"tokenVision"`
	Drawings            []*SceneDrawing `json:"drawings"`
	Tokens              []*Token        `json:"tokens"`
	Lights              []*SceneLight   `json:"lights"`
	Notes               []*Note         `json:"notes"`
	Sounds              []*ScenesSound  `json:"sounds"`
	Walls               []*Wall         `json:"walls"`
	Playlist            string          `json:"playlist"`
	PlaylistSound       string          `json:"playlistSound"`
	Journal             string          `json:"journal"`
	JournalEntryPage    string          `json:"journalEntryPage"`
	Weather             string          `json:"weather"`
	Sort                int             `json:"sort"`
	Ownership           map[string]int  `json:"ownership,omitempty"`
	Stats               Stats           `json:"_stats"`
	ID                  string          `json:"_id"`
	Fog                 SceneFog        `json:"fog,omitempty"`
	Environment         Environment     `json:"environment"`

	// ScenesTemplates     []any             `json:"templates"`
	// Tiles               []any             `json:"tiles"`
	// ScenesFlags       any       `json:"flags,omitempty"`
	// Regions           []any             `json:"regions"`
}

func (s *Scene) ToDB(dest *db.Scene) bool {
	if dest == nil {
		return false
	}

	dest.Folder = s.Folder
	dest.Name = s.Name
	dest.Active = s.Active
	dest.Navigation = s.Navigation
	dest.NavOrder = s.NavOrder
	dest.NavName = s.NavName
	dest.Foreground = s.Foreground
	dest.ForegroundElevation = s.ForegroundElevation
	dest.Thumb = s.Thumb
	dest.Width = s.Width
	dest.Height = s.Height
	dest.Padding = s.Padding
	dest.BackgroundColor = s.BackgroundColor
	dest.TokenVision = s.TokenVision
	dest.Playlist = s.Playlist
	dest.PlaylistSound = s.PlaylistSound
	dest.Journal = s.Journal
	dest.JournalEntryPage = s.JournalEntryPage
	dest.Weather = s.Weather
	dest.Sort = s.Sort
	dest.ID = s.ID

	s.Background.ToDB(&dest.Background)
	s.Initial.ToDB(&dest.Initial)
	s.Grid.ToDB(&dest.Grid)
	s.Stats.ToDB(&dest.Stats)
	s.Fog.ToDB(&dest.Fog)
	s.Environment.ToDB(&dest.Environment)

	OwnershipToDB(&dest.Ownership, s.Ownership)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &dest.Drawings, s.Drawings)
	CopySliceToDBParallel(&wg, &dest.Tokens, s.Tokens)
	CopySliceToDBParallel(&wg, &dest.Lights, s.Lights)
	CopySliceToDBParallel(&wg, &dest.Notes, s.Notes)
	CopySliceToDBParallel(&wg, &dest.Sounds, s.Sounds)
	CopySliceToDBParallel(&wg, &dest.Walls, s.Walls)
	wg.Wait()

	return true
}

type SceneBackground struct {
	Src            string  `json:"src"`
	ScaleX         float64 `json:"scaleX"`
	ScaleY         float64 `json:"scaleY"`
	OffsetX        float64 `json:"offsetX"`
	OffsetY        float64 `json:"offsetY"`
	Rotation       int     `json:"rotation"`
	Tint           string  `json:"tint"`
	AnchorX        float64 `json:"anchorX"`
	AnchorY        float64 `json:"anchorY"`
	Fit            string  `json:"fit"`
	AlphaThreshold int     `json:"alphaThreshold"`
}

func (s *SceneBackground) ToDB(dest *db.SceneBackground) bool {
	if dest == nil {
		return false
	}

	dest.Src = s.Src
	dest.ScaleX = s.ScaleX
	dest.ScaleY = s.ScaleY
	dest.OffsetX = s.OffsetX
	dest.OffsetY = s.OffsetY
	dest.Rotation = s.Rotation
	dest.Tint = s.Tint
	dest.AnchorX = s.AnchorX
	dest.AnchorX = s.AnchorY
	dest.Fit = s.Fit
	dest.AlphaThreshold = s.AlphaThreshold

	return true
}

type SceneInitial struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Scale float64 `json:"scale"`
}

func (s *SceneInitial) ToDB(dest *db.SceneInitial) bool {
	if dest == nil {
		return false
	}

	dest.X = s.X
	dest.Y = s.Y
	dest.Scale = s.Scale

	return true
}

type SceneGrid struct {
	Type      int     `json:"type"`
	Size      int     `json:"size"`
	Color     string  `json:"color"`
	Alpha     float64 `json:"alpha"`
	Distance  int     `json:"distance"`
	Units     string  `json:"units"`
	Style     string  `json:"style"`
	Thickness int     `json:"thickness"`
}

func (s *SceneGrid) ToDB(dest *db.SceneGrid) bool {
	if dest == nil {
		return false
	}

	dest.Type = s.Type
	dest.Size = s.Size
	dest.Color = s.Color
	dest.Alpha = s.Alpha
	dest.Distance = s.Distance
	dest.Units = s.Units
	dest.Style = s.Style
	dest.Thickness = s.Thickness

	return true
}

type SceneDrawing struct {
	Author       string            `json:"author"`
	Shape        SceneDrawingShape `json:"shape"`
	X            float64           `json:"x"`
	Y            float64           `json:"y"`
	Rotation     int               `json:"rotation"`
	BezierFactor int               `json:"bezierFactor"`
	FillType     int               `json:"fillType"`
	FillColor    string            `json:"fillColor"`
	FillAlpha    float64           `json:"fillAlpha"`
	StrokeWidth  int               `json:"strokeWidth"`
	StrokeColor  string            `json:"strokeColor"`
	StrokeAlpha  int               `json:"strokeAlpha"`
	Texture      string            `json:"texture"`
	Text         string            `json:"text"`
	FontFamily   string            `json:"fontFamily"`
	FontSize     int               `json:"fontSize"`
	TextColor    string            `json:"textColor"`
	TextAlpha    int               `json:"textAlpha"`
	Hidden       bool              `json:"hidden"`
	Locked       bool              `json:"locked"`
	ID           string            `json:"_id"`
	Interface    bool              `json:"interface"`
	Elevation    int               `json:"elevation"`
	Sort         int               `json:"sort"`
	// Flags        any     `json:"flags"`
}

func (s *SceneDrawing) ToDB(dest *db.SceneDrawing) bool {
	if dest == nil {
		return false
	}

	dest.Author = s.Author
	dest.X = s.X
	dest.Y = s.Y
	dest.Rotation = s.Rotation
	dest.BezierFactor = s.BezierFactor
	dest.FillType = s.FillType
	dest.FillColor = s.FillColor
	dest.FillAlpha = s.FillAlpha
	dest.StrokeWidth = s.StrokeWidth
	dest.StrokeColor = s.StrokeColor
	dest.StrokeAlpha = s.StrokeAlpha
	dest.Texture = s.Texture
	dest.Text = s.Text
	dest.FontFamily = s.FontFamily
	dest.FontSize = s.FontSize
	dest.TextColor = s.TextColor
	dest.TextAlpha = s.TextAlpha
	dest.Hidden = s.Hidden
	dest.Locked = s.Locked
	dest.ID = s.ID
	dest.Interface = s.Interface
	dest.Elevation = s.Elevation
	dest.Sort = s.Sort

	s.Shape.ToDB(&dest.Shape)

	return true
}

type SceneDrawingShape struct {
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Radius any    `json:"radius"`
	Points []any  `json:"points"`
}

func (s *SceneDrawingShape) ToDB(dest *db.SceneDrawingShape) bool {
	if dest == nil {
		return false
	}

	dest.Type = s.Type
	dest.Width = s.Width
	dest.Height = s.Height

	return true
}

type SceneLight struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Rotation  int     `json:"rotation"`
	Walls     bool    `json:"walls"`
	Vision    bool    `json:"vision"`
	Config    Light   `json:"config"`
	Hidden    bool    `json:"hidden"`
	Flags     any     `json:"flags"`
	Id        string  `json:"_id"`
	Elevation int     `json:"elevation"`
}

func (s *SceneLight) ToDB(dest *db.SceneLight) bool {
	if dest == nil {
		return false
	}

	dest.X = s.X
	dest.Y = s.Y
	dest.Rotation = s.Rotation
	dest.Walls = s.Walls
	dest.Vision = s.Vision
	dest.Hidden = s.Hidden
	dest.ID = s.Id
	dest.Elevation = s.Elevation

	s.Config.ToDB(&dest.Config)

	return true
}

type ScenesSound struct {
	Path      string               `json:"path"`
	X         float64              `json:"x"`
	Y         float64              `json:"y"`
	Radius    float64              `json:"radius"`
	Easing    bool                 `json:"easing"`
	Walls     bool                 `json:"walls"`
	Volume    float64              `json:"volume"`
	Darkness  ScenesSoundsDarkness `json:"darkness"`
	ID        string               `json:"_id"`
	Repeat    bool                 `json:"repeat"`
	Hidden    bool                 `json:"hidden"`
	Elevation float64              `json:"elevation"`
	Effects   ScenesSoundsEffects  `json:"effects"`
	// Flags               any                  `json:"flags"`
}

func (s *ScenesSound) ToDB(dest *db.ScenesSound) bool {
	if dest == nil {
		return false
	}

	dest.Path = s.Path
	dest.X = s.X
	dest.Y = s.Y
	dest.Radius = s.Radius
	dest.Easing = s.Easing
	dest.Walls = s.Walls
	dest.Volume = s.Volume
	dest.ID = s.ID
	dest.Repeat = s.Repeat
	dest.Hidden = s.Hidden
	dest.Elevation = s.Elevation

	s.Darkness.ToDB(&dest.Darkness)
	s.Effects.ToDB(&dest.Effects)

	return true
}

type ScenesSoundsDarkness struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func (s *ScenesSoundsDarkness) ToDB(dest *db.ScenesSoundsDarkness) bool {
	if dest == nil {
		return false
	}

	dest.Min = s.Min
	dest.Max = s.Max

	return true
}

type ScenesSoundsEffects struct {
	Base    ScenesSoundsEffectsBase `json:"base"`
	Muffled ScenesSoundsEffectsBase `json:"muffled"`
}

func (s *ScenesSoundsEffects) ToDB(dest *db.ScenesSoundsEffects) bool {
	if dest == nil {
		return false
	}

	s.Base.ToDB(&dest.Base)
	s.Muffled.ToDB(&dest.Muffled)

	return true
}

type ScenesSoundsEffectsBase struct {
	Intensity int `json:"intensity"`
}

func (s *ScenesSoundsEffectsBase) ToDB(dest *db.ScenesSoundsEffectsBase) bool {
	if dest == nil {
		return false
	}

	dest.Intensity = s.Intensity

	return true
}

type SceneFog struct {
	Exploration bool           `json:"exploration"`
	Reset       int64          `json:"reset"`
	Overlay     string         `json:"overlay"`
	Colors      SceneFogColors `json:"colors"`
}

func (s *SceneFog) ToDB(dest *db.SceneFog) bool {
	if dest == nil {
		return false
	}

	dest.Exploration = s.Exploration
	dest.Reset = s.Reset
	dest.Overlay = s.Overlay

	s.Colors.ToDB(&dest.Colors)

	return true
}

type SceneFogColors struct {
	Explored   string `json:"explored"`
	Unexplored string `json:"unexplored"`
}

func (s *SceneFogColors) ToDB(dest *db.SceneFogColors) bool {
	if dest == nil {
		return false
	}

	dest.Explored = s.Explored
	dest.Unexplored = s.Unexplored

	return true
}
