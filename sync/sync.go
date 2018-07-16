package sync

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mostafah/fsync"
)

// Syncer represents a way to sync a list of files.
type Syncer interface {
	Queue(src string, dest string) error
	// Flushes any queued operations that are not completed, before returning.
	// Flush() error
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

	return fsync.Sync(dest, src)
}

// // Flush is a no-op for GoSyncer.
// func (s GoSyncer) Flush() error {
// 	return nil
// }

// ImmediateRsync shells out queued files to rsync.
type ImmediateRsync struct{}

// Queue syncs the given file in an immediate rsync operation (no batching).
func (r ImmediateRsync) Queue(src string, dest string) error {
	os.MkdirAll(filepath.Dir(dest), 0700)

	cmd := exec.Command("rsync", "-a", "--", src, dest)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// type fileToSync struct {
//  src  string
//  dest string
// }

// // BatchedRsync shells out queued files to rsync, only executing rsync when `Flush` is called.
// type BatchedRsync struct {
//  files []fileToSync
// }

// // Queue queues the given file.
// func (r BatchedRsync) Queue(src string, dest string) error {
//  r.files = append(r.files, fileToSync{src: src, dest: dest})
//  return nil
// }

// // Flush executesRsync
// func (r BatchedRsync) Flush() error {
//  return errors.New("Unimplemented!")
// }
