package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	json_model "github.com/lbenedar/fresta/internal/models/json"
)

func (conf *FoundryHttpRequest) GetSessionId() error {
	path := fmt.Sprintf("http://%s%s", conf.Host, AuthPath)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Jar: conf.Jar}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	_, err = conf.SetSessionTokenFromHeader()
	return err
}

// func (conf *FoundryHttpRequest) GetSessionToken() (bool, error) {
// 	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", conf.Host, AuthPath), nil)
// 	if err != nil {
// 		return false, err
// 	}

// 	header := http.Header{}
// 	header.Set("Cookie", fmt.Sprintf("session=%s", conf.SessionID))
// 	header.Set("Connection", "keep-alive")
// 	header.Set("Host", conf.Host)
// 	header.Set("Origin", fmt.Sprintf("http://%s", conf.Host))
// 	header.Set("Referer", fmt.Sprintf("http://%s%s", conf.Host, AuthPath))

// 	req.Header = header
// 	client := &http.Client{Jar: conf.Jar}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return false, err
// 	}

// 	return conf.SetSessionTokenFromHeader(resp.Header)
// }

func (conf *FoundryHttpRequest) GetStatus() (*json_model.Status, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/status", conf.Host), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Jar: conf.Jar}
	getResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer getResp.Body.Close()

	statusByte := bytes.Buffer{}
	_, err = io.Copy(&statusByte, getResp.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}

	var status json_model.Status
	err = json.Unmarshal(statusByte.Bytes(), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (conf *FoundryHttpRequest) GetGame() (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", conf.Host, GamePath), nil)
	if err != nil {
		return false, err
	}

	header := conf.GetRequestHeader(GamePath)

	req.Header = *header
	client := &http.Client{Jar: conf.Jar}

	_, err = client.Do(req)
	if err != nil {
		return false, err
	}

	return true, nil
}
