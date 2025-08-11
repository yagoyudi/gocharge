# gocharge

Simplest way to get notified through `notify-send` when your battery is low.
Designed for window managers/compositors environments, mainly barless setups.

> [!Warning]
> Only works on Linux. Maybe in BSDs.

## Dependencies

- libnotify

## Usage

1. Use the -f flag to specify the batery file.
2. Add this binary in cron or systemd-timer to run every 5 minutes.

