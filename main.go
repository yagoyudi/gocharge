package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	defaultBatteryThreshold            = 20
	defaultNotificationDurationSeconds = 10
	defaultBatteryPath                 = "/sys/class/power_supply/BAT1"
)

var (
	batteryThreshold            int
	notificationDurationSeconds int
	batteryPath                 string
)

func init() {
	flag.IntVar(
		&batteryThreshold,
		"t",
		defaultBatteryThreshold,
		"Battery threshold (in %)",
	)
	flag.IntVar(
		&notificationDurationSeconds,
		"d",
		defaultNotificationDurationSeconds,
		"Duration of notification (in seconds)",
	)
	flag.StringVar(
		&batteryPath,
		"f",
		defaultBatteryPath,
		"Path to battery file",
	)
	flag.Parse()
}

func main() {
	b := newBattery(batteryPath)

	currentCapacity, err := b.currentCapacity()
	if err != nil {
		exitWithError(err)
	}

	// If battery is above threshold, then delete lock file and print current
	// capacity:
	if currentCapacity >= batteryThreshold {
		_ = defaultLock.delete()
		fmt.Printf("Battery: %d%%\n", currentCapacity)
		return
	}

	// If the battery is below threshold and lock already exists, then just
	// print capacity:
	if defaultLock.exist() {
		fmt.Printf("Battery: %d%%\n", currentCapacity)
		return
	}

	n, err := newNotification(
		withTitle("gocharge"),
		withBody("Your battery is low. Go charge it!"),
		withDuration(notificationDurationSeconds),
	)
	if err != nil {
		exitWithError(err)
	}

	err = n.send()
	if err != nil {
		exitWithError(err)
	}

	err = defaultLock.create()
	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("Battery: %d%%\n", currentCapacity)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "gocharge: %v\n", err)
	os.Exit(1)
}
