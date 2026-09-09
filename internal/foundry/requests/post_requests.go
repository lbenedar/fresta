package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (conf *FoundryHttpRequest) PostAuthenticationData() (*http.Response, error) {
	authData := url.Values{}
	authData.Add("adminPassword", conf.Password)
	authData.Add("action", "adminAuth")

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", conf.Host, AuthPath), bytes.NewBuffer([]byte(authData.Encode())))
	if err != nil {
		return nil, err
	}

	header := conf.GetRequestHeader(AuthPath)
	header.Set("Content-Length", fmt.Sprintf("%d", len(authData.Encode())))

	req.Header = header.Clone()

	client := &http.Client{Jar: conf.Jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (conf *FoundryHttpRequest) PostJoinWorld(userId, password string) (*http.Response, error) {
	data := struct {
		UserId   string `json:"userid"`
		Password string `json:"password"`
		Action   string `json:"action"`
		Session  string `json:"session"`
	}{
		UserId:   userId,
		Password: password,
		Action:   "join",
		Session:  *conf.SessionID,
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", conf.Host, JoinPath), bytes.NewBuffer(byteData))
	if err != nil {
		return nil, err
	}

	header := conf.GetRequestHeaderJson(JoinPath)
	header.Set("Content-Length", fmt.Sprintf("%d", len(byteData)))

	req.Header = header.Clone()

	client := &http.Client{Jar: conf.Jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (conf *FoundryHttpRequest) PostLaunchWorld(world string) (*http.Response, error) {
	launchData := url.Values{}
	launchData.Add("world", world)
	launchData.Add("action", "launchWorld")

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", conf.Host, SetupPath), bytes.NewBuffer([]byte(launchData.Encode())))
	if err != nil {
		return nil, err
	}

	header := conf.GetRequestHeader(SetupPath)
	header.Set("Content-Length", fmt.Sprintf("%d", len(launchData.Encode())))

	req.Header = header.Clone()

	client := &http.Client{Jar: conf.Jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (conf *FoundryHttpRequest) PostReturnToSetup() (*http.Response, error) {
	launchData := url.Values{}
	launchData.Add("action", "shutdown")

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", conf.Host, JoinPath), bytes.NewBuffer([]byte(launchData.Encode())))
	if err != nil {
		return nil, err
	}

	header := conf.GetRequestHeader(JoinPath)
	header.Set("Content-Length", fmt.Sprintf("%d", len(launchData.Encode())))

	req.Header = header.Clone()

	client := &http.Client{Jar: conf.Jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
