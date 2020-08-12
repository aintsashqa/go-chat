package api

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ErrJoinChatWrapper = "HttpHandler.JoinChat"
)

type httpHandler struct {
	templateDirectory string

	logger *logrus.Logger
}

func NewHttpHandler(templateDirectory string, logger *logrus.Logger) HttpHandlerInterface {
	return &httpHandler{
		templateDirectory: templateDirectory,

		logger: logger,
	}
}

func (h *httpHandler) JoinChat(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("%s/index.html", h.templateDirectory)

	temp := template.Must(template.ParseFiles(filename))
	if err := temp.Execute(w, nil); err != nil {
		h.logger.Error(errors.Wrap(err, ErrJoinChatWrapper))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
