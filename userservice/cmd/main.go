package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net"
	"user/api"
	"user/internal/app/user/handlers"
	"user/pkg/midlewhare"

	"log"
	"net/http"
	"user/configs"
	specializeRepo "user/internal/app/specialize/repository/postgresql"
	userRepo "user/internal/app/user/repository/postgresql"
	"user/internal/app/user/usecase/impl"
	"user/pkg/database/postgresql"
	"user/pkg/logger"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/toml/server.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err = logger.InitLogger("stdout"); err != nil {
		log.Fatal(err)
	}

	postgres, err := postgresql.NewPostgres(config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	userRepository := &userRepo.Repository{
		Db: postgres.GetPostgres(),
	}
	specializeRepository := &specializeRepo.Repository{
		Db: postgres.GetPostgres(),
	}

	userUseCase := impl.New(userRepository, specializeRepository)

	userHandler := handlers.New(userUseCase)
	router := mux.NewRouter()
	router.Use(midlewhare.LoggingRequest)
	router.HandleFunc("profile/{id:[0-9]+}\"", userHandler.GetUserInfo).Methods(http.MethodGet)
	router.HandleFunc("profile/{id:[0-9]+}\"", userHandler.ChangeProfile).Methods(http.MethodPut)

	c := midlewhare.CorsMiddleware(config.Origin)
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	go func() {
		log.Println("GRPC server start")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer()
	srv := handlers.NewGRPCServer(userUseCase)
	api.RegisterUserServer(s, srv)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server start")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
