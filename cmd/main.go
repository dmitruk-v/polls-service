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

func run() error {
	dsn := "postgres://postgres:postgres@localhost:5432/mydb"
	db := postgres.MustConnectWithRetry(dsn, 10*time.Minute)
	defer db.Close(context.Background())

	pollStorage := postgres.NewPollStorage(db)

	memcachedClient, err := memcached.Connect("localhost:11211")
	if err != nil {
		return err
	}
	defer memcachedClient.Close()

	pollCacheClient := memcached.NewPollCache(memcachedClient)

	// TODO: Set Addr from ENV variable
	webServerCfg := web.ServerConfig{
		Addr:         "localhost:8080",
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
