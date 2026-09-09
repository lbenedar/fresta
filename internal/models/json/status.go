package json

type Status struct {
	Active        bool   `json:"active"`
	Version       string `json:"version"`
	World         string `json:"world,omitempty"`
	System        string `json:"system,omitempty"`
	SystemVersion string `json:"systemVersion,omitempty"`
	Users         int    `json:"users,omitempty"`
	Uptime        int64  `json:"uptime,omitempty"`
}

func (s *Status) Copy(copyStatus *Status) {
	s.Active = copyStatus.Active
	s.System = copyStatus.System
	s.SystemVersion = copyStatus.SystemVersion
	s.Uptime = copyStatus.Uptime
	s.Users = copyStatus.Users
	s.Version = copyStatus.Version
	s.World = copyStatus.World
}
