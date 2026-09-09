package json

import "github.com/lbenedar/fresta/internal/models/db"

type Language struct {
	Lang string `json:"lang"`
	Name string `json:"name"`
	Path string `json:"path"`
	// SystemLanguagesFlags SystemLanguagesFlags `json:"flags"`
}

func (l *Language) ToDB(dest **db.Language) bool {
	if dest == nil {
		return false
	}

	lang := &db.Language{
		Lang: l.Lang,
		Name: l.Name,
		Path: l.Path,
	}

	*dest = lang

	return true
}

type SetupLanguage struct {
	ID      string                 `json:"id"`
	Label   string                 `json:"label"`
	Modules []*SetupLanguageModule `json:"modules"`
}

func (s *SetupLanguage) ToDB(dest **db.SetupLanguage) bool {
	if dest == nil {
		return false
	}

	lang := &db.SetupLanguage{
		ID:    s.ID,
		Label: s.Label,
	}

	CopySliceToDB(&lang.Modules, s.Modules)

	*dest = lang

	return true
}

type SetupLanguageModule struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Path  string `json:"path"`
}

func (s *SetupLanguageModule) ToDB(dest *db.SetupLanguageModule) bool {
	if dest == nil {
		return false
	}

	dest.ID = s.ID
	dest.Label = s.Label
	dest.Path = s.Path

	return true
}
