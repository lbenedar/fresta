package config

import (
	"errors"
	"fmt"
)

var (
	ErrWrongKindPassed = errors.New("Passed wrond kind of object")
)

type ConfigError struct {
	configName string
	err        error
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("Config %s got error: %s", ce.configName, ce.err)
}
