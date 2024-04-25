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
	cfg := config.Config{}

	method, uri, _ := SonnenBatterieStatus(&cfg)
	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "/", uri)
}

type MockSonnenBatterieClient struct {
	status sonnenbatterie.Status
	err    error
}

func (m *MockSonnenBatterieClient) GetStatus(ctx context.Context) (*sonnenbatterie.Status, error) {
	return &m.status, m.err
}

type MockRateLimiter struct {
	hit bool
}

func (m *MockRateLimiter) Hit() bool {
	return m.hit
}

func Test_sbs_sonmnenBatterieController(t *testing.T) {
	cfg := config.Config{}

	type fields struct {
		cfg         *config.Config
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
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(sbc(tt.fields.cfg, tt.fields.sbClient, tt.fields.rateLimiter))
			request, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/", nil)
			resp, err := srv.Client().Do(request)
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

func sbc(cfg *config.Config,
	sbClient sonnenbatterie.SonnenClient, tokenBucket RateLimiter,
) http.HandlerFunc {
	tr := &sbs{
		cfg:         cfg,
		sbClient:    sbClient,
		tokenBucket: tokenBucket,
	}
	return tr.sonnenBatterieController
}
