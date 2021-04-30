package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net"
	"user/api"
	impl2 "user/internal/app/specialize/usecase/impl"
	"user/internal/app/user/handlers"
	handlers2 "user/internal/session/handlers"
	"user/internal/session/repository/grpcrepository"
	impl3 "user/internal/session/usecase/impl"
	"user/pkg/middleware"

	"log"
	"net/http"
	"user/configs"
	specHandler "user/internal/app/specialize/handler"
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

	//connect to auth service
	grpcConn, err := grpc.Dial(":8085", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()
	client := api.NewSessionClient(grpcConn)
	sessionRepository := grpcrepository.New(client)

	userUseCase := impl.New(userRepository, specializeRepository)
	userHandler := handlers.New(userUseCase)

	specializeUseCase := impl2.New(specializeRepository)
	specializeHandler := specHandler.New(specializeUseCase, userUseCase)

	sessionUseCase := impl3.New(sessionRepository)
	sessionMiddleWare := handlers2.New(sessionUseCase)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Use(sessionMiddleWare.CheckSession)
	router.Use(middleware.LoggingRequest)
	router.Use(csrfMiddleware)
	router.HandleFunc("/profile/{id:[0-9]+}", userHandler.GetUserInfo).Methods(http.MethodGet)
	router.HandleFunc("/profile/{id:[0-9]+}", userHandler.ChangeProfile).Methods(http.MethodPut)
	router.HandleFunc("/profile/{id:[0-9]+}/specialize", specializeHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/profile/{id:[0-9]+}/specialize", specializeHandler.Remove).Methods(http.MethodDelete)

	c := middleware.CorsMiddleware(config.Origin)
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	go func() {
		log.Println("Server start on port", config.BindAddr)
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

	log.Println("GRPC Server start on port :8081")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
