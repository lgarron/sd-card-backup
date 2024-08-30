package sync

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/mostafah/fsync"
)

const BYTES_IN_MEGABYTE = 1000 * 1000

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

	cmd := exec.Command("cp", src, dest)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type MacOSNativeCpUsingFilesizeAndBirthTime struct {
}

func NewMacOSNativeCpUsingFilesizeAndBirthTime() MacOSNativeCpUsingFilesizeAndBirthTime {
	if runtime.GOOS != "darwin" {
		fmt.Printf("Running a macOS-specific sync implementation outside macOS. Exiting.")
		os.Exit(1)
	}
	return MacOSNativeCpUsingFilesizeAndBirthTime{}
}

func (s MacOSNativeCpUsingFilesizeAndBirthTime) Queue(src string, dest string) error {
	same, srcStat, err := s.fileIsSameHeuristic(src, dest)
	if err != nil {
		return err
	}

	if same {
		fmt.Printf(" (skipping: already backed up)")
		return nil
	}

	fmt.Printf("\n↪ %s (%d MB)", dest, srcStat.Size/BYTES_IN_MEGABYTE)

	os.MkdirAll(filepath.Dir(dest), 0700)

	{
		// TODO: output progress using https://unix.stackexchange.com/questions/66795/how-to-check-progress-of-running-cp#:~:text=On%20recent%20versions%20of%20Mac,written%20to%20the%20standard%20output.%22

		// `-p` copies modification and access time, but not creation (birth) time.
		cmd := exec.Command("cp", src, dest)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			// TODO: more info about the file we failed on?
			cmd := exec.Command("open", "-R", dest)
			cmd.Stderr = os.Stderr
			err = cmd.Run()

			return err
		}
	}

	{
		// TODO: output progress using https://unix.stackexchange.com/questions/66795/how-to-check-progress-of-running-cp#:~:text=On%20recent%20versions%20of%20Mac,written%20to%20the%20standard%20output.%22

		// `-p` copies modification and access time, but not creation (birth) time.
		cmd := exec.Command("touch", "-r", src, dest)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			// TODO: more info about the file we failed on?
			return err
		}
	}

	{
		cmd := exec.Command("GetFileInfo", "-d", src)
		cmd.Stderr = os.Stderr
		birthTimeStringBytesFromMacOS, err := cmd.Output()
		if err != nil {
			// TODO: more info about the file we failed on?
			return err
		}
		birthTimeStringFromMacOS := strings.TrimSuffix(string(birthTimeStringBytesFromMacOS), "\n")

		formattedTimeFromStat := time.Unix(srcStat.Birthtimespec.Sec, srcStat.Birthtimespec.Nsec).Format("01/02/2006 15:04:05")

		if birthTimeStringFromMacOS != formattedTimeFromStat {
			// TODO: remove the `birthTimeStringFromMacOS` calculation once these have been stress tested across time zones.
			return fmt.Errorf("incompatible times: (%v, %v)", birthTimeStringFromMacOS, formattedTimeFromStat)
		}

		{
			cmd2 := exec.Command("SetFile", "-d", string(birthTimeStringFromMacOS), dest)
			cmd2.Stderr = os.Stderr
			err = cmd2.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Returns src stat if there was no error.
func (s MacOSNativeCpUsingFilesizeAndBirthTime) fileIsSameHeuristic(src string, dest string) (bool, *syscall.Stat_t, error) {
	if filepath.Base(src) != filepath.Base((dest)) {
		return false, nil, errors.New("heuristic encountered two files with different base names")
	}

	srcStat := syscall.Stat_t{}
	err := syscall.Stat(src, &srcStat)
	if err != nil {
		return false, nil, err
	}

	destStat := syscall.Stat_t{}
	err = syscall.Stat(dest, &destStat)
	if err != nil {
		if err.Error() == "no such file or directory" {
			return false, &srcStat, nil
		}
		return false, nil, err
	}

	if srcStat.Size != destStat.Size {
		fmt.Printf("\n↪️ file size differs: %d src bytes vs. %d dest bytes", srcStat.Size, destStat.Size)
		return false, &srcStat, nil
	}

	if srcStat.Birthtimespec.Sec != destStat.Birthtimespec.Sec {
		fmt.Printf("\n↪️ birth time differs: %d src vs. %d dest", srcStat.Birthtimespec.Sec, destStat.Birthtimespec.Sec)
		return false, &srcStat, nil
	}

	return true, &srcStat, nil
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
