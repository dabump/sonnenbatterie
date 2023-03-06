package notification

import (
	"context"
	"fmt"

	"github.com/containrrr/shoutrrr"
	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
	"github.com/dabump/sonnenbatterie/internal/trend"
)

type notificationEngine struct {
	config              *config.Config
	context             context.Context
	notificationChannel <-chan []*sonnenbatterie.Status
}

func NewDaemon(ctx context.Context, cfg *config.Config,
	chn <-chan []*sonnenbatterie.Status) *notificationEngine {

	ne := notificationEngine{
		config:              cfg,
		context:             ctx,
		notificationChannel: chn,
	}

	ne.start()
	return &ne
}

func (n *notificationEngine) start() {
	rulesEngine := newRulesEngine()
	go func() {
		for {
			select {
			case <-n.context.Done():
				fmt.Print("notification engine stopped\n")
				return
			case event := <-n.notificationChannel:

				var values []int
				for _, s := range event {
					values = append(values, s.Usoc)
				}
				if rulesEngine.dispatchNotification(values) {
					var message string
					if values[0] == 100 {
						message = "sonnenbatterie fully charged"
					} else if values[0] == 0 {
						message = "sonnenbatterie depleted"
					} else {
						message = fmt.Sprintf("sonnenbatterie at %v%% with %s trend", values[0], trend.Calculate(values))
					}
					fmt.Printf("message: %v\n", message)
					err := shoutrrr.Send(n.config.ShoutrrrURL, message)
					if err != nil {
						fmt.Printf("err: %v\n", err)
					}
				}
			}
		}
	}()
}
