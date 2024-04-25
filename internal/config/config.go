package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dabump/sonnenbatterie/internal/logger"
)

const (
	defaultConfigFile = "config.cfg"
)

type Config struct {
	SonnenBatterieIP                string `json:"ipAddress"`
	SonnenBatterieStatusPath        string `json:"statusPath"`
	SonnenBatterieProtocolScheme    string `json:"protocolScheme"`
	SonnenBatteriePollingTimeInMins uint   `json:"pollingTimeInMinutes"`

	ShoutrrrURLs         []string `json:"shoutrrrUrls"`
	HttpTimeoutInMinutes uint     `json:"timeoutInMinutes"`
	HttpServerPort       uint     `json:"httpServerPort"`
}

type reader func(name string) ([]byte, error)

func readConfigFile(fileName string, fn reader) (*Config, error) {
	var cfg Config
	b, err := fn(fileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error during unmarshalling of config file: %w", err)
	}

	return &cfg, nil
}

func LoadConfig(ctx context.Context) *Config {
	args := os.Args[1:]
	if len(args) >= 1 {
		cfg, err := readConfigFile(args[0], os.ReadFile)
		if err != nil {
			logger.LoggerFromContext(ctx).Errorf("config file not located %s", args[0])
			panic(err)
		}
		return cfg
	} else {
		cfg, err := readConfigFile(defaultConfigFile, os.ReadFile)
		if err != nil {
			logger.LoggerFromContext(ctx).Errorf("default config file not located ./%s", defaultConfigFile)
			panic(err)
		}
		return cfg
	}
}
