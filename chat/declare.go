package chat

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Hub struct {
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
	Process([]byte) []byte
	GetMessageCollection() [][]byte
}

type MessageRepositoryInterface interface {
	AddMessage(*Message) error
	GetMessageCollection() ([]Message, error)
}
