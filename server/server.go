package server

import (
	"FL_2/cache/tarantoolcache"
	"FL_2/store/postgresstore"
	"log"
	"net/http"
)

func Start(config *Config, https bool) error {
	store  := postgresstore.New(config.DSN)
	cache, err := tarantoolcache.New(config.DatabaseURL)
	if err != nil{
		log.Fatal(err)
	}
	if err = store.Open(); err != nil {
		log.Fatal(err)
	}

	s := newServer(store,cache, config)
	if https {
		return http.ListenAndServeTLS(config.BindAddr,
			"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
			"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem",
			s)
	}

	return http.ListenAndServe(config.BindAddr, s)
}

