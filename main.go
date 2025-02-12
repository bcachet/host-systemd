package main

import (
	"flag"
	"fmt"
	"log"
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
	flag.Parse()
	if flag.NArg() < 2 {
		Usage()
		os.Exit(Error)
	}
	var command string = flag.Arg(0)
	var unit string = flag.Arg(1)

	ctx := context.Background()
	systemdConnection, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to systemd: %v\n", err)
		panic(err)
	}
	defer systemdConnection.Close()

	completedCh := make(chan string)
	timeout := time.AfterFunc(30*time.Second, func() {
		close(completedCh) // Close the channel to trigger the select
	})
	_, err = executeSystemdCommand(systemdConnection, ctx, command, unit, completedCh)
	if err != nil {
		log.Fatalf("Failed to %s unit: %v", command, err)
	}

	// Wait for the operation to complete or timeout
	select {
	case <-completedCh:
		log.Printf("%s job completed for unit: %s", command, unit)
		os.Exit(Success)
	case <-timeout.C:
		log.Printf("Timed out waiting for %s job to complete for unit: %s", command, unit)
		os.Exit(Error)
	}
}

func executeSystemdCommand(c *dbus.Conn, ctx context.Context, command string, unit string, completedCh chan string) (int, error) {
	switch command {
	case "restart":
		return c.RestartUnitContext(ctx, unit, "replace", completedCh)
	case "reload":
		return c.ReloadUnitContext(ctx, unit, "replace", completedCh)
	default:
		return 0, fmt.Errorf("invalid command: %v. Valid commands: reload|restart", command)
	}
}

func Usage() {
	fmt.Printf("usage: %s COMMAND UNIT\n\n", os.Args[0])
	flag.PrintDefaults()
}
