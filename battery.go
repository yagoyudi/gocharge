package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type battery struct {
	path string
}

func newBattery(path string) *battery {
	return &battery{
		path: path,
	}
}

func (b *battery) currentCapacity() (int, error) {
	path := filepath.Join(b.path, "capacity")
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read battery capacity file: %w", err)
	}

	s := strings.TrimSpace(string(data))
	c, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse battery capacity: %w", err)
	}
	return c, nil
}
