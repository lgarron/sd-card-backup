package backup

import (
	"errors"
	"fmt"
)

type folderMapping struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type commandLineOptions struct {
	DryRun bool
}

// Operation defines the config file format for the `sd-card-backup`
// commandline tool.
type Operation struct {
	DestinationRoot  string          `json:"destination_root"`
	SDCardMountPoint string          `json:"sd_card_mount_point"`
	SDCardNames      []string        `json:"sd_card_names"`
	FolderMapping    []folderMapping `json:"folder_mapping"`
	// TODO: the following should be a tuple, but Go is inadequate for that.
	CommandToRunBefore []string `json:"command_to_run_before"` // Contains a command and arguments as entries.
	Options            commandLineOptions
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
	// TODO: condense calculations similar to table-driven tests.
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
			return errors.New("contains empty card name")
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
