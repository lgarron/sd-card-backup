package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mostafah/fsync"
)

type folderOperation struct {
	Operation     Operation
	SourceRoot    string
	CardName      string
	FolderMapping folderMapping
}

func folderForClassification(classification int) (string, error) {
	switch classification {
	case imageFile:
		return "Images", nil
	case videoFile:
		return "Videos", nil
	case unclassifiedFile:
		return "Unsorted", nil
	default:
		return "", fmt.Errorf("Unknown classification: %d", classification)
	}
}

func monthFolderName(path string) string {
	stat := syscall.Stat_t{}
	syscall.Stat(path, &stat)
	ctime := time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return ctime.Format("2006-01")
}

func (fo folderOperation) targetPath(path string, f os.FileInfo) (string, error) {
	classificationFolder, err := folderForClassification(classifyPath(path))
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(fo.SourceRoot, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(
		fo.Operation.DestinationRoot,
		classificationFolder,
		monthFolderName(path),
		fo.CardName,
		fo.FolderMapping.Destination,
		relPath,
	), nil
}

func (fo folderOperation) syncFile(dest string, src string) error {
	fmt.Printf("\r%s      ", dest)
	if fo.Operation.Options.DryRun {
		fmt.Println()
	} else {
		os.MkdirAll(filepath.Dir(dest), 0700)
		err := fsync.Sync(dest, src)
		if err != nil {
			return err
		}

		err = fsync.Sync(dest, src)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fo folderOperation) visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if f.IsDir() {
		return nil
	}

	targetPath, err := fo.targetPath(path, f)
	if err != nil {
		return err
	}

	return fo.syncFile(targetPath, path)
}

// Backups up:
//
//   [op.SDCardMountPoint]/[cardName]/[fm.Source]/[filePath]
//
// to:
//
//   [op.DestinationRoot]/[classification]/[month]/[cardName]/[fm.Destination]/[filePath]
//
func (op Operation) backupFolder(cardName string, fm folderMapping) error {
	folderSourceRoot := filepath.Join(op.SDCardMountPoint, cardName, fm.Source)
	fo := &folderOperation{
		Operation:     op,
		SourceRoot:    folderSourceRoot,
		CardName:      cardName,
		FolderMapping: fm,
	}
	return filepath.Walk(folderSourceRoot, fo.visit)
}

// BackupCard backups up the given SD card.
func (op Operation) BackupCard(cardName string) error {
	for _, fm := range op.FolderMapping {
		err := op.backupFolder(cardName, fm)
		if err != nil {
			return err
		}
	}
	return nil
}
