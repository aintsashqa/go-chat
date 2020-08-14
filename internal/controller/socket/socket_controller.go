package socket

import (
	"net/http"

	"github.com/aintsashqa/go-chat/chat"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	WrapperSocketControllerUpgradeMethod        = "SocketController.Upgrade"
	WrapperSocketControllerReceiveMessageMethod = "SocketController.ReceiveMessage"
)

var (
	upgrader = &websocket.Upgrader{}
)

type SocketController struct {
	community *chat.Community
	logger    *logrus.Logger
}

func NewSocketController(community *chat.Community, logger *logrus.Logger) *SocketController {
	return &SocketController{
		community: community,
		logger:    logger,
	}
}

func (c *SocketController) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.Wrap(err, WrapperSocketControllerUpgradeMethod)
	}
	return ws, nil
}

func (c *SocketController) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("hub_id")
	hub, err := c.community.FindHubWithId(id)
	if err != nil {
		c.logger.Error(errors.Wrap(err, WrapperSocketControllerReceiveMessageMethod))
		return
	}

	ws, err := c.upgrade(w, r)
	if err != nil {
		c.logger.Error(errors.Wrap(err, WrapperSocketControllerReceiveMessageMethod))
		return
	}

	client := chat.NewClient(hub, ws, c.logger)
	client.Register()

	go client.WritePump()
	go client.ReadPump()
}
