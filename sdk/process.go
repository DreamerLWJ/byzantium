package sdk

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

// SignProcess send signal to process by pid
func SignProcess(pid int, signal syscall.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Signal(signal)
	if err != nil {
		return err
	}
	return nil
}

// IsProcessAlive check if pid is alive
func IsProcessAlive(pid int) (alive bool, err error) {
	err = SignProcess(pid, syscall.Signal(0))
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
