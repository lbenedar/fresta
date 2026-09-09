package types

import (
	"time"

	"github.com/lbenedar/fresta/internal/models/json"
)

type FoundryStatus struct {
	IsActive      bool          `json:"on_world"`
	Version       string        `json:"version"`
	World         string        `json:"world,omitempty"`
	System        string        `json:"system,omitempty"`
	SystemVersion string        `json:"-"`
	Users         int           `json:"-"`
	Uptime        time.Duration `json:"-"`
}

func NewFoundryStatus(status *json.Status) *FoundryStatus {
	return &FoundryStatus{
		IsActive:      status.Active,
		Version:       status.Version,
		World:         status.World,
		System:        status.System,
		SystemVersion: status.SystemVersion,
		Users:         status.Users,
		Uptime:        time.Duration(status.Uptime),
	}
}

func (fs *FoundryStatus) Update(status *json.Status) {
	fs.IsActive = status.Active
	fs.Version = status.Version
	fs.World = status.World
	fs.System = status.System
	fs.SystemVersion = status.SystemVersion
	fs.Users = status.Users
	fs.Uptime = time.Duration(status.Uptime)
}
