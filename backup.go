package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mostafah/fsync"
)

type fileFilter = func(fileClassification) bool

var classificationBackupOrder = []fileClassification{
	imageFile,
	videoFile,
	unclassifiedFile,
}

func filterClassification(want fileClassification) fileFilter {
	return func(have fileClassification) bool {
		return want == have
	}
}

type folderOperation struct {
	Operation     Operation
	SourceRoot    string
	CardName      string
	FolderMapping folderMapping
	FileFilter    fileFilter
}

func folderForClassification(classification fileClassification) (string, error) {
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
	if !fo.FileFilter(classifyPath(path)) {
		return nil
	}

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
func (op Operation) backupFolder(cardName string, fm folderMapping, ff fileFilter) error {
	folderSourceRoot := filepath.Join(op.SDCardMountPoint, cardName, fm.Source)
	fo := &folderOperation{
		Operation:     op,
		SourceRoot:    folderSourceRoot,
		CardName:      cardName,
		FolderMapping: fm,
		FileFilter:    ff,
	}
	return filepath.Walk(folderSourceRoot, fo.visit)
}

// BackupCard backups up the given SD card.
func (op Operation) BackupCard(cardName string) error {
	for _, fc := range classificationBackupOrder {
		for _, fm := range op.FolderMapping {
			err := op.backupFolder(cardName, fm, filterClassification(fc))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// BackupAllCards backs up all cards in `op.SDCardNames`.
func (op Operation) BackupAllCards() error {
	fmt.Printf("Backing up to: %s\n", op.DestinationRoot)
	for _, s := range op.SDCardNames {
		err := op.BackupCard(s)
		if err != nil {
			return err
		}
	}
	return nil
}
