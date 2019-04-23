package tail

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/viant/toolbox/url"
	"os"
)

const defaultBatchSize = 512

type Config struct {
	Threads      int
	BatchSize    int
	ProcessedURL string
	Routes       Routes
}

func NewConfigFromEnv(key string) (*Config, error) {
	if key == "" {
		return nil, fmt.Errorf("os config key was empty %v", key)
	}
	config := &Config{}
	encodedConfig := os.Getenv(key)
	if encodedConfig == "" {
		return nil, fmt.Errorf("os config was empty %v", key)
	}
	err := json.Unmarshal([]byte(encodedConfig), config)
	if err == nil {
		if err = config.Init(); err == nil {
			err = config.Validate()
		}
	}
	return config, err
}

func NewConfigFromURL(URL string) (*Config, error) {
	config := &Config{}
	resource := url.NewResource(URL)
	err := resource.Decode(config)
	if err == nil {
		if err = config.Init(); err == nil {
			err = config.Validate()
		}
	}
	return config, err
}

func (c *Config) Init() error {
	if c.Threads == 0 {
		c.Threads = 1
	}
	if c.BatchSize == 0 {
		c.BatchSize = defaultBatchSize
	}
	if c.ProcessedURL != "" {
		c.ProcessedURL = url.NewResource(c.ProcessedURL).URL
	}
	return nil
}

func (c *Config) Validate() error {
	if len(c.Routes) == 0 {
		return errors.New("routes were empty")
	}
	if c.ProcessedURL == "" {
		return errors.New("processedURL was empty")
	}
	return nil
}
