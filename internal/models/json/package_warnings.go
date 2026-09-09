package json

import "github.com/lbenedar/fresta/internal/models/db"

type PackageWarningsData struct {
	Id            string   `json:"id"`
	Type          string   `json:"type"`
	Warning       []string `json:"warning"`
	Error         []string `json:"error"`
	Reinstallable bool     `json:"reinstallable"`
	Manifest      string   `json:"manifest,omitempty"`
}

func (p *PackageWarningsData) ToDB(dest **db.PackageWarningsData) bool {
	if dest == nil {
		return false
	}

	packageData := &db.PackageWarningsData{
		ID:            p.Id,
		Type:          p.Type,
		Reinstallable: p.Reinstallable,
		Manifest:      p.Manifest,
	}

	packageData.Warning = make([]string, len(p.Warning))
	copy(packageData.Warning, p.Warning)

	packageData.Error = make([]string, len(p.Error))
	copy(packageData.Error, p.Error)

	*dest = packageData

	return true
}
