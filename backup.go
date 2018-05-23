package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/lgarron/sd-card-backup/sync"
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
	Syncer        sync.Syncer
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

func dateFolderNames(path string) (year string, date string) {
	stat := syscall.Stat_t{}
	syscall.Stat(path, &stat)
	ctime := time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return ctime.Format("2006"), ctime.Format("2006-01-02")
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

	year, date := dateFolderNames(path)

	return filepath.Join(
		fo.Operation.DestinationRoot,
		classificationFolder,
		year,
		date,
		fo.CardName,
		fo.FolderMapping.Destination,
		relPath,
	), nil
}

func (fo folderOperation) syncFile(src string, dest string) error {
	fmt.Printf("\r%s      ", dest)
	if fo.Operation.Options.DryRun {
		fmt.Println()
	} else {
		return fo.Syncer.Queue(src, dest)
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

	return fo.syncFile(path, targetPath)
}

func folderExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
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
		Syncer:        &sync.GoSyncer{},
	}
	return filepath.Walk(folderSourceRoot, fo.visit)
}

// BackupCard backups up the given SD card.
func (op Operation) BackupCard(cardName string) error {
	sdCardPath := filepath.Join(op.SDCardMountPoint, cardName)
	// Check if source folder exists is mounted
	exists, err := folderExists(sdCardPath)
	if err != nil {
		return err
	}
	if !exists {
		fmt.Printf("\r[%s] Skipping SD card (unmounted)\n", cardName)
		return nil
	}

	fmt.Printf("\r[%s] Backing up SD card\n", cardName)

	for _, fc := range classificationBackupOrder {
		for _, fm := range op.FolderMapping {

			folderSourceRoot := filepath.Join(op.SDCardMountPoint, cardName, fm.Source)

			// Check if source folder exists
			exists, err := folderExists(folderSourceRoot)
			if err != nil {
				return err
			}
			if !exists {
				continue
			}

			err = op.backupFolder(cardName, fm, filterClassification(fc))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// BackupAllCards backs up all cards in `op.SDCardNames`.
func (op Operation) BackupAllCards() error {
	// Check if source folder exists
	exists, err := folderExists(op.SDCardMountPoint)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("SD card mount point does not exist: %s", op.DestinationRoot)
	}

	// Check if destination folder exists
	exists, err = folderExists(op.DestinationRoot)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("destination folder does not exist: %s", op.DestinationRoot)
	}

	fmt.Printf("--------\n")
	fmt.Printf("Backing up from:\n  %s\n", op.SDCardMountPoint)
	fmt.Printf("Backing up to:\n  %s\n", op.DestinationRoot)
	fmt.Printf("--------\n")
	for _, s := range op.SDCardNames {
		err := op.BackupCard(s)
		if err != nil {
			return err
		}
	}
	return nil
}
