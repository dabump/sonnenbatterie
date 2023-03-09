package notification

import (
	"context"
	"fmt"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
	"github.com/dabump/sonnenbatterie/internal/trend"
)

type MessageDispatcher interface {
	Send(message string) []error
}

type notificationEngine struct {
	config              *config.Config
	context             context.Context
	dispatcher          MessageDispatcher
	notificationChannel <-chan []*sonnenbatterie.Status
}

func NewDaemon(ctx context.Context, cfg *config.Config,
	chn <-chan []*sonnenbatterie.Status, dispatcher MessageDispatcher,
) *notificationEngine {
	ne := notificationEngine{
		config:              cfg,
		context:             ctx,
		dispatcher:          dispatcher,
		notificationChannel: chn,
	}

	ne.start()
	return &ne
}

func (n *notificationEngine) start() {
	rulesEngine := NewRulesEngine()
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
					err := n.dispatcher.Send(message)
					if err != nil {
						fmt.Printf("err: %v\n", err)
					}
				}
			}
		}
	}()
}
