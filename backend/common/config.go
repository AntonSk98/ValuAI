package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadProperties loads YAML configuration from the specified file path into the provided struct.
func LoadProperties(path string, cfg interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}
