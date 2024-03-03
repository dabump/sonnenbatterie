package dispatch

import (
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type shoutrrrDispatcher struct {
	senders []*router.ServiceRouter
}

func NewShoutrrrDispatcher(shoutrrrUrl ...string) *shoutrrrDispatcher {
	senders := make([]*router.ServiceRouter, len(shoutrrrUrl))
	for index, url := range shoutrrrUrl {
		sender, _ := shoutrrr.CreateSender(url)
		senders[index] = sender
	}
	return &shoutrrrDispatcher{
		senders: senders,
	}
}

func (s *shoutrrrDispatcher) Send(message string) error {
	for _, sender := range s.senders {
		err := sender.Send(message, &types.Params{})
		if err != nil && err[0] != nil {
			return err[0]
		}
	}
	return nil
}
