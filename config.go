package backup

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (fm folderMapping) validate() error {
	if fm.Source == "" {
		return fmt.Errorf("missing `source` in folder mapping: %+v", fm)
	}
	if fm.Destination == "" {
		return fmt.Errorf("missing `destination` in folder mapping: %+v", fm)
	}
	return nil
}

func (o Operation) validate() error {
	if o.DestinationRoot == "" {
		return errors.New("missing `destination_root`")
	}
	if o.SDCardMountPoint == "" {
		return errors.New("missing `sd_card_mount_point`")
	}
	if o.SDCardNames == nil {
		return errors.New("missing `sd_card_names`")
	}
	if len(o.SDCardNames) == 0 {
		return errors.New("empty `sd_card_names`")
	}
	for _, c := range o.SDCardNames {
		if c == "" {
			return errors.New("contains empty SD card name")
		}
	}
	if o.FolderMapping == nil {
		return errors.New("missing `folder_mapping`")
	}
	if len(o.FolderMapping) == 0 {
		return errors.New("empty `folder_mapping`")
	}
	for _, fm := range o.FolderMapping {
		err := fm.validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func operationFromBytes(b []byte) (*Operation, error) {
	config := Operation{}
	// TODO: Enforce required fields to catch config typos?
	err := json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func xdgConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "sd-card-backup", "config.json")
}

// OperationFromConfig reads from the global config path.
func OperationFromConfig() (*Operation, error) {
	file, err := ioutil.ReadFile(xdgConfigPath())
	if err != nil {
		return nil, err
	}

	return operationFromBytes(file)
}
