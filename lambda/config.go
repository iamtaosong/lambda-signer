package lambda

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config holder object
type Config struct {
	Bucket          string `json:"bucket"`
	EnvironmentName string `json:"environment_name"`
	KMSKeyID        string `json:"kms_key_id"`
}

// ReadFromFile reads a file from disk
func (c *Config) ReadFromFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	defer f.Close()

	config, err := ioutil.ReadAll(f)

	if err := json.Unmarshal(config, c); err != nil {
		return err
	}

	return nil
}

// LoadConfig from a file on disk
func LoadConfig(name string) (*Config, error) {
	config := new(Config)
	if err := config.ReadFromFile(name); err != nil {
		return nil, err
	}

	return config, nil
}
