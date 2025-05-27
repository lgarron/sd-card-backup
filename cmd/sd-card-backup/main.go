package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	backup "github.com/lgarron/sd-card-backup"
)

var dryRun = flag.Bool("dry-run", false, "Print what would happen, but don't modify the filesystem.")

func main() {
	// Try to parse flags before doing anything.
	flag.Parse()

	op, err := backup.OperationFromConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config file: %s\n", err)
		os.Exit(1)
	}

	op.Options.DryRun = *dryRun

	if len(op.CommandToRunBefore) > 0 {
		if op.Options.DryRun {
			fmt.Printf("Skipping the following `command_to_run_before` due to dry run: %#v\n", op.CommandToRunBefore)
		} else {
			// TODO: use https://github.com/lgarron/printable-shell-command once we port this.
			fmt.Printf("Running command: %#v\n", op.CommandToRunBefore)
			cmd := exec.Command(op.CommandToRunBefore[0], op.CommandToRunBefore[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}
	}

	err = op.BackupAllCards()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error backing up: %s\n", err)
		os.Exit(1)
	}
}
