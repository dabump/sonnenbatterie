package controller

import (
	"context"
	"net/http"

	"github.com/dabump/sonnenbatterie/internal/config"
)

func ServiceStatus(context.Context, *config.Config) (string, string, http.HandlerFunc) {
	return http.MethodGet, "/status", statusController
}

func statusController(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/text")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}
