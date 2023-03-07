package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dmitruk-v/4service/schema"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PollStorage interface {
	InsertPoll(ctx context.Context, id string, poll schema.Poll) error
	GetPollByID(ctx context.Context, id string) (schema.Poll, error)
}

type PollHandler struct {
	cacheClient *memcache.Client
	pollStorage PollStorage
}

func NewPollHandler(cacheClient *memcache.Client, pollStorage PollStorage) *PollHandler {
	return &PollHandler{
		cacheClient: cacheClient,
		pollStorage: pollStorage,
	}
}

func (h *PollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "poll-id")
	item, err := h.cacheClient.Get(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Got Item: %v\n", item)))
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("reading request body, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := uuid.New().String()

	var poll schema.Poll
	if err := json.Unmarshal(body, &poll); err != nil {
		log.Printf("unmarshal request body to json, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.pollStorage.InsertPoll(r.Context(), id, poll); err != nil {
		log.Printf("inserting poll, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: get expiration from ENV variable
	item := &memcache.Item{Key: id, Value: body, Expiration: 24 * 60 * 60}
	log.Println("---", item.Key)
	if err := h.cacheClient.Set(item); err != nil {
		log.Printf("setting memcached value, with error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	uri := url.Values{}
	uri.Add("питання-1", "відповідь-1")
	uri.Add("питання-2", "відповідь-2")
	uri.Add("питання-3", "відповідь-3")
	result := "http://localhost:8080/polls/12345?" + uri.Encode()
	w.Write([]byte(result))
}
