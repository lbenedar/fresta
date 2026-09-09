package db

type Environment struct {
	ID uint

	GlobalLight           EnvironmentGlobalLight
	DarknessLevel         float64
	DarknessLock          bool
	Cycle                 bool
	ScenesEnvironmentBase EnvironmentBase
	Dark                  EnvironmentBase
}

type EnvironmentGlobalLight struct {
	ID uint

	Enabled    bool
	Darkness   GlobalLightDarkness
	Alpha      float64
	Bright     bool
	Color      string
	Coloration float64
	Luminosity float64
	Saturation float64
	Contrast   float64
	Shadows    float64
}

type EnvironmentBase struct {
	ID uint

	Hue        float64
	Intensity  float64
	Luminosity float64
	Saturation float64
	Shadows    float64
}

type GlobalLightDarkness struct {
	ID uint

	Max int
	Min int
}
