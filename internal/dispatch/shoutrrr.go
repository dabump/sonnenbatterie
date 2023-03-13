package dispatch

import (
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/dabump/sonnenbatterie/internal/config"
)

type shoutrrrDispatcher struct {
	sender *router.ServiceRouter
}

func NewShoutrrrDispatcher(cfg *config.Config) *shoutrrrDispatcher {
	router, _ := shoutrrr.CreateSender(cfg.ShoutrrrURL)
	return &shoutrrrDispatcher{
		sender: router,
	}
}

func (s *shoutrrrDispatcher) Send(message string) error {
	err := s.sender.Send(message, &types.Params{})
	if err != nil && err[0] != nil {
		return err[0]
	}
	return nil
}
