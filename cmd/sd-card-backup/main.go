package main

import (
	"fmt"
	"os"

	backup "github.com/lgarron/sd-card-backup"
)

func main() {
	op, err := backup.OperationFromConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Backing up to: %s\n", op.DestinationRoot)
	for _, s := range op.SDCardNames {
		err := op.BackupCard(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error backing up: %s\n", err)
			os.Exit(1)
		}
	}
}
