package driver

import (
	"fmt"
	"net/http"

	"github.com/ghazlabs/idn-remote-scheduler/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"
)

type API struct {
	APIConfig
}

type APIConfig struct {
	Service        core.Service `validate:"nonnil"`
	ClientUsername string       `validate:"nonzero"`
	ClientPassword string       `validate:"nonzero"`
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid API config: %w", err)
	}
	return &API{APIConfig: cfg}, nil
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(BasicAuth(a.ClientUsername, a.ClientPassword))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/health", a.serveHealthCheck)
	r.Get("/messages", a.serveGetMessages)
	r.Post("/messages", a.serveScheduleMessage)
	r.Post("/messages/{id}/retry", a.serveRetryMessage)

	return r
}

func (a *API) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's working!"))
}

func (a *API) serveGetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := a.Service.GetAllMessages(r.Context())
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}

	resp := NewSuccessResp(messages)
	render.Render(w, r, resp)
}

func (a *API) serveScheduleMessage(w http.ResponseWriter, r *http.Request) {
	resp := NewSuccessResp(nil)
	render.Render(w, r, resp)
}

func (a *API) serveRetryMessage(w http.ResponseWriter, r *http.Request) {
	resp := NewSuccessResp(nil)
	render.Render(w, r, resp)
}
