package json

import "github.com/lbenedar/fresta/internal/models/db"

type Setting struct {
	Key   string `json:"key"`
	User  any    `json:"user"`
	Value string `json:"value"`
	ID    string `json:"_id"`
	Stats Stats  `json:"_stats"`
}

func (s *Setting) ToDB(dest **db.Setting) bool {
	if dest == nil {
		return false
	}

	setting := &db.Setting{
		Key:   s.Key,
		Value: s.Value,
		ID:    s.ID,
	}

	s.Stats.ToDB(&setting.Stats)

	*dest = setting

	return true
}
