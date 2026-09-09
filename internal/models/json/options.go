package json

import "github.com/lbenedar/fresta/internal/models/db"

type GameOptions struct {
	Language      string `json:"language"`
	Port          int    `json:"port"`
	RoutePrefix   any    `json:"routePrefix"`
	UpdateChannel string `json:"updateChannel"`
}

func (g *GameOptions) ToDB(dest *db.GameOptions) bool {
	if dest == nil {
		return false
	}

	dest.Language = g.Language
	dest.Port = g.Port
	dest.UpdateChannel = g.UpdateChannel

	return true
}

type SetupOptions struct {
	AwsConfig      any    `json:"awsConfig"`
	CompressSocket bool   `json:"compressSocket"`
	CompressStatic bool   `json:"compressStatic"`
	CSSTheme       string `json:"cssTheme"`
	DataPath       string `json:"dataPath"`
	Fullscreen     bool   `json:"fullscreen"`
	Hostname       string `json:"hostname"`
	HotReload      bool   `json:"hotReload"`
	Language       string `json:"language"`
	LocalHostname  string `json:"localHostname"`
	Port           int    `json:"port"`
	ProxySSL       bool   `json:"proxySSL"`
	Telemetry      bool   `json:"telemetry"`
	UpdateChannel  string `json:"updateChannel"`
	Upnp           bool   `json:"upnp"`
	DeleteNEDB     bool   `json:"deleteNEDB"`
	NoBackups      bool   `json:"noBackups"`
	// PasswordSalt      any    `json:"passwordSalt"`
	// Protocol          any    `json:"protocol"`
	// ProxyPort         any    `json:"proxyPort"`
	// RoutePrefix       any    `json:"routePrefix"`
	// SslCert           any    `json:"sslCert"`
	// SslKey            any    `json:"sslKey"`
	// UpnpLeaseDuration any    `json:"upnpLeaseDuration"`
	// World             any    `json:"world"`
	// AdminPassword string `json:"adminPassword"`
}

func (s *SetupOptions) ToDB(dest **db.SetupOptions) bool {
	if dest == nil {
		return false
	}

	options := &db.SetupOptions{
		CompressSocket: s.CompressSocket,
		CompressStatic: s.CompressStatic,
		CSSTheme:       s.CSSTheme,
		DataPath:       s.DataPath,
		Fullscreen:     s.Fullscreen,
		Hostname:       s.Hostname,
		HotReload:      s.HotReload,
		Language:       s.Language,
		LocalHostname:  s.LocalHostname,
		Port:           s.Port,
		ProxySSL:       s.ProxySSL,
		Telemetry:      s.Telemetry,
		UpdateChannel:  s.UpdateChannel,
		Upnp:           s.Upnp,
		DeleteNEDB:     s.DeleteNEDB,
		NoBackups:      s.NoBackups,
	}

	*dest = options

	return true
}
