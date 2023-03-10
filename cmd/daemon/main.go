package main

import (
	"context"
	"crypto/tls"
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
	shottrDispatcher := dispatch.NewShoutrrrDispatcher(cfg)
	err := shottrDispatcher.Send("sonnenbatterie daemon started")
	if len(err) >= 1 && err[0] != nil {
		logger.LoggerFromContext(ctx).Errorf("could not invoke notification dispatcher err: %v", err)
	}

	client := http.Client{
		Timeout: time.Duration(cfg.HttpTimeoutInMinutes) * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	sonnenClient := sonnenbatterie.NewClient(ctx, &client, cfg)

	notificationChannel := make(chan []*sonnenbatterie.Status)

	sonnenbatterie.NewDeamon(ctx, sonnenClient, cfg, notificationChannel)
	notification.NewDaemon(ctx, cfg, notificationChannel, shottrDispatcher)

	rtr := router.New(ctx, cfg)
	rtr.AddController(controller.ServiceStatus)
	rtr.AddController(controller.SonnenBatterieStatus)
	rtr.ListenAndServe(":8888")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	cancel()

	logger.LoggerFromContext(ctx).Info("sonnen batterie deamon stopping")
	shottrDispatcher.Send("sonnenbatterie daemon stopped")
	time.Sleep(2 * time.Second)
}
