package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
)

func main() {
	cfg := config.LoadConfig()

	client := http.Client{
		Timeout: time.Duration(cfg.HttpTimeoutInMinutes) * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	sonnenBatterieClient := sonnenbatterie.NewClient(&client, cfg)
	status, err := sonnenBatterieClient.GetStatus()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	ms, _ := json.MarshalIndent(status, "", " ")
	fmt.Printf("%s\n", ms)
}
