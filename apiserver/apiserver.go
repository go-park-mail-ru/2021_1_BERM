package apiserver

import (
	"fl_ru/store/tarantoolstore"
	"log"
	"net/http"
)

func Start(config *Config, https bool) error {
	store, err := tarantoolstore.New(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	s := newServer(store, config)
	if https {
		return http.ListenAndServeTLS(config.BindAddr,
			"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
			"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem",
			s)
	}

	return http.ListenAndServe(config.BindAddr, s)
}
