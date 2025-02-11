package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"context"

	"github.com/coreos/go-systemd/v22/dbus"
)

const (
	Success = 0
	Error   = 1
)

func main() {
	var command string
	var unit string
	flag.StringVar(&command, "command", "reload", "Operation to perform [reload|restart]")
	flag.StringVar(&unit, "unit", "", "Targeted Systemd unit/service")
	flag.Parse()

	ctx := context.Background()
	systemdConnection, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to systemd: %v\n", err)
		panic(err)
	}
	defer systemdConnection.Close()

	completedCh := make(chan string)
	if command == "restart" {
		jobID, err := systemdConnection.RestartUnitContext(ctx, unit, "replace", completedCh)
		if err != nil {
			fmt.Printf("Failed to restart unit: %v\n", err)
			panic(err)
		}
		fmt.Printf("Restart job id: %d\n", jobID)
	} else {
		jobID, err := systemdConnection.ReloadUnitContext(ctx, unit, "replace", completedCh)
		if err != nil {
			fmt.Printf("Failed to reload unit: %v\n", err)
			panic(err)
		}
		fmt.Printf("Reload job id: %d\n", jobID)
	}

	// Wait for the reload to complete
	select {
	case <-completedCh:
		fmt.Printf("Reload job completed for unit: %s\n", unit)
		os.Exit(Success)
	case <-time.After(30 * time.Second):
		fmt.Printf("Timed out waiting for restart job to complete for unit: %s\n", unit)
		os.Exit(Error)
	}
}
