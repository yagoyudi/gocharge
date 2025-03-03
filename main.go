package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	lowBatteryWarning = "Your battery is low. Go charge it!"
	lockFile          = "/tmp/gocharge.lock"
)

func main() {
	threshold, duration, capacityPath := parseFlags()

	capacity, err := readBatteryCapacity(capacityPath)
	if err != nil {
		exitWithError(err)
	}

	if capacity < threshold {
		handleLowBattery(duration)
	} else {
		removeLockFileIfExists()
	}

	fmt.Printf("Battery: %d%%\n", capacity)
}

// parseFlags reads command-line flags and returns the battery threshold,
// notification duration, and battery capacity file path.
func parseFlags() (int, int, string) {
	lowBattery := flag.Int("l", 20, "When the battery is considered low (in %)")
	notificationDuration := flag.Int("d", 10, "Duration of the notification (in seconds)")
	capacityPath := flag.String("f", "/sys/class/power_supply/BAT1/capacity", "Path to the battery capacity file")
	flag.Parse()
	return *lowBattery, *notificationDuration, *capacityPath
}

// readBatteryCapacity reads and parses the battery capacity from the given
// file path.
func readBatteryCapacity(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read battery capacity file: %w", err)
	}

	capacityStr := strings.TrimSpace(string(data))
	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse battery capacity: %w", err)
	}
	return capacity, nil
}

// handleLowBattery sends a notification if the battery is low and no previous
// notification was sent (lock file not present).
func handleLowBattery(duration int) {
	if _, err := os.Stat(lockFile); os.IsNotExist(err) {
		if err := sendNotification(duration); err != nil {
			exitWithError(err)
		}
		if err := createLockFile(); err != nil {
			exitWithError(err)
		}
	}
}

// sendNotification uses the "notify-send" command to alert the user that the
// battery is low.
func sendNotification(duration int) error {
	durationMs := strconv.Itoa(duration * 1000)
	notifySend, err := exec.LookPath("notify-send")
	if err != nil {
		return fmt.Errorf("notify-send not found: %w", err)
	}
	cmd := exec.Command(
		notifySend,
		"-t", durationMs,
		"gocharge",
		lowBatteryWarning,
	)
	return cmd.Run()
}

// createLockFile creates the lock file so that the notification is not sent
// repeatedly.
func createLockFile() error {
	file, err := os.Create(lockFile)
	if err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}
	return file.Close()
}

// removeLockFileIfExists removes the lock file if it exists.
func removeLockFileIfExists() {
	if _, err := os.Stat(lockFile); err == nil {
		_ = os.Remove(lockFile)
	}
}

// exitWithError prints the error message to stderr and exits the program.
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "gocharge: %v\n", err)
	os.Exit(1)
}
