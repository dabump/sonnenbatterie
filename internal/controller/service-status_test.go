package controller

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestServiceStatus(t *testing.T) {
	cfg := config.Config{}

	method, uri, _ := ServiceStatus(&cfg)

	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "/status", uri)
}

func TestServiceStatusHandler(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(schf())
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/status", nil)
	resp, err := srv.Client().Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(bodyBytes))

	ctHeader := resp.Header.Get("Content-Type")
	assert.Equal(t, "application/text", ctHeader)
}

func schf() http.HandlerFunc {
	return statusController
}
