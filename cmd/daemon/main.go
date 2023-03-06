package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/dispatch"
	"github.com/dabump/sonnenbatterie/internal/notification"
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
	sonnenClient := sonnenbatterie.NewClient(&client, cfg)
	shottrDispatcher := dispatch.NewShoutrrrDispatcher(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	notificationChannel := make(chan []*sonnenbatterie.Status)

	sonnenbatterie.NewDeamon(ctx, sonnenClient, cfg, notificationChannel)
	notification.NewDaemon(ctx, cfg, notificationChannel, shottrDispatcher)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	cancel()
	fmt.Print("sonnen batterie deamon stopping...\n")
	time.Sleep(2 * time.Second)
}
