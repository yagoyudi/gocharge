package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

type option func(options *options) error

type options struct {
	title           *string
	body            *string
	durationSeconds *int
}

func withTitle(t string) option {
	return func(options *options) error {
		options.title = &t
		return nil
	}
}

func withBody(b string) option {
	return func(options *options) error {
		options.body = &b
		return nil
	}
}

func withDuration(d int) option {
	return func(options *options) error {
		if d < 0 {
			return errors.New("duration can't be negative")
		}
		options.durationSeconds = &d
		return nil
	}
}

type notification struct {
	title           string
	body            string
	durationSeconds int
	durationMs      string
}

func newNotification(opts ...option) (*notification, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var title string
	if options.title == nil {
		title = "gocharge"
	} else {
		title = *options.title
	}

	var body string
	if options.title == nil {
		body = "Your battery is low. Go charge it!"
	} else {
		body = *options.body
	}

	var duration int
	if options.durationSeconds == nil {
		duration = 20
	} else {
		duration = *options.durationSeconds
	}

	n := &notification{
		title:           title,
		body:            body,
		durationSeconds: duration,
		durationMs:      strconv.Itoa(duration * 1000),
	}

	return n, nil
}

// send uses the "notify-send" to alert the user that the battery is low.
func (n *notification) send() error {
	notifySend, err := exec.LookPath("notify-send")
	if err != nil {
		return fmt.Errorf("notify-send not found: %w", err)
	}
	cmd := exec.Command(
		notifySend,
		"-t", n.durationMs,
		n.title,
		n.body,
	)
	return cmd.Run()
}
