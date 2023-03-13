package controller

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dabump/sonnenbatterie/internal/common"
	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
	"github.com/dabump/tokenbucket"
)

const (
	ratePerSecondLimit = 1
)

type RateLimiter interface {
	Hit() bool
}

type sbs struct {
	cfg         *config.Config
	ctx         context.Context
	sbClient    sonnenbatterie.SonnenClient
	tokenBucket RateLimiter
}

func SonnenBatterieStatus(ctx context.Context, cfg *config.Config) (string, string, http.HandlerFunc) {
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	sbClient := sonnenbatterie.NewClient(ctx, &client, cfg)

	tokenBucket := tokenbucket.NewBucket("sonnen-status", ratePerSecondLimit)
	tokenBucketDaemon := tokenbucket.NewDaemon(tokenBucket, tokenbucket.NA)
	tokenBucketDaemon.Start()

	sbs := sbs{
		cfg:         cfg,
		ctx:         ctx,
		sbClient:    sbClient,
		tokenBucket: tokenBucketDaemon,
	}
	return http.MethodGet, "/", sbs.sonmnenBatterieController
}

func (t *sbs) sonmnenBatterieController(resp http.ResponseWriter, req *http.Request) {
	hit := t.tokenBucket.Hit()
	if !hit {
		common.TooManyRequests(resp)
		return
	}

	status, err := t.sbClient.GetStatus()
	if err != nil {
		common.InternalServerError(resp, err)
		return
	}

	pj, err := json.Marshal(status)
	if err != nil {
		common.InternalServerError(resp, err)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(pj))
}
