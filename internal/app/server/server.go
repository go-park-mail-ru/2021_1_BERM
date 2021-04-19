package server

import (
	"ff/internal/app/databases"
	"ff/configs"
	"log"
	"net/http"
)

func Start(config *configs.Config, https bool) error {
	store, err := databases.NewPostgres(config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	cache, err := databases.NewTarantool(config.BindAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := newServer(config, store, cache)
	if https {
		return http.ListenAndServeTLS(config.BindAddr,
			"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
			"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem",
			s)
	}

	return http.ListenAndServe(config.BindAddr, s)
}
