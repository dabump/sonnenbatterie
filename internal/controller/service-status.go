package controller

import (
	"net/http"

	"github.com/dabump/sonnenbatterie/internal/config"
)

func ServiceStatus(*config.Config) (string, string, http.HandlerFunc) {
	return http.MethodGet, "/status", statusController
}

func statusController(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/text")
	resp.WriteHeader(http.StatusOK)
	_,_ = resp.Write([]byte("ok"))
}
