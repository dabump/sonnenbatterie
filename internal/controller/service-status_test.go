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
	ctx := context.Background()
	cfg := config.Config{}

	method, uri, _ := ServiceStatus(ctx, &cfg)

	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "/status", uri)
}

func TestServiceStatusHandler(t *testing.T) {
	srv := httptest.NewServer(schf())
	resp, err := srv.Client().Get(srv.URL + "/status")

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
