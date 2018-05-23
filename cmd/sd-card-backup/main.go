package main

import (
	"flag"
	"fmt"
	"os"

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

	err = op.BackupAllCards()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error backing up: %s\n", err)
		os.Exit(1)
	}
}
