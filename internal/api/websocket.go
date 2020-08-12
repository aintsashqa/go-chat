package api

import (
	"net/http"

	"github.com/aintsashqa/go-chat/chat"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ErrUpgradeWrapper        = "WSHandler.Upgrade"
	ErrReceiveMessageWrapper = "WSHandler.ReceiveMessage"
)

var (
	upgrader = &websocket.Upgrader{}
)

type wsHandler struct {
	hub    *chat.Hub
	logger *logrus.Logger
}

func NewWSHandler(hub *chat.Hub, logger *logrus.Logger) WSHandlerInterface {
	return &wsHandler{
		hub:    hub,
		logger: logger,
	}
}

func (h *wsHandler) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrUpgradeWrapper)
	}
	return ws, nil
}

func (h *wsHandler) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrade(w, r)
	if err != nil {
		h.logger.Error(errors.Wrap(err, ErrReceiveMessageWrapper))
		return
	}

	client := chat.NewClient(h.hub, ws, h.logger)
	client.Register()

	go client.WritePump()
	go client.ReadPump()
}
