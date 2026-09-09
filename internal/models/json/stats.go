package json

import "github.com/lbenedar/fresta/internal/models/db"

type Stats struct {
	CoreVersion    string `json:"coreVersion"`
	SystemID       string `json:"systemId"`
	SystemVersion  string `json:"systemVersion"`
	LastModifiedBy string `json:"lastModifiedBy"`
	ModifiedTime   int64  `json:"modifiedTime"`
	// CompendiumSource string `json:"compendiumSource,omitempty"`
	// DuplicateSource  string `json:"duplicateSource,omitempty"`
	// ExportSource     string `json:"exportSource,omitempty"`
}

func (s *Stats) ToDB(dest *db.Stats) bool {
	if dest == nil {
		return false
	}

	dest.CoreVersion = s.CoreVersion
	dest.SystemID = s.SystemID
	dest.SystemVersion = s.SystemVersion
	dest.LastModifiedBy = s.LastModifiedBy
	dest.ModifiedTime = s.ModifiedTime

	return true
}
