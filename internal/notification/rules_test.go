package notification

import (
	"context"
	"testing"
	"time"

	"github.com/dabump/sonnenbatterie/internal/trend"
)

func Test_ruleEngine_dispatchNotification(t *testing.T) {
	twentyFourHourDuration, _ := time.ParseDuration("-24h")
	before24Hours := time.Now().Add(twentyFourHourDuration)

	type fields struct {
		notifiedOnFull        bool
		notifiedOnEmpty       bool
		lastReset             time.Time
		lastNotificationTrend trend.Trend
	}
	type args struct {
		values []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "100% not yet notified",
			fields: fields{
				notifiedOnFull:        false,
				notifiedOnEmpty:       false,
				lastReset:             time.Now(),
				lastNotificationTrend: trend.Upward,
			},
			want: true,
			args: args{
				values: []int{100, 99},
			},
		},
		{
			name: "100% already notified",
			fields: fields{
				notifiedOnFull:        true,
				notifiedOnEmpty:       false,
				lastReset:             time.Now(),
				lastNotificationTrend: trend.Upward,
			},
			want: false,
			args: args{
				values: []int{100, 99},
			},
		},
		{
			name: "100% already notified - 24 hours since last notification",
			fields: fields{
				notifiedOnFull:        true,
				notifiedOnEmpty:       false,
				lastReset:             before24Hours,
				lastNotificationTrend: trend.Upward,
			},
			want: true,
			args: args{
				values: []int{100, 99},
			},
		},
		{
			name: "0% not yet notified",
			fields: fields{
				notifiedOnFull:        false,
				notifiedOnEmpty:       false,
				lastReset:             time.Now(),
				lastNotificationTrend: trend.Downward,
			},
			want: true,
			args: args{
				values: []int{0, 1, 2},
			},
		},
		{
			name: "0% already notified",
			fields: fields{
				notifiedOnFull:        false,
				notifiedOnEmpty:       true,
				lastReset:             time.Now(),
				lastNotificationTrend: trend.Downward,
			},
			want: false,
			args: args{
				values: []int{0, 1, 2},
			},
		},
		{
			name: "0% already notified - 24 hours since last notification",
			fields: fields{
				notifiedOnFull:        false,
				notifiedOnEmpty:       true,
				lastReset:             before24Hours,
				lastNotificationTrend: trend.Downward,
			},
			want: true,
			args: args{
				values: []int{0, 1, 2},
			},
		},
		{
			name: "Upper threshold reacher",
			fields: fields{
				notifiedOnFull:        false,
				notifiedOnEmpty:       false,
				lastNotificationTrend: trend.Downward,
			},
			want: true,
			args: args{
				values: []int{upperThresholdNotification, upperThresholdNotification - 1, upperThresholdNotification - 2},
			},
		},
		{
			name: "Lower threshold reacher",
			fields: fields{
				notifiedOnFull:        false,
				notifiedOnEmpty:       false,
				lastNotificationTrend: trend.Upward,
			},
			want: true,
			args: args{
				values: []int{lowerThresholdNotification, lowerThresholdNotification + 1, lowerThresholdNotification + 2},
			},
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			r := &ruleEngine{
				notifiedOnFull:        tt.fields.notifiedOnFull,
				notifiedOnEmpty:       tt.fields.notifiedOnEmpty,
				lastNotificationTrend: tt.fields.lastNotificationTrend,
				lastReset:             tt.fields.lastReset,
			}
			if got := r.dispatchNotification(ctx, tt.args.values); got != tt.want {
				t.Errorf("ruleEngine.dispatchNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}
