package web

import (
	"log"
	"net/http"
	"time"

	"github.com/dmitruk-v/4service/schema"
)

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Clients struct {
	CacheClient schema.PollCache
}

type Storages struct {
	PollStorage schema.PollStorage
}

type Server struct {
	cfg      ServerConfig
	clients  Clients
	storages Storages
}

func NewServer(cfg ServerConfig, clients Clients, storages Storages) *Server {
	return &Server{
		cfg:      cfg,
		clients:  clients,
		storages: storages,
	}
}

func (s *Server) Run() error {
	server := http.Server{
		Addr:         s.cfg.Addr,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		Handler:      s.initRoutes(),
	}
	defer server.Close()
	log.Printf("Server listen at %v\n", s.cfg.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
