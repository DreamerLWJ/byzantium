package process_checker

import "github.com/pkg/errors"

var (
	ErrTerminated = errors.Errorf("checker is terminated")
)
