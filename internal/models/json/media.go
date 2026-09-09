package json

import "github.com/lbenedar/fresta/internal/models/db"

type Media struct {
	Type    string `json:"type"`
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

func (m *Media) ToDB(dest **db.Media) bool {
	if dest == nil {
		return false
	}

	media := &db.Media{
		Type:    m.Type,
		URL:     m.URL,
		Caption: m.Caption,
	}

	*dest = media

	return true
}
