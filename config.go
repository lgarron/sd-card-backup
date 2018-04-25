package backup

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/adrg/xdg"
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
