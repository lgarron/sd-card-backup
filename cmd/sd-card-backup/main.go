package main

import (
	"fmt"
	"os"

	backup "github.com/lgarron/sd-card-backup"
)

func main() {
	op, err := backup.OperationFromConfig()
	if err != nil {
		fmt.Printf("Could not read config file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Backing up to: %s", op.DestinationRoot)
	for _, s := range op.SDCardNames {
		op.BackupCard(s)
	}
}
