package app

import (
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
)

type PollHandler struct {
	memcClient *memcache.Client
}

func NewPollHandler(memcClient *memcache.Client) *PollHandler {
	return &PollHandler{
		memcClient: memcClient,
	}
}

func (h *PollHandler) GetPoll(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(time.Second * 10)
	w.Write([]byte("get poll!"))
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create poll!"))
}
