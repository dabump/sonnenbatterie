package notification

import (
	"fmt"
	"time"

	"github.com/dabump/sonnenbatterie/internal/trend"
)

const (
	lowerThresholdNotification int = 20
	upperThresholdNotification int = 80
)

type ruleEngine struct {
	notifiedOnFull        bool
	notifiedOnEmpty       bool
	lastReset             time.Time
	lastNotificationTrend trend.Trend
}

func newRulesEngine() *ruleEngine {
	return &ruleEngine{
		notifiedOnFull:  false,
		notifiedOnEmpty: false,
		lastReset:       time.Now(),
	}
}

func (r *ruleEngine) dispatchNotification(values []int) bool {
	// Determine initial trend
	t := trend.Calculate(values)
	fmt.Printf("trend: %v - %v%% \n", t, values[0])

	if has24HoursPassed(r.lastReset) {
		r.lastReset = time.Now()
		r.notifiedOnFull = false
		r.notifiedOnEmpty = false
		r.lastNotificationTrend = ""
	}

	// If battery fully charged (100%) and not yet been notified, then enable dispatching
	if t == trend.Upward && values[0] == 100 && !r.notifiedOnFull {
		r.notifiedOnFull = true
		r.notifiedOnEmpty = false
		return true
	}

	// If battery fully drained (0%) and not yet been notified, then enable dispatching
	if t == trend.Downward && values[0] == 0 && !r.notifiedOnEmpty {
		r.notifiedOnFull = false
		r.notifiedOnEmpty = true
		return true
	}

	// Notify if trend is upwards climb past the upward threshold, and has net yet received a notification
	if t == trend.Upward && values[0] >= upperThresholdNotification && r.lastNotificationTrend != trend.Upward {
		r.lastNotificationTrend = t
		return true
	}

	// Notify if trend is downwards falling past the downward threshold, and has net yet received a notification
	if t == trend.Downward && values[0] <= lowerThresholdNotification && r.lastNotificationTrend != trend.Downward {
		r.lastNotificationTrend = t
		return true
	}

	return false

}

func has24HoursPassed(lastChecked time.Time) bool {
	now := time.Now()
	elapsed := now.Sub(lastChecked)
	return elapsed >= 24*time.Hour
}
