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
	"authorizationservice/pkg/midlewhare"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
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
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		//FIXME поставить нормального юзера
		User:          "test",
		Pass:          "test",
	}
	conn, err := tarantool.Connect(config.DatabaseURL, opts)
	defer conn.Close()
	sessionRepository := tarantoolrepository.New(conn)

	//FIXME написать адрес
	grpcConn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBlock())
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

	router := mux.NewRouter()
	router.Use(midlewhare.LoggingRequest)
	router.HandleFunc("/profile/authorized", sessionHandler.CheckLogin).Methods(http.MethodGet)
	router.HandleFunc("/logout", sessionHandler.CheckLogin).Methods(http.MethodPut)
	router.HandleFunc("/profile", profileHandler.RegistrationProfile).Methods(http.MethodPost)
	router.HandleFunc("/login", sessionHandler.LogOut).Methods(http.MethodPost)

	c := midlewhare.CorsMiddleware(config.Origin)
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	go func() {
		log.Println("HTTP server start")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer()
	srv := handlers.NewGRPCServer(sessionUseCase)
	api.RegisterSessionServer(s, srv)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server start")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
