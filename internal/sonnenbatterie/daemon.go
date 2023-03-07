package sonnenbatterie

import (
	"context"
	"fmt"
	"time"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/queue"
)

const (
	reverseLimitQueueSize = 30
	defaultTickerTimeInSecords time.Duration = 2 * time.Second
)

type SonnenClient interface {
	GetStatus() (*Status, error)
}

type Daemon struct {
	ctx                 context.Context
	config              *config.Config
	sonnenClient        SonnenClient
	lastStatusCheck     time.Time
	notificationChannel chan []*Status
	reverseLimitQueue   *queue.ReversedLimitedQueue
}

func NewDeamon(ctx context.Context, sonnenClient SonnenClient,
	config *config.Config, notificationChannel chan []*Status) *Daemon {

	daemon := Daemon{
		ctx:                 ctx,
		config:              config,
		sonnenClient:        sonnenClient,
		notificationChannel: notificationChannel,
		reverseLimitQueue:   queue.NewReversedLimitedQueue(reverseLimitQueueSize),
	}

	daemon.start()
	return &daemon
}

func (d *Daemon) start() {
	d.lastStatusCheck = time.Now()
	ticker := time.NewTicker(defaultTickerTimeInSecords)
	pollingTime := time.Duration(d.config.SonnenBatteriePollingTimeInMins) * time.Minute

	go func() {
		for {
			select {
			case <-d.ctx.Done():
				fmt.Print("sonnen batterie daemon stopped\n")
				return
			case <-ticker.C:
				if d.lastStatusCheck.Add(pollingTime).Before(time.Now()) {
					d.lastStatusCheck = time.Now()
					status, err := d.sonnenClient.GetStatus()
					if err != nil {
						fmt.Printf("error during communication to sonnen batterie: %v\n", err)
						continue
					}

					d.reverseLimitQueue.Enqueue(status)

					statusItems := make([]*Status, len(d.reverseLimitQueue.Queue()))
					for i, v := range d.reverseLimitQueue.Queue() {
						statusItems[i] = v.(*Status)
					}
					d.notificationChannel <- statusItems
				}
			}
		}
	}()
}
