package aeroagg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	IP        string
	Port      int
	Namespace string
	Dataset   string
}

//NewConfigFromEnv creates a new config from env
func NewConfigFromEnv(key string) (*Config, error) {
	payload := os.Getenv(key)
	if payload == "" {
		return nil, fmt.Errorf("config payload was empty")
	}
	config := &Config{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(config)
	return config, err
}
