package types

import "fmt"

type FoundryError struct {
	Err       error
	Direction DirectionCode
	Type      TransportCode
	IsFatal   bool
}

func (e *FoundryError) Error() string {
	return fmt.Sprintf("Received from %d: %s", e.Type, e.Err.Error())
}

func (e *FoundryError) Unwrap() error { return e.Err }

// func NewWriterError(err error, isFatal bool) *FoundryError {
// 	return &FoundryError{Direction: WriterCode, Err: err, IsFatal: isFatal}
// }

// func NewReaderError(err error, isFatal bool) *FoundryError {
// 	return &FoundryError{Type: ReaderCode, Err: err, IsFatal: isFatal}
// }
