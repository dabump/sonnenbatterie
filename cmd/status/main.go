package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/logger"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
)

func main() {
	ctx := context.Background()
	ctx = logger.NewContextLogger(ctx)

	cfg := config.LoadConfig(ctx)

	client := http.Client{
		Timeout: time.Duration(cfg.HttpTimeoutInMinutes) * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	sonnenBatterieClient := sonnenbatterie.NewClient(ctx, &client, cfg)
	status, err := sonnenBatterieClient.GetStatus()
	if err != nil {
		logger.LoggerFromContext(ctx).Errorf("error: %v", err)
		os.Exit(1)
	}
	ms, _ := json.MarshalIndent(status, "", " ")
	fmt.Printf("%s\n", ms)
}
