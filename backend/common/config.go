package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads a YAML configuration file into the provided struct.
func LoadConfig(path string, cfg interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}
