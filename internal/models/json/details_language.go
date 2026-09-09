package json

import "github.com/lbenedar/fresta/internal/models/db"

type DetailsLanguages struct {
	Details string `json:"details"`
}

func (d *DetailsLanguages) ToDB(dest *db.DetailsLanguages) bool {
	if dest == nil {
		return false
	}

	dest.Details = d.Details

	return true
}
