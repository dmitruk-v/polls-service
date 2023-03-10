package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitruk-v/poll-service/cache/memcached"
	"github.com/dmitruk-v/poll-service/db/postgres"
	"github.com/dmitruk-v/poll-service/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const maxTimeout = 10 * time.Minute

func run() error {
	// Get config from ENV variables
	config := ParseEnv()

	// Init database and storages
	db := postgres.MustConnectWithRetry(config.PosgresDSN, maxTimeout)
	defer db.Close(context.Background())
	postgres.MustSeedString(db, postgres.SeedSQL)

	pollStorage := postgres.NewPollStorage(db)

	// Init cache
	memcachedClient := memcached.MustConnectWithRetry(config.MemcachedServers, maxTimeout)
	defer memcachedClient.Close()

	pollCacheClient := memcached.NewPollCache(memcachedClient)

	// Init exit context
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-exitCh
		cancel()
	}()

	// Init static-server
	go func() {
		if err := web.RunStaticServer(ctx, ":8081"); err != nil {
			log.Println(err)
		}
	}()

	// Init HTML-renderer
	htmlRenderer := web.NewBaseHTMLRender()

	// Init web-server
	webServerCfg := web.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	clients := web.Clients{
		PollCache: pollCacheClient,
	}
	storages := web.Storages{
		PollStorage: pollStorage,
	}
	webServer := web.NewServer(webServerCfg, clients, storages, htmlRenderer)
	return webServer.Run(ctx)
}
