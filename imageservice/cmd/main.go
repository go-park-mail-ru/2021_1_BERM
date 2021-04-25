package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"image/api"
	"image/configs"
	imageHandlers "image/internal/app/image/handlers"
	imageUCase "image/internal/app/image/usecase"
	"image/internal/app/logger"
	"image/internal/app/middleware"
	"log"
	"net/http"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/server.toml", "path to config file")
}

func main() {
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.InitLogger("stdout"); err != nil {
		log.Fatal(err)
	}

	//TODO: поправить порт коннекта
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log .Fatal(err)
	}
	defer conn.Close()
	userRepo := api.NewUserClient(conn)

	imageUseCase := imageUCase.NewUseCase(userRepo)

	imageHandler := imageHandlers.NewHandler(*imageUseCase)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Use(middleware.LoggingRequest)
	router.Use(csrfMiddleware)

	router.HandleFunc("/profile/avatar", imageHandler.PutAvatar).Methods(http.MethodPut)

	c := middleware.CorsMiddleware(config.Origin)

	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	if config.HTTPS {
		if err := server.ListenAndServeTLS(
			"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
			"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem"); err != nil {
			log.Fatal(err)
		}
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	//s := grpc.NewServer()
	//srv := &imageserver.ImageServer{}
	//api.RegisterImageServer(s, srv)
	//
	//l, err := net.Listen("tcp", ":8080")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := s.Serve(l); err != nil {
	//	log.Fatal(err)
	//}
}