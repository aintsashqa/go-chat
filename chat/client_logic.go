package chat

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ErrWritePumpWrapper = "Client.WritePump"
	ErrReadPumpWrapper  = "Client.ReadPump"
)

func NewClient(hub *Hub, ws *websocket.Conn, logger *logrus.Logger) *Client {
	return &Client{
		hub:       hub,
		ws:        ws,
		broadcast: make(chan []byte),

		logger: logger,
	}
}

func (c *Client) Register() {
	c.hub.register <- c
}

func (c *Client) WritePump() {
	defer func() {
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.broadcast:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.Error(errors.Wrap(err, ErrWritePumpWrapper))
				return
			}

			w.Write(message)

			length := len(c.broadcast)
			for i := 0; i < length; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.broadcast)
			}

			if err = w.Close(); err != nil {
				c.logger.Error(errors.Wrap(err, ErrWritePumpWrapper))
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.logger.Error(errors.Wrap(err, ErrReadPumpWrapper))
			break
		}
		c.hub.broadcast <- message
	}
}
