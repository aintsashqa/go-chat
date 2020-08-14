package chat

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	WrapperHubRunMehtod = "Hub.Run"
)

func NewHub(handler MessageServiceInterface, logger *logrus.Logger) *Hub {
	return &Hub{
		id:         uuid.NewV4(),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),

		handler: handler,
		logger:  logger,
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
	messages := h.handler.GetMessageCollection(h.id)
	for _, m := range messages {
		client.broadcast <- m
	}
}

func (h *Hub) broadcastMessageToClients(message []byte) {
	result := h.handler.Process(h.id, message)
	for client := range h.clients {
		client.broadcast <- result
	}
}

func (h *Hub) Getid() uuid.UUID {
	return h.id
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.logger.Debugf("[%s]: Hub {%s} register a new client.", WrapperHubRunMehtod, h.id.String())
			h.registerClient(client)

		case client := <-h.unregister:
			h.logger.Debugf("[%s]: Hub {%s} has delete a client.", WrapperHubRunMehtod, h.id.String())
			delete(h.clients, client)

		case message := <-h.broadcast:
			h.broadcastMessageToClients(message)
		}
	}
}

func (h *Hub) Terminate() {
	for c := range h.clients {
		c.Terminate()
	}
}
