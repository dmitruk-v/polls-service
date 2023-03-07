package main

import (
	"context"
	"log"
	"os"
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

	// Init database and storages
	// ---------------------------------------------

	dsn, ok := os.LookupEnv("POSTGRES_DSN")
	if !ok {
		panic("env variable POSTGRES_DSN is not found")
	}
	db := postgres.MustConnectWithRetry(dsn, maxTimeout)
	defer db.Close(context.Background())
	postgres.MustSeedString(db, postgres.SeedSQL)

	pollStorage := postgres.NewPollStorage(db)

	// Init cache
	// ---------------------------------------------

	memcachedAddr, ok := os.LookupEnv("MEMCACHED")
	if !ok {
		panic("env variable MEMCACHED is not found")
	}
	memcachedClient := memcached.MustConnectWithRetry(maxTimeout, memcachedAddr)
	defer memcachedClient.Close()

	pollCacheClient := memcached.NewPollCache(memcachedClient)

	// Init web-server
	// ---------------------------------------------

	// TODO: Set Addr from ENV variable
	webServerCfg := web.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	clients := web.Clients{
		CacheClient: pollCacheClient,
	}

	storages := web.Storages{
		PollStorage: pollStorage,
	}

	webServer := web.NewServer(webServerCfg, clients, storages)
	return webServer.Run()
}
