package controller

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
)

type sbs struct {
	cfg      *config.Config
	ctx      context.Context
	sbClient sonnenbatterie.SonnenClient
}

func SonnenBatterieStatus(ctx context.Context, cfg *config.Config) (string, string, http.HandlerFunc) {
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	sbClient := sonnenbatterie.NewClient(ctx, &client, cfg)

	sbs := sbs{
		cfg:      cfg,
		ctx:      ctx,
		sbClient: sbClient,
	}
	return http.MethodGet, "/", sbs.sonmnenBatterieController
}

func (t *sbs) sonmnenBatterieController(resp http.ResponseWriter, req *http.Request) {
	status, err := t.sbClient.GetStatus()
	if err != nil {
		internalServerError(resp, err)
		return
	}

	pj, err := json.Marshal(status)
	if err != nil {
		internalServerError(resp, err)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(pj))
}

func internalServerError(resp http.ResponseWriter, err error) {
	resp.Header().Set("Content-Type", "application/text")
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(err.Error()))
}
