package chat

import (
	"errors"

	erro "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	WrapperCommunityRegisterMethod      = "Community.Register"
	WrapperCommunityFindHubWithIdMethod = "Community.FindHubWithId"
	WrapperCommunityUnregisterHubMethod = "Community.UnregisterHub"
)

var (
	ErrHubNotFound = errors.New("Hub not found.")
)

func NewCommunity(handler MessageServiceInterface, logger *logrus.Logger) *Community {
	return &Community{
		hubs: make(map[*Hub]bool),

		handler: handler,
		logger:  logger,
	}
}

func (c *Community) RegisterHub() *Hub {
	hub := NewHub(c.handler, c.logger)
	go hub.Run()
	c.hubs[hub] = true

	c.logger.Infof("[%s]: Community has created a new hub {%s}.", WrapperCommunityRegisterMethod, hub.Getid().String())

	return hub
}

func (c *Community) FindHubWithId(id string) (*Hub, error) {
	for h := range c.hubs {
		if h.id.String() == id {
			return h, nil
		}
	}

	return nil, erro.Wrap(ErrHubNotFound, WrapperCommunityFindHubWithIdMethod)
}

func (c *Community) UnregisterHub(id string) error {
	hub, err := c.FindHubWithId(id)
	if err != nil {
		return erro.Wrap(err, WrapperCommunityUnregisterHubMethod)
	}

	hub.Terminate()
	delete(c.hubs, hub)

	return nil
}
