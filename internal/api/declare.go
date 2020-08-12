package api

import (
	"net/http"
)

type HttpHandlerInterface interface {
	JoinChat(http.ResponseWriter, *http.Request)
}

type WSHandlerInterface interface {
	ReceiveMessage(http.ResponseWriter, *http.Request)
}
