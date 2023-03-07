package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dmitruk-v/4service/schema"
	"github.com/go-chi/chi/v5"
)

type PollHandler struct {
	JsonHandler
	pollCache   schema.PollCache
	pollStorage schema.PollStorage
}

func NewPollHandler(pollCache schema.PollCache, pollStorage schema.PollStorage) *PollHandler {
	return &PollHandler{
		pollCache:   pollCache,
		pollStorage: pollStorage,
	}
}

func (h *PollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "survey_id")
	poll, err := h.pollCache.GetPoll(surveyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Got poll: %v\n", poll)))
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var poll schema.Poll
	if err := h.readJSON(r, &poll); err != nil {
		log.Printf("reading request body from json, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.pollStorage.InsertPoll(r.Context(), poll); err != nil {
		log.Printf("inserting poll, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uri := url.Values{}
	uri.Add("survey_id", strconv.Itoa(int(poll.SurveyID)))
	for k, v := range poll.PreSetValues {
		uri.Add(k, v)
	}
	result := "http://localhost:8080/polls?" + uri.Encode()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))
}