package json

import "github.com/lbenedar/fresta/internal/models/db"

type Addresses struct {
	Local              string `json:"local"`
	Remote             string `json:"remote"`
	RemoteIsAccessible bool   `json:"remoteIsAccessible"`
}

func (a *Addresses) ToDB(dest *db.Addresses) bool {
	if dest == nil {
		return false
	}

	dest.Local = a.Local
	dest.Remote = a.Remote
	dest.RemoteIsAccessible = a.RemoteIsAccessible

	return true
}
