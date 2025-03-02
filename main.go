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
	lowBattery := flag.Int("l", 20, "When the battery is considered low (in %)")
	notificationDuration := flag.Int("d", 10, "Duration of the notification (in seconds)")
	capacityPath := flag.String("f", "/sys/class/power_supply/BAT1/capacity", "Path to the battery capacity file")
	flag.Parse()

	capacityBytes, err := os.ReadFile(*capacityPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gocharge: %v\n", err)
		os.Exit(1)
	}

	capacityStr := strings.TrimSpace(string(capacityBytes))
	capacity, err := strconv.Atoi(capacityStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gocharge: %v\n", err)
		os.Exit(1)
	}

	if capacity < *lowBattery {
		if _, err := os.Stat(lockFile); os.IsNotExist(err) {
			if err := notify(*notificationDuration); err != nil {
				fmt.Fprintf(os.Stderr, "gocharge: %v\n", err)
				os.Exit(1)
			}
			_, err := os.Create(lockFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "gocharge: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		if _, err := os.Stat(lockFile); err == nil {
			os.Remove(lockFile)
		}
	}

	fmt.Printf("Battery: %d%%\n", capacity)
}

func notify(duration int) error {
	durationString := strconv.Itoa(duration * 1000)
	cmd := exec.Command("notify-send", "-t", durationString, "gocharge", lowBatteryWarning)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
