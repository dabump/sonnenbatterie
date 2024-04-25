package sonnenbatterie

import (
	"context"
	"time"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/logger"
	"github.com/dabump/sonnenbatterie/internal/queue"
)

const (
	reverseLimitQueueSize                    = 30
	defaultTickerTimeInSecords time.Duration = 2 * time.Second
)

type SonnenClient interface {
	GetStatus(ctx context.Context) (*Status, error)
}

type Daemon struct {
	config              *config.Config
	sonnenClient        SonnenClient
	lastStatusCheck     time.Time
	notificationChannel chan []*Status
	reverseLimitQueue   *queue.ReversedLimitedQueue
}

func NewDeamon(ctx context.Context, sonnenClient SonnenClient,
	config *config.Config, notificationChannel chan []*Status,
) *Daemon {
	daemon := Daemon{
		config:              config,
		sonnenClient:        sonnenClient,
		notificationChannel: notificationChannel,
		reverseLimitQueue:   queue.NewReversedLimitedQueue(reverseLimitQueueSize),
	}

	daemon.start(ctx)
	return &daemon
}

func (d *Daemon) start(ctx context.Context) {
	d.lastStatusCheck = time.Now()
	ticker := time.NewTicker(defaultTickerTimeInSecords)
	pollingTime := time.Duration(d.config.SonnenBatteriePollingTimeInMins) * time.Minute

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.LoggerFromContext(ctx).Info("sonnen batterie daemon stopped")
				return
			case <-ticker.C:
				if d.lastStatusCheck.Add(pollingTime).Before(time.Now()) {
					d.lastStatusCheck = time.Now()
					status, err := d.sonnenClient.GetStatus(ctx)
					if err != nil {
						logger.LoggerFromContext(ctx).Errorf("error during communication to sonnen batterie: %v", err)
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
