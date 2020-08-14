package chat

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Community struct {
	hubs map[*Hub]bool

	handler MessageServiceInterface
	logger  *logrus.Logger
}

type Hub struct {
	id         uuid.UUID
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	handler MessageServiceInterface
	logger  *logrus.Logger
}

type Client struct {
	hub       *Hub
	ws        *websocket.Conn
	broadcast chan []byte

	logger *logrus.Logger
}

type Message struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type MessageServiceInterface interface {
	Process(uuid.UUID, []byte) []byte
	GetMessageCollection(uuid.UUID) [][]byte
}

type MessageRepositoryInterface interface {
	AddMessage(string, *Message) error
	GetMessageCollection(string) ([]Message, error)
}
