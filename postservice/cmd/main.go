package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	pb "post/api"
	"post/configs"
	"post/internal/app/database"
	"post/internal/app/logger"
	"post/internal/app/middleware"

	orderHandlers "post/internal/app/order/handlers"
	orderRepo "post/internal/app/order/repository"
	orderUCase "post/internal/app/order/usecase"

	responseHandlers "post/internal/app/response/handlers"
	responseRepo "post/internal/app/response/repository"
	responseUCase "post/internal/app/response/usecase"

	vacancyHandlers "post/internal/app/vacancy/handlers"
	vacancyRepo "post/internal/app/vacancy/repository"
	vacancyUCase "post/internal/app/vacancy/usecase"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/server.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.InitLogger("stdout"); err != nil {
		log.Fatal(err)
	}

	postgres, err := database.NewPostgres(config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//TODO: поправить порт коннекта
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log .Fatal(err)
	}
	defer conn.Close()
	userRepo := pb.NewUserClient(conn)

	orderRepository := orderRepo.NewRepo(postgres.GetPostgres())
	vacancyRepository := vacancyRepo.NewRepo(postgres.GetPostgres())
	responseRepository := responseRepo.NewRepo(postgres.GetPostgres())

	orderUseCase := orderUCase.NewUseCase(*orderRepository, userRepo)
	vacancyUseCase := vacancyUCase.NewUseCase(*vacancyRepository, userRepo)
	responseUseCase := responseUCase.NewUseCase(*responseRepository, userRepo)

	orderHandler := orderHandlers.NewHandler(*orderUseCase)
	vacancyHandler := vacancyHandlers.NewHandler(*vacancyUseCase)
	responseHandler := responseHandlers.NewHandler(*responseUseCase)

	router := mux.NewRouter()
	router.Use(middleware.LoggingRequest)

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	order := router.PathPrefix("/order").Subrouter()
	order.Use(csrfMiddleware)
	order.HandleFunc("", orderHandler.CreateOrder).Methods(http.MethodPost)
	order.HandleFunc("", orderHandler.GetActualOrder).Methods(http.MethodGet)

	//order.HandleFunc("/{id:[0-9]+}", s.handleChangeOrder).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}", orderHandler.GetOrder).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.CreatePostResponse).Methods(http.MethodPost)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.GetAllPostResponses).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.ChangePostResponse).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.DelPostResponse).Methods(http.MethodDelete)
	order.HandleFunc("/{id:[0-9]+}/select", orderHandler.SelectExecutor).Methods(http.MethodPut)
	order.HandleFunc("/{id:[0-9]+}/select", orderHandler.DeleteExecutor).Methods(http.MethodDelete)
	order.HandleFunc("/profile/{id:[0-9]+}", orderHandler.GetAllUserOrders).Methods(http.MethodGet)

	vacancy := router.PathPrefix("/vacancy").Subrouter()
	vacancy.Use(csrfMiddleware)
	vacancy.HandleFunc("", vacancyHandler.CreateVacancy).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}", vacancyHandler.GetVacancy).Methods(http.MethodGet)
	vacancy.HandleFunc("/{id:[0-9]+}/response", responseHandler.CreatePostResponse).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}/response", responseHandler.GetAllPostResponses).Methods(http.MethodGet)

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
}
