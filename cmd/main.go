package main

import (
	"context"
	"log"
	"time"

	"github.com/dmitruk-v/4service/cache/memcached"
	"github.com/dmitruk-v/4service/db/postgres"
	"github.com/dmitruk-v/4service/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const maxTimeout = 10 * time.Minute

func run() error {

	// Get config from ENV variables
	// ---------------------------------------------
	config := ParseEnv()

	// Init database and storages
	// ---------------------------------------------
	db := postgres.MustConnectWithRetry(config.PosgresDSN, maxTimeout)
	defer db.Close(context.Background())
	postgres.MustSeedString(db, postgres.SeedSQL)

	pollStorage := postgres.NewPollStorage(db)

	// Init cache
	// ---------------------------------------------
	memcachedClient := memcached.MustConnectWithRetry(maxTimeout, config.MemcachedServers...)
	defer memcachedClient.Close()

	pollCacheClient := memcached.NewPollCache(memcachedClient)

	// Init web-server
	// ---------------------------------------------
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

	webServer := web.NewServer(webServerCfg, clients, storages)
	return webServer.Run()
}
