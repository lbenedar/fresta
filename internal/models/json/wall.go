package json

import "github.com/lbenedar/fresta/internal/models/db"

type Wall struct {
	C         []int     `json:"c"`
	Light     int       `json:"light"`
	Move      int       `json:"move"`
	Sight     int       `json:"sight"`
	Sound     int       `json:"sound"`
	Dir       int       `json:"dir"`
	Door      int       `json:"door"`
	Ds        int       `json:"ds"`
	ID        string    `json:"_id"`
	Threshold Threshold `json:"threshold"`
	Animation any       `json:"animation"`
	// WallsFlags any `json:"flags"`
}

func (w *Wall) ToDB(dest *db.Wall) bool {
	if dest == nil {
		return false
	}

	dest.Light = w.Light
	dest.Move = w.Move
	dest.Sight = w.Sight
	dest.Sound = w.Sound
	dest.Dir = w.Dir
	dest.Door = w.Door
	dest.Ds = w.Ds
	dest.ID = w.ID
	w.Threshold.ToDB(&dest.Threshold)

	dest.C = make([]int, len(w.C))
	copy(dest.C, w.C)

	return true
}

type Threshold struct {
	Light       int  `json:"light"`
	Sight       int  `json:"sight"`
	Sound       int  `json:"sound"`
	Attenuation bool `json:"attenuation"`
}

func (t *Threshold) ToDB(dest *db.Threshold) bool {
	if dest == nil {
		return false
	}

	dest.Light = t.Light
	dest.Sight = t.Sight
	dest.Sound = t.Sound
	dest.Attenuation = t.Attenuation

	return true
}
