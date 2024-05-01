package notification

import (
	"context"
	"fmt"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/logger"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
	"github.com/dabump/sonnenbatterie/internal/trend"
)

type MessageDispatcher interface {
	Send(message string) error
}

type notificationEngine struct {
	config              *config.Config
	dispatcher          MessageDispatcher
	notificationChannel <-chan []*sonnenbatterie.Status
}

func NewDaemon(ctx context.Context, cfg *config.Config,
	chn <-chan []*sonnenbatterie.Status, dispatcher MessageDispatcher,
) *notificationEngine {
	ne := notificationEngine{
		config:              cfg,
		dispatcher:          dispatcher,
		notificationChannel: chn,
	}

	ne.start(ctx)
	return &ne
}

func (n *notificationEngine) start(ctx context.Context) {
	rulesEngine := NewRulesEngine()
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.LoggerFromContext(ctx).Info("notification engine stopped")
				return
			case event := <-n.notificationChannel:

				var values []float64
				for _, s := range event {
					values = append(values, s.Usoc)
				}
				if rulesEngine.dispatchNotification(ctx, values) {
					var message string
					if values[0] == 100 {
						message = "sonnenbatterie fully charged"
					} else if values[0] == 0 {
						message = "sonnenbatterie depleted"
					} else {
						message = fmt.Sprintf("sonnenbatterie at %v%% with %s trend", values[0], trend.Calculate(values))
					}
					logger.LoggerFromContext(ctx).Infof("message: %v", message)

					err := n.dispatcher.Send(message)
					if err != nil {
						logger.LoggerFromContext(ctx).Errorf("err: %v", err)
					}
				}
			}
		}
	}()
}
