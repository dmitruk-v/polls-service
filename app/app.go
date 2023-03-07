package app

import (
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
)

type AppConfig struct {
	HTTPServer struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
}

type App struct {
	cfg         AppConfig
	cacheClient *memcache.Client
	pollStorage PollStorage
}

func NewApp(cfg AppConfig, cacheClient *memcache.Client, pollStorage PollStorage) *App {
	return &App{
		cfg:         cfg,
		cacheClient: cacheClient,
		pollStorage: pollStorage,
	}
}

func (app *App) Run() error {
	server := http.Server{
		Addr:         app.cfg.HTTPServer.Addr,
		ReadTimeout:  app.cfg.HTTPServer.ReadTimeout,
		WriteTimeout: app.cfg.HTTPServer.WriteTimeout,
		Handler:      app.initRoutes(),
	}
	defer server.Close()
	log.Printf("Server listen at %v\n", app.cfg.HTTPServer.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (app *App) initRoutes() http.Handler {
	router := chi.NewRouter()

	pollHandler := NewPollHandler(app.cacheClient, app.pollStorage)
	router.Get("/polls/{poll-id}", pollHandler.GetPoll)
	router.Post("/polls", pollHandler.CreatePoll)

	return router
}
