package json

import (
	"time"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Release struct {
	Generation          int    `json:"generation"`
	Channel             string `json:"channel"`
	Suffix              string `json:"suffix"`
	Build               int    `json:"build"`
	NodeVersion         int    `json:"node_version"`
	MaxGeneration       int    `json:"maxGeneration"`
	MaxStableGeneration int    `json:"maxStableGeneration"`
	Time                int64  `json:"time"`
	// Flags               Flags  `json:"flags"`
}

func (r *Release) ToDB(dest *db.Release) bool {
	if dest == nil {
		return false
	}

	dest.Generation = r.Generation
	dest.Channel = r.Channel
	dest.Suffix = r.Suffix
	dest.Build = r.Build
	dest.NodeVersion = r.NodeVersion
	dest.MaxGeneration = r.MaxGeneration
	dest.MaxStableGeneration = r.MaxStableGeneration
	dest.Time = time.Unix(r.Time, 0)

	return true
}
