package json

import "github.com/lbenedar/fresta/internal/models/db"

type Environment struct {
	GlobalLight           EnvironmentGlobalLight `json:"globalLight"`
	DarknessLevel         float64                `json:"darknessLevel"`
	DarknessLock          bool                   `json:"darknessLock"`
	Cycle                 bool                   `json:"cycle"`
	ScenesEnvironmentBase EnvironmentBase        `json:"base"`
	Dark                  EnvironmentBase        `json:"dark"`
}

func (e *Environment) ToDB(dest *db.Environment) bool {
	if dest == nil {
		return false
	}

	e.GlobalLight.ToDB(&dest.GlobalLight)
	dest.DarknessLevel = e.DarknessLevel
	dest.DarknessLock = e.DarknessLock
	dest.Cycle = e.Cycle
	e.ScenesEnvironmentBase.ToDB(&dest.ScenesEnvironmentBase)
	e.Dark.ToDB(&dest.Dark)

	return true
}

type EnvironmentGlobalLight struct {
	Enabled    bool                `json:"enabled"`
	Darkness   GlobalLightDarkness `json:"darkness"`
	Alpha      float64             `json:"alpha"`
	Bright     bool                `json:"bright"`
	Color      string              `json:"color"`
	Coloration float64             `json:"coloration"`
	Luminosity float64             `json:"luminosity"`
	Saturation float64             `json:"saturation"`
	Contrast   float64             `json:"contrast"`
	Shadows    float64             `json:"shadows"`
}

func (e *EnvironmentGlobalLight) ToDB(dest *db.EnvironmentGlobalLight) bool {
	if dest == nil {
		return false
	}

	dest.Enabled = e.Enabled
	e.Darkness.ToDB(&dest.Darkness)
	dest.Alpha = e.Alpha
	dest.Bright = e.Bright
	dest.Color = e.Color
	dest.Coloration = e.Coloration
	dest.Luminosity = e.Luminosity
	dest.Saturation = e.Saturation
	dest.Contrast = e.Contrast
	dest.Shadows = e.Shadows

	return true
}

type EnvironmentBase struct {
	Hue        float64 `json:"hue"`
	Intensity  float64 `json:"intensity"`
	Luminosity float64 `json:"luminosity"`
	Saturation float64 `json:"saturation"`
	Shadows    float64 `json:"shadows"`
}

func (e *EnvironmentBase) ToDB(dest *db.EnvironmentBase) bool {
	if dest == nil {
		return false
	}

	dest.Hue = e.Hue
	dest.Intensity = e.Intensity
	dest.Luminosity = e.Luminosity
	dest.Saturation = e.Saturation
	dest.Shadows = e.Shadows

	return true
}

type GlobalLightDarkness struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

func (g *GlobalLightDarkness) ToDB(dest *db.GlobalLightDarkness) bool {
	if dest == nil {
		return false
	}

	dest.Max = g.Max
	dest.Min = g.Min

	return true
}
