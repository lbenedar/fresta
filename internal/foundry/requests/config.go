package requests

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	ErrorSidWrongFormat = errors.New("Session id wrong format")
	ErrorSidNotFound    = errors.New("Session id didn't find in response header")
)

const (
	AuthPath    = "/auth"
	LicensePath = "/license"
	JoinPath    = "/join"
	PlayersPath = "/players"
	SetupPath   = "/setup"
	UpdatePath  = "/update"
	GamePath    = "/game"

	SocketPath = "/socket.io"
)

var WsTypeData = map[string]string{
	AuthPath:    "getAuthData",
	LicensePath: "getAuthData",
	JoinPath:    "getJoinData",
	PlayersPath: "getPlayersData",
	SetupPath:   "getSetupData",
	UpdatePath:  "getUpdateData",
}

type FoundryHttpRequest struct {
	Host      string
	Password  string
	SessionID *string
	Jar       *cookiejar.Jar
}

func (conf *FoundryHttpRequest) SetSessionTokenFromHeader() (bool, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s", conf.Host))
	if err != nil {
		return false, err
	}

	cookies := conf.Jar.Cookies(u)
	for i := range cookies {
		if cookies[i].Name == "session" {
			conf.SessionID = &cookies[i].Value
			return true, nil
		}
	}

	return false, nil
}

func (conf *FoundryHttpRequest) GetRequestHeader(path string) *http.Header {
	header := http.Header{}
	header.Set("Cookie", fmt.Sprintf("session=%s", *conf.SessionID))
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	header.Set("Connection", "keep-alive")
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("Host", conf.Host)
	header.Set("Origin", fmt.Sprintf("http://%s", conf.Host))
	header.Set("Referer", fmt.Sprintf("http://%s%s", conf.Host, path))

	return &header
}

func (conf *FoundryHttpRequest) GetRequestHeaderJson(path string) *http.Header {
	header := http.Header{}
	header.Set("Accept", "*/*")
	header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	header.Set("Connection", "keep-alive")
	header.Set("Content-Type", "application/json")
	header.Set("Cookie", fmt.Sprintf("session=%s", *conf.SessionID))
	header.Set("Host", conf.Host)
	header.Set("Origin", fmt.Sprintf("http://%s", conf.Host))
	header.Set("Priority", "u=0")
	header.Set("Referer", fmt.Sprintf("http://%s%s", conf.Host, path))

	return &header
}
