package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	return filepath.Join(homeDir, ".config", "sd-card-backup", "config.json")
}

// OperationFromConfig reads from the global config path.
func OperationFromConfig() (*Operation, error) {
	file, err := ioutil.ReadFile(xdgConfigPath())
	if err != nil {
		return nil, err
	}

	return operationFromBytes(file)
}
