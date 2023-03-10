package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dmitruk-v/poll-service/schema"
)

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Clients struct {
	PollCache schema.PollCache
}

type Storages struct {
	PollStorage schema.PollStorage
}

type Server struct {
	cfg          ServerConfig
	clients      Clients
	storages     Storages
	htmlRenderer HTMLRenderer
}

func NewServer(cfg ServerConfig, clients Clients, storages Storages, htmlRenderer HTMLRenderer) *Server {
	return &Server{
		cfg:          cfg,
		clients:      clients,
		storages:     storages,
		htmlRenderer: htmlRenderer,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := http.Server{
		Addr:         s.cfg.Addr,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		Handler:      s.initRoutes(),
	}
	go func() {
		<-ctx.Done()
		log.Println("Closing backend http server...")
		if err := server.Close(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("Server listen at %v\n", s.cfg.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
