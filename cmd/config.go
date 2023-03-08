package main

import (
	"os"
	"strings"
)

type EnvConfig struct {
	PosgresDSN       string
	MemcachedServers []string
}

func ParseEnv() EnvConfig {
	cfg := EnvConfig{}
	var ok bool
	cfg.PosgresDSN, ok = os.LookupEnv("POSTGRES_DSN")
	if !ok {
		panic("env variable POSTGRES_DSN is not found")
	}
	memcachedServers, ok := os.LookupEnv("MEMCACHED")
	if !ok {
		panic("env variable MEMCACHED is not found")
	}
	cfg.MemcachedServers = strings.Split(strings.TrimSpace(memcachedServers), ",")
	return cfg
}
