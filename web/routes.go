package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) initRoutes() http.Handler {
	router := chi.NewRouter()
	pollHandler := NewPollHandler(s.clients.PollCache, s.storages.PollStorage)
	router.Get("/polls/{survey_id}", pollHandler.GetPoll)
	router.Post("/polls", pollHandler.CreatePoll)
	return router
}
