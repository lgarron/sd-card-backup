package backup

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// From https://go-review.googlesource.com/c/go/+/1591
func copyFile(dst string, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dst)
}

// smartCopy looks at the destination:
//
// Destination empty: copy, preserving timestamp
// Destination exists, same timestamp and size: do nothing
// Destination exists, same timestamp but not same size: copy, preserving timestamp
// Destination exists, newer timestamp: HOLD THE PRESS
// Destination exists, older timestamp: HOLD THE PRESS
func smartCopy(dst string, src string, perm os.FileMode) {

}
