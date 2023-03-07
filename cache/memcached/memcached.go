package memcached

import (
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func Connect(servers ...string) (*memcache.Client, error) {
	sl := new(memcache.ServerList)
	if err := sl.SetServers(servers...); err != nil {
		return nil, fmt.Errorf("unable to resolve MEMCACHED servers: %w", err)
	}
	log.Println("Connecting to MEMCACHED...")
	memcacheClient := memcache.NewFromSelector(sl)
	if err := memcacheClient.Ping(); err != nil {
		return nil, fmt.Errorf("MEMCACHED servers down: %w", err)
	}
	log.Println("Successfuly connected to MEMCACHED!")
	return memcacheClient, nil
}

func MustConnectWithRetry(timeout time.Duration, servers ...string) *memcache.Client {
	timer := time.NewTimer(timeout)
	delay := time.Second
	for {
		select {
		case <-timer.C:
			msg := fmt.Sprintf("cant connect to MEMCACHED, timeout expired: %v", timeout)
			panic(msg)
		default:
			conn, err := Connect(servers...)
			if err == nil {
				return conn
			}
			log.Println(err)
		}
		log.Printf("Retry after %v...\n", delay)
		time.Sleep(delay)
		delay *= 2
	}
}
