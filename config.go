package backup

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/adrg/xdg"
)

type folderMapping struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// Operation defines the config file format for the `sd-card-backup`
// commandline tool.
type Operation struct {
	DestinationRoot  string          `json:"destination_root"`
	SDCardMountPoint string          `json:"sd_card_mount_point"`
	SDCardNames      []string        `json:"sd_card_names"`
	FolderMapping    []folderMapping `json:"folder_mapping"`
}

func operationFromBytes(b []byte) (*Operation, error) {
	config := Operation{}
	// TODO: Enforce required fields to catch config typos?
	err := json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// OperationFromConfig reads from the global XDG config path.
func OperationFromConfig() (*Operation, error) {
	configPath := filepath.Join(xdg.ConfigHome, "sd-card-backup", "config.json")
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return operationFromBytes(file)
}
