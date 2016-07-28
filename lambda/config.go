package lambda

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type configLoader struct {
	ConfigURL string `json:"config_url"`
}

// ReadFromFile reads a file from disk
func (c *configLoader) ReadFromFile(name string) error {
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
func LoadConfig(name string) (*configLoader, error) {
	config := new(configLoader)
	if err := config.ReadFromFile(name); err != nil {
		return nil, err
	}

	return config, nil
}

// Config holder object
type Config struct {
	Bucket          string `json:"bucket"`
	EnvironmentName string `json:"environment_name"`
	KMSKeyID        string `json:"kms_key_id"`
}

// ReadFromFile reads a file from disk
func (c *Config) ReadFromReader(r io.Reader) error {
	config, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(config, c); err != nil {
		return err
	}

	return nil
}
