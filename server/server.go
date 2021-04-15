package server

import (
	"FL_2/store/diskmediastore"
	"FL_2/store/postgresstore"
	"FL_2/store/tarantoolcache"
	"FL_2/usecase/implementation"
	"log"
	"net/http"
)

func Start(config *Config, https bool) error {
	store := postgresstore.New(config.DSN)
	cache, err := tarantoolcache.New(config.DatabaseURL)
	mediaStore := diskmediastore.New(config.ContentDir)
	useCase := implementation.New(store, cache, mediaStore)
	if err != nil {
		log.Fatal(err)
	}
	if err = store.Open(); err != nil {
		log.Fatal(err)
	}

	s := newServer(useCase, config)
	if https {
		return http.ListenAndServeTLS(config.BindAddr,
			"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
			"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem",
			s)
	}

	return http.ListenAndServe(config.BindAddr, s)
}
