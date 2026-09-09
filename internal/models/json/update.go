package json

import "github.com/lbenedar/fresta/internal/models/db"

type CoreUpdate struct {
	HasUpdate          bool   `json:"hasUpdate"`
	CanUpdate          bool   `json:"canUpdate"`
	CouldReachWebsite  bool   `json:"couldReachWebsite"`
	SlowResponse       bool   `json:"slowResponse"`
	Version            string `json:"version"`
	Channel            string `json:"channel"`
	WillDisableModules bool   `json:"willDisableModules"`
}

func (c *CoreUpdate) ToDB(dest *db.CoreUpdate) bool {
	if dest == nil {
		return false
	}

	dest.HasUpdate = c.HasUpdate
	dest.CanUpdate = c.CanUpdate
	dest.CouldReachWebsite = c.CouldReachWebsite
	dest.SlowResponse = c.SlowResponse
	dest.Version = c.Version
	dest.Channel = c.Channel
	dest.WillDisableModules = c.WillDisableModules

	return true
}

type SystemUpdate struct {
	HasUpdate bool   `json:"hasUpdate"`
	Version   string `json:"version"`
}

func (s *SystemUpdate) ToDB(dest *db.SystemUpdate) bool {
	if dest == nil {
		return false
	}

	dest.HasUpdate = s.HasUpdate
	dest.Version = s.Version

	return true
}
