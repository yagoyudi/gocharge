package main

import (
	"fmt"
	"os"
)

const (
	defaultLockPath = "/tmp/gocharge.lock"
)

var defaultLock = &lock{path: defaultLockPath}

type lock struct {
	path string
}

// create creates the lock file.
func (l lock) create() error {
	file, err := os.Create(l.path)
	if err != nil {
		return err
	}
	// The content of a lock file doesn't matter. Thus, just create the file
	// and close it:
	return file.Close()
}

// delete removes the lock file.
func (l lock) delete() error {
	if !l.exist() {
		return fmt.Errorf("lock does not exist")
	}
	return os.Remove(l.path)
}

// exist returns true is lock file exists.
func (l lock) exist() bool {
	_, err := os.Stat(l.path)
	return err == nil
}
