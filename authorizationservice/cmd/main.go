package main

import (
	"authorizationservice/api"
	"authorizationservice/configs"
	profHandler "authorizationservice/internal/profile/handlers"
	"authorizationservice/internal/profile/repository/grpcrepository"
	profileImpl "authorizationservice/internal/profile/usecase/impl"
	"authorizationservice/internal/session/handlers"
	sessHandler "authorizationservice/internal/session/handlers"
	"authorizationservice/internal/session/repository/tarantoolrepository"
	"authorizationservice/internal/session/usecase/impl"
	"authorizationservice/pkg/logger"
	"authorizationservice/pkg/middleware"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"google.golang.org/grpc"
	"net"
	"net/http"
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

	opts := tarantool.Opts{
		User: "guest",
	}
	conn, err := tarantool.Connect(config.DatabaseURL, opts)
	defer conn.Close()
	sessionRepository := tarantoolrepository.New(conn)

	// grpc connect to UserService
	grpcConn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()
	client := api.NewUserClient(grpcConn)
	profileRepository := grpcrepository.New(client)

	sessionUseCase := impl.New(sessionRepository)
	profileUseCase := profileImpl.New(profileRepository)

	sessionHandler := sessHandler.New(sessionUseCase)
	profileHandler := profHandler.New(sessionUseCase, profileUseCase)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Use(middleware.LoggingRequest)
	apiRout := router.PathPrefix("/api").Subrouter()
	apiRout.HandleFunc("/logout", sessionHandler.LogOut).Methods(http.MethodDelete)
	apiRout.HandleFunc("/profile", profileHandler.RegistrationProfile).Methods(http.MethodPost)
	apiRout.HandleFunc("/login", profileHandler.AuthorisationProfile).Methods(http.MethodPost)

	profile := apiRout.PathPrefix("/profile").Subrouter()
	profile.Use(csrfMiddleware)
	profile.HandleFunc("/authorized", sessionHandler.CheckLogin).Methods(http.MethodGet)

	c := middleware.CorsMiddleware(config.Origin)
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	go func() {
		log.Println("HTTP server start on port", config.BindAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer()
	srv := handlers.NewGRPCServer(sessionUseCase)
	api.RegisterSessionServer(s, srv)
	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GRPC server start on port :8085")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
