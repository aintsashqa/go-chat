package chat

import (
	"github.com/sirupsen/logrus"
)

func NewHub(handler MessageServiceInterface, logger *logrus.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),

		handler: handler,
		logger:  logger,
	}
}

func (h *Hub) registerClient(client *Client) {
	h.logger.Debugf("Hub register new client [%v]\n", client)
	h.clients[client] = true
	h.logger.Debugf("Hub triyng to send all messages from database to current client [%v]", client)
	messages := h.handler.GetMessageCollection()
	for _, m := range messages {
		client.broadcast <- m
	}
}

func (h *Hub) broadcastMessageToClients(message []byte) {
	result := h.handler.Process(message)
	for client := range h.clients {
		client.broadcast <- result
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			delete(h.clients, client)

		case message := <-h.broadcast:
			h.broadcastMessageToClients(message)
		}
	}
}
