package json

import "github.com/lbenedar/fresta/internal/models/db"

type Playlist struct {
	Name        string         `json:"name"`
	ID          string         `json:"_id"`
	Sounds      []*Sound       `json:"sounds"`
	Mode        int            `json:"mode"`
	Playing     bool           `json:"playing"`
	Fade        int            `json:"fade,omitempty"`
	Folder      string         `json:"folder"`
	Sorting     string         `json:"sorting"`
	Seed        int            `json:"seed,omitempty"`
	Sort        int            `json:"sort"`
	Ownership   map[string]int `json:"ownership,omitempty"`
	Stats       Stats          `json:"_stats"`
	Description string         `json:"description,omitempty"`
	Channel     string         `json:"channel"`
	// Flags                any                  `json:"flags,omitempty"`
}

func (p *Playlist) ToDB(dest **db.Playlist) bool {
	if dest == nil {
		return false
	}

	playlist := &db.Playlist{
		Name:        p.Name,
		ID:          p.ID,
		Mode:        p.Mode,
		Playing:     p.Playing,
		Fade:        p.Fade,
		Folder:      p.Folder,
		Sorting:     p.Sorting,
		Seed:        p.Seed,
		Sort:        p.Sort,
		Description: p.Description,
		Channel:     p.Channel,
	}

	p.Stats.ToDB(&playlist.Stats)

	OwnershipToDB(&playlist.Ownership, p.Ownership)

	CopySliceToDB(&playlist.Sounds, p.Sounds)

	*dest = playlist

	return true
}

type Sound struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	ID          string  `json:"_id"`
	Playing     bool    `json:"playing"`
	PausedTime  float64 `json:"pausedTime"`
	Repeat      bool    `json:"repeat"`
	Volume      float64 `json:"volume"`
	Fade        int     `json:"fade"`
	Sort        int     `json:"sort"`
	Channel     string  `json:"channel"`
	Description string  `json:"description,omitempty"`
	// Flags any     `json:"flags"`
}

func (s *Sound) ToDB(dest **db.Sound) bool {
	if dest == nil {
		return false
	}

	sound := &db.Sound{
		Name:        s.Name,
		Path:        s.Path,
		ID:          s.ID,
		Playing:     s.Playing,
		PausedTime:  s.PausedTime,
		Repeat:      s.Repeat,
		Volume:      s.Volume,
		Fade:        s.Fade,
		Sort:        s.Sort,
		Channel:     s.Channel,
		Description: s.Description,
	}

	*dest = sound

	return true
}
