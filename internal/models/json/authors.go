package json

import "github.com/lbenedar/fresta/internal/models/db"

type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Email   string `json:"email,omitempty"`
	Discord string `json:"discord,omitempty"`
	// SystemAuthorsFlags SystemAuthorsFlags `json:"flags"`
}

func (a *Author) ToDB(dest **db.Author) bool {
	if dest == nil {
		return false
	}

	author := &db.Author{
		Name:    a.Name,
		URL:     a.URL,
		Email:   a.Email,
		Discord: a.Discord,
	}

	*dest = author

	return true
}
