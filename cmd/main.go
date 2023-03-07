package main

import (
	"log"
	"time"

	"github.com/dmitruk-v/4service/app"
	"github.com/dmitruk-v/4service/cache"
	"github.com/dmitruk-v/4service/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dsn := "postgres://postgres:mysecretpassword@localhost:5432/polls-db"
	db := postgres.MustConnectWithRetry(dsn, 10*time.Minute)
	pollStorage := postgres.NewPollStorage(db)

	cacheClient, err := cache.MemcacheConnect("localhost:11211")
	if err != nil {
		return err
	}
	defer cacheClient.Close()

	cfg := app.AppConfig{}
	cfg.HTTPServer.Addr = "localhost:8080"
	cfg.HTTPServer.ReadTimeout = 5 * time.Second
	cfg.HTTPServer.WriteTimeout = 5 * time.Second

	app := app.NewApp(cfg, cacheClient, pollStorage)
	return app.Run()
}
