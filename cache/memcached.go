package cache

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func MemcacheConnect(servers ...string) (*memcache.Client, error) {
	sl := new(memcache.ServerList)
	if err := sl.SetServers(servers...); err != nil {
		return nil, fmt.Errorf("unable to resolve memcached servers: %w", err)
	}
	log.Println("Connecting to memcached...")
	memcacheClient := memcache.NewFromSelector(sl)
	if err := memcacheClient.Ping(); err != nil {
		return nil, fmt.Errorf("memcached servers down: %w", err)
	}
	log.Println("Successfuly connected to memcached!")
	return memcacheClient, nil
}
