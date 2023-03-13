package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/sonnenbatterie"
	"github.com/stretchr/testify/assert"
)

func TestSonnenBatterieStatus(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{}

	method, uri, _ := SonnenBatterieStatus(ctx, &cfg)
	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "/", uri)
}

type MockSonnenBatterieClient struct {
	status sonnenbatterie.Status
	err    error
}

func (m *MockSonnenBatterieClient) GetStatus() (*sonnenbatterie.Status, error) {
	return &m.status, m.err
}

type MockRateLimiter struct {
	hit bool
}

func (m *MockRateLimiter) Hit() bool {
	return m.hit
}

func Test_sbs_sonmnenBatterieController(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{}

	type fields struct {
		cfg         *config.Config
		ctx         context.Context
		sbClient    sonnenbatterie.SonnenClient
		rateLimiter RateLimiter
	}
	tests := []struct {
		name              string
		fields            fields
		wantErr           bool
		statusCode        int
		contentTypeHeader string
	}{
		{
			name: "success-200-ok",
			fields: fields{
				cfg: &cfg,
				ctx: ctx,
				sbClient: &MockSonnenBatterieClient{
					status: sonnenbatterie.Status{},
					err:    nil,
				},
				rateLimiter: &MockRateLimiter{
					hit: true,
				},
			},
			wantErr:           false,
			statusCode:        http.StatusOK,
			contentTypeHeader: "application/json",
		},
		{
			name: "unsuccessful-5xx-ise",
			fields: fields{
				cfg: &cfg,
				ctx: ctx,
				sbClient: &MockSonnenBatterieClient{
					status: sonnenbatterie.Status{},
					err:    fmt.Errorf("somethingWentWrong"),
				},
				rateLimiter: &MockRateLimiter{
					hit: true,
				},
			},
			wantErr:           false,
			statusCode:        http.StatusInternalServerError,
			contentTypeHeader: "application/text",
		},
		{
			name: "unsuccessful-429-ratelimit",
			fields: fields{
				cfg: &cfg,
				ctx: ctx,
				sbClient: &MockSonnenBatterieClient{
					status: sonnenbatterie.Status{},
					err:    fmt.Errorf("somethingWentWrong"),
				},
				rateLimiter: &MockRateLimiter{
					hit: false,
				},
			},
			wantErr:           false,
			statusCode:        http.StatusTooManyRequests,
			contentTypeHeader: "application/text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(sbc(tt.fields.ctx, tt.fields.cfg, tt.fields.sbClient, tt.fields.rateLimiter))
			resp, err := srv.Client().Get(srv.URL + "/")
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.statusCode != resp.StatusCode {
				t.Errorf("http.Response.StatusCode() got = %v, want %v", resp.StatusCode, tt.statusCode)
				return
			}
			if tt.contentTypeHeader != resp.Header.Get("Content-Type") {
				t.Errorf("invalid content-type header value got = %v, want %v", resp.Header.Get("Content-Type"), tt.contentTypeHeader)
				return
			}
		})
	}
}

func sbc(ctx context.Context, cfg *config.Config,
	sbClient sonnenbatterie.SonnenClient, tokenBucket RateLimiter,
) http.HandlerFunc {
	tr := &sbs{
		cfg:         cfg,
		ctx:         ctx,
		sbClient:    sbClient,
		tokenBucket: tokenBucket,
	}
	return tr.sonmnenBatterieController
}
