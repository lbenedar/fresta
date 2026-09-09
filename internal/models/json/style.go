package json

import "github.com/lbenedar/fresta/internal/models/db"

type Style struct {
	Src string `json:"src"`
}

func (s *Style) ToDB(dest **db.Style) bool {
	if dest == nil {
		return false
	}

	style := &db.Style{
		Src: s.Src,
	}

	*dest = style

	return true
}
