package json

import "github.com/lbenedar/fresta/internal/models/db"

type Token struct {
	DisplayName      int             `json:"displayName"`
	DisplayBars      int             `json:"displayBars"`
	Disposition      int             `json:"disposition"`
	Sight            TokenSight      `json:"sight"`
	Name             string          `json:"name"`
	ActorLink        bool            `json:"actorLink"`
	AppendNumber     bool            `json:"appendNumber"`
	PrependAdjective bool            `json:"prependAdjective"`
	Texture          TokenTexture    `json:"texture"`
	Width            float64         `json:"width"`
	Height           float64         `json:"height"`
	LockRotation     bool            `json:"lockRotation"`
	Rotation         int             `json:"rotation"`
	Alpha            int             `json:"alpha"`
	Bar1             TokenBar        `json:"bar1"`
	Bar2             TokenBar        `json:"bar2"`
	Light            Light           `json:"light"`
	RandomImg        bool            `json:"randomImg"`
	Occludable       TokenOccludable `json:"occludable"`
	Ring             Ring            `json:"ring"`
	TurnMarker       TokenTurnMarker `json:"turnMarker"`
	MovementAction   any             `json:"movementAction"`
	// Flags            any          `json:"flags"`
	// DetectionModes   []any        `json:"detectionModes"`
}

func (t *Token) ToDB(dest **db.Token) bool {
	if dest == nil {
		return false
	}

	token := &db.Token{
		DisplayName:      t.DisplayName,
		DisplayBars:      t.DisplayBars,
		Disposition:      t.Disposition,
		Name:             t.Name,
		ActorLink:        t.ActorLink,
		AppendNumber:     t.AppendNumber,
		PrependAdjective: t.PrependAdjective,
		Width:            t.Width,
		Height:           t.Height,
		LockRotation:     t.LockRotation,
		Rotation:         t.Rotation,
		Alpha:            t.Alpha,
		RandomImg:        t.RandomImg,
	}

	t.Ring.ToDB(&token.Ring)
	t.Sight.ToDB(&token.Sight)
	t.Texture.ToDB(&token.Texture)
	t.Bar1.ToDB(&token.Bar1)
	t.Bar2.ToDB(&token.Bar2)
	t.Light.ToDB(&token.Light)
	t.Occludable.ToDB(&token.Occludable)
	t.TurnMarker.ToDB(&token.TurnMarker)

	*dest = token

	return true
}

type TokenTexture struct {
	Src            string  `json:"src"`
	ScaleX         float64 `json:"scaleX"`
	ScaleY         float64 `json:"scaleY"`
	OffsetX        float64 `json:"offsetX"`
	OffsetY        float64 `json:"offsetY"`
	Rotation       float64 `json:"rotation"`
	AnchorX        float64 `json:"anchorX"`
	AnchorY        float64 `json:"anchorY"`
	Fit            string  `json:"fit"`
	Tint           string  `json:"tint"`
	AlphaThreshold float64 `json:"alphaThreshold"`
}

func (t *TokenTexture) ToDB(dest **db.TokenTexture) bool {
	if dest == nil {
		return false
	}

	texture := &db.TokenTexture{
		Src:            t.Src,
		ScaleX:         t.ScaleX,
		ScaleY:         t.ScaleY,
		OffsetX:        t.OffsetX,
		OffsetY:        t.OffsetY,
		Rotation:       t.Rotation,
		AnchorX:        t.AnchorX,
		AnchorY:        t.AnchorY,
		Fit:            t.Fit,
		Tint:           t.Tint,
		AlphaThreshold: t.AlphaThreshold,
	}

	*dest = texture

	return true
}

type TokenSight struct {
	Color       string
	Enabled     bool
	Range       int
	Angle       int
	VisionMode  string
	Attenuation float64
	Brightness  float64
}

func (t *TokenSight) ToDB(dest **db.TokenSight) bool {
	if dest == nil {
		return false
	}

	sight := &db.TokenSight{
		Color:       t.Color,
		Enabled:     t.Enabled,
		Range:       t.Range,
		Angle:       t.Angle,
		VisionMode:  t.VisionMode,
		Attenuation: t.Attenuation,
		Brightness:  t.Brightness,
	}

	*dest = sight

	return true
}

type TokenBar struct {
	Attribute string `json:"attribute"`
}

func (t *TokenBar) ToDB(dest *db.TokenBar) bool {
	if dest == nil {
		return false
	}

	dest.Attribute = t.Attribute

	return true
}

type TokenOccludable struct {
	Radius int `json:"radius"`
}

func (t *TokenOccludable) ToDB(dest *db.TokenOccludable) bool {
	if dest == nil {
		return false
	}

	dest.Radius = t.Radius

	return true
}

type TokenTurnMarker struct {
	Mode        int    `json:"mode"`
	Animation   string `json:"animation"`
	Src         string `json:"src"`
	Disposition bool   `json:"disposition"`
}

func (t *TokenTurnMarker) ToDB(dest *db.TokenTurnMarker) bool {
	if dest == nil {
		return false
	}

	dest.Mode = t.Mode
	dest.Animation = t.Animation
	dest.Src = t.Src
	dest.Disposition = t.Disposition

	return true
}
