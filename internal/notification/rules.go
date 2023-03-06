package notification

import (
	"fmt"

	"github.com/dabump/sonnenbatterie/internal/trend"
)

const (
	lowerThresholdNotification int = 51
	upperThresholdNotification int = 53
)

type ruleEngine struct {
	notifiedOnFull        bool
	notifiedOnEmpty       bool
	lastNotificationTrend trend.Trend
}

func newRulesEngine() *ruleEngine {
	return &ruleEngine{
		notifiedOnFull: false,
		notifiedOnEmpty: false,
	}
}

func (r *ruleEngine) dispatchNotification(values []int) bool {
	// Determine initial trend
	t := trend.Calculate(values)
	if r.lastNotificationTrend != trend.Upward &&
		r.lastNotificationTrend != trend.Downward &&
		r.lastNotificationTrend != trend.NoTrend {
		r.lastNotificationTrend = t
	}
	fmt.Printf("trend: %v - %v%% \n", t, values[0])

	// If battery fully charged (100%) and not yet been notified, then enable dispatching
	if t == trend.Upward && values[0] == 100 && !r.notifiedOnFull {
		r.notifiedOnFull = true
		r.notifiedOnEmpty = false
		return true
	}

	// If battery fully frained (0%) and not yet been notified, then enable dispatching
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
