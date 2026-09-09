package actions

import "errors"

var (
	ErrActionNotFound      = errors.New("Action doesn't found")
	ErrObjTypeNotMatch     = errors.New("Provided value type didn't match obj field type")
	ErrObjNotPointer       = errors.New("Passed obj is not pointer")
	ErrUserChannelIsClosed = errors.New("User channel is closed")
)
