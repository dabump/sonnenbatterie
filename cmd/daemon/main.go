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
	"github.com/dabump/sonnenbatterie/internal/controller"
	"github.com/dabump/sonnenbatterie/internal/dispatch"
	"github.com/dabump/sonnenbatterie/internal/logger"
	"github.com/dabump/sonnenbatterie/internal/notification"
	"github.com/dabump/sonnenbatterie/internal/router"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.NewContextLogger(ctx)

	cfg := config.LoadConfig(ctx)

	shottrDispatcher := dispatch.NewShoutrrrDispatcher(cfg.ShoutrrrURLs...)
	err := shottrDispatcher.Send("sonnenbatterie daemon started...")
	if err != nil {
		logger.LoggerFromContext(ctx).Errorf("could not invoke notification dispatcher err: %v", err)
	}

	client := http.Client{
		Timeout: time.Duration(cfg.HttpTimeoutInMinutes) * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	sonnenClient := sonnenbatterie.NewClient(&client, cfg)

	notificationChannel := make(chan []*sonnenbatterie.Status)

	sonnenbatterie.NewDeamon(ctx, sonnenClient, cfg, notificationChannel)
	notification.NewDaemon(ctx, cfg, notificationChannel, shottrDispatcher)

	rtr := router.New(cfg)
	rtr.AddController(controller.ServiceStatus)
	rtr.AddController(controller.SonnenBatterieStatus)
	rtr.ListenAndServe(":" + fmt.Sprint(cfg.HttpServerPort))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	cancel()

	logger.LoggerFromContext(ctx).Info("sonnen batterie deamon stopping")
	err = shottrDispatcher.Send("sonnenbatterie daemon stopped")
	if err != nil {
		logger.LoggerFromContext(ctx).Errorf("unable to send message via shottr: %w", err)
	}
	time.Sleep(2 * time.Second)
}
