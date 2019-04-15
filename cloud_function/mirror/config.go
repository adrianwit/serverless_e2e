package split

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//Config reprsents service config
type Config struct {
	DestURL string
	Secrets string
	SecretURL string//TODO load secrsts from URL is non empty
	SecretKey string
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