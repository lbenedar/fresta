package json

import "github.com/lbenedar/fresta/internal/models/db"

type Light struct {
	Alpha          float64        `json:"alpha"`
	Angle          int            `json:"angle"`
	Bright         float64        `json:"bright"`
	Coloration     float64        `json:"coloration"`
	Dim            float64        `json:"dim"`
	Attenuation    float64        `json:"attenuation"`
	Luminosity     float64        `json:"luminosity"`
	Saturation     float64        `json:"saturation"`
	Contrast       float64        `json:"contrast"`
	Shadows        float64        `json:"shadows"`
	LightAnimation LightAnimation `json:"animation"`
	LightDarkness  LightDarkness  `json:"darkness"`
	Negative       bool           `json:"negative"`
	Priority       int            `json:"priority"`
	Color          string         `json:"color"`
}

func (l *Light) ToDB(dest **db.Light) bool {
	if dest == nil {
		return false
	}

	light := &db.Light{
		Alpha:       l.Alpha,
		Angle:       l.Angle,
		Bright:      l.Bright,
		Coloration:  l.Coloration,
		Dim:         l.Dim,
		Attenuation: l.Attenuation,
		Luminosity:  l.Luminosity,
		Saturation:  l.Saturation,
		Contrast:    l.Contrast,
		Shadows:     l.Shadows,
		Negative:    l.Negative,
		Priority:    l.Priority,
		Color:       l.Color,
	}

	l.LightAnimation.ToDB(&light.LightAnimation)
	l.LightDarkness.ToDB(&light.LightDarkness)

	*dest = light

	return true
}

type LightAnimation struct {
	Type      any  `json:"type"`
	Speed     int  `json:"speed"`
	Intensity int  `json:"intensity"`
	Reverse   bool `json:"reverse"`
}

func (l *LightAnimation) ToDB(dest *db.LightAnimation) bool {
	if dest == nil {
		return false
	}

	dest.Speed = l.Speed
	dest.Intensity = l.Intensity
	dest.Reverse = l.Reverse

	return true
}

type LightDarkness struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func (l *LightDarkness) ToDB(dest *db.LightDarkness) bool {
	if dest == nil {
		return false
	}

	dest.Min = l.Min
	dest.Max = l.Max

	return true
}
