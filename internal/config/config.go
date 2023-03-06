package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultConfigFile = "config.cfg"
)

type Config struct {
	SonnenBatterieIP                string `json:"ipAddress"`
	SonnenBatterieStatusPath        string `json:"statusPath"`
	SonnenBatterieProtocolScheme    string `json:"protocolScheme"`
	SonnenBatteriePollingTimeInMins uint   `json:"pollingTimeInMinutes"`

	ShoutrrrURL          string `json:"shoutrrrUrl`
	HttpTimeoutInMinutes uint   `json:"timeoutInMinutes"`
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
		return nil, fmt.Errorf("error during unmarshalling of config file: %v", err)
	}

	return &cfg, nil
}

func LoadConfig() *Config {
	args := os.Args[1:]
	if len(args) >= 1 {
		cfg, err := readConfigFile(args[0], os.ReadFile)
		if err != nil {
			fmt.Printf("config file not located %s\n", args[0])
			panic(err)
		}
		return cfg
	} else {
		cfg, err := readConfigFile(defaultConfigFile, os.ReadFile)
		if err != nil {
			fmt.Printf("default config file not located ./%s\n", defaultConfigFile)
			panic(err)
		}
		return cfg
	}
}
