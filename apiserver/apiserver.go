package apiserver

import (
	"fl_ru/store/tarantoolstore"
	"log"
	"net/http"
)

func Start(config *Config) error{
	store, err := tarantoolstore.New(config.DatabaseUrl)
	if err != nil{
		log.Fatal(err)
	}
	s:=newServer(store)
	return http.ListenAndServe(config.BindAddr, s)
}

