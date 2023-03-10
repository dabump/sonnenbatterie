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

type T struct {
	cfg        *config.Config
	ctx        context.Context
	httpClient http.Client
}

func SonnenBatterieStatus(ctx context.Context, cfg *config.Config) (string, string, http.HandlerFunc) {
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	t := T{
		cfg:        cfg,
		ctx:        ctx,
		httpClient: client,
	}
	return http.MethodGet, "/", t.sonmnenBatterieController
}

func (t *T) sonmnenBatterieController(resp http.ResponseWriter, req *http.Request) {
	cl := sonnenbatterie.NewClient(t.ctx, &t.httpClient, t.cfg)
	status, err := cl.GetStatus()
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
