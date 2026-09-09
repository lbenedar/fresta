package json

import "github.com/lbenedar/fresta/internal/models/db"

type Compatibility struct {
	Minimum  string `json:"minimum,omitempty"`
	Verified string `json:"verified,omitempty"`
	Maximum  string `json:"maximum,omitempty"`
}

func (c *Compatibility) ToDB(dest *db.Compatibility) bool {
	if dest == nil {
		return false
	}

	dest.Minimum = c.Minimum
	dest.Verified = c.Verified
	dest.Maximum = c.Maximum

	return true
}
