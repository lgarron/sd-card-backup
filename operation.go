package backup

import (
	"fmt"
	"os"
	"path/filepath"
)

func (op Operation) visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}

	fmt.Printf("[%d] Visited: %s\n", classifyFile(f), path)
	return nil
}

func (op Operation) backupFolder(cardName string, fm folderMapping) {
	path := filepath.Join(op.SDCardMountPoint, cardName, fm.Source)
	filepath.Walk(path, op.visit)
}

// BackupCard backups up the given SD card.
func (op Operation) BackupCard(cardName string) {
	for _, fm := range op.FolderMapping {
		op.backupFolder(cardName, fm)
	}
}
