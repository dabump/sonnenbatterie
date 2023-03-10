package router

import (
	"context"
	"net/http"

	"github.com/dabump/sonnenbatterie/internal/config"
	"github.com/dabump/sonnenbatterie/internal/logger"
	"github.com/go-chi/chi/v5"
)

type ControllerFn func(ctx context.Context, cfg *config.Config) (string, string, http.HandlerFunc)

type router struct {
	cfg    *config.Config
	ctx    context.Context
	router chi.Router
}

func New(ctx context.Context, cfg *config.Config) *router {
	rtr := chi.NewRouter()
	rtr.Use(logger.MiddlewareLogger)

	return &router{
		router: rtr,
		ctx:    ctx,
		cfg:    cfg,
	}
}

func (r *router) AddController(cf ControllerFn) {
	r.router.Method(cf(r.ctx, r.cfg))
}

func (r *router) ListenAndServe(address string) {
	go http.ListenAndServe(address, r.router)
}
