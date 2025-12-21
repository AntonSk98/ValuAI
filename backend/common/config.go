package common

import (
	"os"

	"github.com/drone/envsubst/v2"
	"gopkg.in/yaml.v3"
)

// LoadProperties loads YAML configuration from the specified file path into the provided struct.
// It automatically expands environment variables in the format ${VAR_NAME} or $VAR_NAME.
func LoadProperties(path string, cfg interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Expand environment variables in the YAML content
	expandedData, err := envsubst.Eval(string(data), func(key string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		// If environment variable is not set, keep the placeholder
		return "${" + key + "}"
	})
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(expandedData), cfg)
}
