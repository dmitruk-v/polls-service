package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dmitruk-v/poll-service/schema"
	"github.com/go-chi/chi/v5"
)

type PollHandler struct {
	JsonHandler
	pollCache    schema.PollCache
	pollStorage  schema.PollStorage
	htmlRenderer HTMLRenderer
}

func NewPollHandler(pollCache schema.PollCache, pollStorage schema.PollStorage, htmlRenderer HTMLRenderer) *PollHandler {
	return &PollHandler{
		pollCache:    pollCache,
		pollStorage:  pollStorage,
		htmlRenderer: htmlRenderer,
	}
}

func (h *PollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	surveyIDParam := chi.URLParam(r, "survey_id")
	surveyID, err := strconv.ParseInt(surveyIDParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var poll schema.Poll
	var inCache bool
	var pollNotFoundErr *schema.ErrPollNotFound
	// Extract poll from cache
	poll, err = h.pollCache.GetPoll(surveyID)
	if err != nil {
		if !errors.As(err, &pollNotFoundErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		inCache = false
	}
	if !inCache {
		// If absent in cache, extract from database
		poll, err = h.pollStorage.GetPollByID(r.Context(), surveyID)
		if err != nil {
			if errors.As(err, &pollNotFoundErr) {
				http.Error(w, pollNotFoundErr.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := h.htmlRenderer.ShowTemplate(w, "form.html", poll); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var poll schema.Poll
	if err := h.readJSON(r, &poll); err != nil {
		log.Printf("reading request body from json, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.pollStorage.InsertPoll(r.Context(), poll); err != nil {
		log.Printf("inserting poll to database, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.pollCache.SetPoll(poll); err != nil {
		log.Printf("saving poll to cache, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uri := url.Values{}
	for q, a := range poll.PreSetValues {
		uri.Add(q, a)
	}
	result := fmt.Sprintf("http://localhost:8080/polls/%v?%v", poll.SurveyID, uri.Encode())
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))
}
