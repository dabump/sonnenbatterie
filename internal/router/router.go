package router

import (
	"net/http"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/logger"
	"github.com/go-chi/chi/v5"
)

type ControllerFn func(cfg *config.Config) (string, string, http.HandlerFunc)

type router struct {
	cfg    *config.Config
	router chi.Router
}

func New(cfg *config.Config) *router {
	rtr := chi.NewRouter()
	rtr.Use(logger.MiddlewareLogger)

	return &router{
		router: rtr,
		cfg:    cfg,
	}
}

func (r *router) AddController(cf ControllerFn) {
	r.router.Method(cf(r.cfg))
}

func (r *router) ListenAndServe(address string) {
	go func() {
		_ = http.ListenAndServe(address, r.router)
	}()
}
