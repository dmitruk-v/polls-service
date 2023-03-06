package main

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dmitruk-v/4service/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	memcacheClient := memcache.New()

	cfg := app.AppConfig{}
	cfg.HTTPServer.Addr = ":8080"
	cfg.HTTPServer.ReadTimeout = 5 * time.Second
	cfg.HTTPServer.WriteTimeout = 5 * time.Second

	app := app.NewApp(cfg, memcacheClient)
	return app.Run()
}
