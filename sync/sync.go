package sync

import (
	"os"
	"path/filepath"

	"github.com/mostafah/fsync"
)

// Syncer represents a way to sync a list of files.
type Syncer interface {
	Queue(src string, dest string) error
}

// GoSyncer is a Syncer implemented in Go.
type GoSyncer struct{}

// Queue syncs immediately using `fsync.Sync`.
func (s GoSyncer) Queue(src string, dest string) error {
	os.MkdirAll(filepath.Dir(dest), 0700)
	err := fsync.Sync(dest, src)
	if err != nil {
		return err
	}

	err = fsync.Sync(dest, src)
	if err != nil {
		return err
	}

	return nil
}
