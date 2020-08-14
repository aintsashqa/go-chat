package web

import (
	"net/http"

	"github.com/aintsashqa/go-chat/chat"
	"github.com/aintsashqa/go-chat/internal/service/render"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	WrapperHubControllerCreateHubMethod   = "HubController.CreateHub"
	WrapperHubControllerJoinHubMethod     = "HubController.JoinHub"
	WrapperHubControllerInviteToHubMethod = "HubController.InviteToHub"
)

type HubController struct {
	community *chat.Community
	cookie    *sessions.CookieStore
	render    *render.HtmlRenderService
	logger    *logrus.Logger
}

func NewHubController(community *chat.Community, cookie *sessions.CookieStore, render *render.HtmlRenderService, logger *logrus.Logger) *HubController {
	return &HubController{
		community: community,
		cookie:    cookie,
		render:    render,
		logger:    logger,
	}
}

func (c *HubController) NewHub(w http.ResponseWriter, r *http.Request) {
	c.render.Render(w, "hub/new_hub", nil)
}

func (c *HubController) CreateHub(w http.ResponseWriter, r *http.Request) {
	session, err := c.cookie.Get(r, "session")
	if err != nil {
		c.logger.Error(errors.Wrap(err, WrapperHubControllerCreateHubMethod))
		c.render.ErrorResponse(w, http.StatusInternalServerError)
	}

	hub := c.community.RegisterHub()
	session.Values["hub_id"] = hub.Getid().String()
	if err = session.Save(r, w); err != nil {
		c.logger.Error(errors.Wrap(err, WrapperHubControllerCreateHubMethod))
		c.render.ErrorResponse(w, http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/hub", http.StatusMovedPermanently)
}

func (c *HubController) JoinHub(w http.ResponseWriter, r *http.Request) {
	session, err := c.cookie.Get(r, "session")
	if err != nil {
		c.logger.Error(errors.Wrap(err, WrapperHubControllerJoinHubMethod))
		c.render.ErrorResponse(w, http.StatusInternalServerError)
	}

	hubID := session.Values["hub_id"].(string)

	data := struct {
		HubID string
	}{
		HubID: hubID,
	}
	c.render.Render(w, "hub/hub", data)
}

func (c *HubController) InviteToHub(w http.ResponseWriter, r *http.Request) {
	hub_id := chi.URLParam(r, "hub_id")

	session, err := c.cookie.Get(r, "session")
	if err != nil {
		c.logger.Error(errors.Wrap(err, WrapperHubControllerInviteToHubMethod))
		c.render.ErrorResponse(w, http.StatusInternalServerError)
	}

	session.Values["hub_id"] = hub_id
	if err = session.Save(r, w); err != nil {
		c.logger.Error(errors.Wrap(err, WrapperHubControllerInviteToHubMethod))
		c.render.ErrorResponse(w, http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/hub", http.StatusMovedPermanently)
}
