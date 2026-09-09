package transport

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

func (tr *FoundryTransport) LaunchWorld(worldName string) error {
	resp, err := tr.Http.PostLaunchWorld(worldName)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("World is not found")
	}

	tr.ExchangeChan.ProgressMsg[worldName] = make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	select {
	case <-tr.ExchangeChan.ProgressMsg[worldName]:
		tr.Logger.Debug("World has been launched(from GetInitialData)", "world", worldName)
	case <-ctx.Done():
		close(tr.ExchangeChan.ProgressMsg[worldName])
		err := ctx.Err()
		if err != nil {
			return err
		}
		return ErrorTimeout
	}
	return nil
}

func (tr *FoundryTransport) ConnectToFoundry() error {
	if tr.Http.Password != "" {
		err := tr.Authenticate()
		if err != nil {
			return err
		}
	}

	tr.Logger.Info("Successefully connected to Foundry", "host", tr.Http.Host)
	return nil
}

func (tr *FoundryTransport) Authenticate() error {
	if tr.Http.SessionID == nil {
		err := tr.Http.GetSessionId()
		if err != nil {
			return err
		}
	}

	resp, err := tr.Http.PostAuthenticationData()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentLen, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	return tr.CheckAuthResponse(resp.Body, contentLen)
}

func (tr *FoundryTransport) CheckAuthResponse(body io.Reader, respLength int) error {
	authBody := make([]byte, respLength)

	_, err := body.Read(authBody)
	if err != nil && err != io.EOF {
		return err
	}

	if !strings.Contains(string(authBody), "Admin authentication successful") {
		return ErrorAuthPassWrong
	}
	tr.IsAuth = true
	return nil
}

func (tr *FoundryTransport) LogInToWorld(userId string, userPass string) error {
	if tr.Http.SessionID == nil {
		err := tr.Http.GetSessionId()
		if err != nil {
			return err
		}
	}

	resp, err := tr.Http.PostJoinWorld(userId, userPass)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentLen, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	return tr.CheckLoginResponse(resp.Body, contentLen)
}

func (tr *FoundryTransport) CheckLoginResponse(body io.Reader, respLength int) error {
	joinBody := make([]byte, respLength)

	_, err := body.Read(joinBody)
	if err != nil && err != io.EOF {
		return err
	}

	tr.Logger.Debug("Received response", "response", string(joinBody))

	if !strings.Contains(string(joinBody), "\"status\":\"success\"") {
		return ErrorAuthPassWrong
	}
	// tr.IsAuth = true
	return nil
}

func (tr *FoundryTransport) HasSessionId() bool {
	return tr.Http.SessionID != nil
}

func (tr *FoundryTransport) InitSessionId() error {
	return tr.Http.GetSessionId()
}
