package transport

import "errors"

var (
	ErrorNotAuth            = errors.New("Admin is not authenticated")
	ErrorIsNotReady         = errors.New("Connection is not ready for communication")
	ErrorTimeout            = errors.New("Answer has not been received after timeout")
	ErrorStateTypeNotExist  = errors.New("State type does not exists")
	ErrorAuthPassWrong      = errors.New("Authentication password is wrong")
	ErrorChannelIsClosed    = errors.New("Channel has been closed")
	ErrorUserIsNotConnected = errors.New("User is not connected")
	ErrorUserDoesntExist    = errors.New("User does not exist")
)
