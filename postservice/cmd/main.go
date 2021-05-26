package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"post/api"
	pb "post/api"
	"post/configs"
	"post/internal/app/session/handlers"
	"post/internal/app/session/repository"
	"post/internal/app/session/usecase"
	"post/pkg/logger"
	"post/pkg/middleware"
	"post/pkg/postgresql"

	pMetric "post/pkg/metric"

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
	flag.StringVar(&configPath, "config-path", "configs/toml/server.toml", "path to config file")
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

	postgres, err := postgresql.NewPostgres(config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "PostService",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// connect to user service
	conn, err := grpc.Dial(
		config.Host + ":8081",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	userRepo := pb.NewUserClient(conn)

	// connect to auth service
	grpcConn, err := grpc.Dial(
		config.Host + ":8085",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()
	client := api.NewSessionClient(grpcConn)
	sessionRepository := repository.New(client)
	sessionUseCase := usecase.New(sessionRepository)
	sessionMiddleWare := handlers.New(sessionUseCase)

	orderRepository := orderRepo.NewRepo(postgres.GetPostgres())
	vacancyRepository := vacancyRepo.NewRepo(postgres.GetPostgres())
	responseRepository := responseRepo.NewRepo(postgres.GetPostgres())

	orderUseCase := orderUCase.NewUseCase(orderRepository, userRepo)
	vacancyUseCase := vacancyUCase.NewUseCase(vacancyRepository, userRepo)
	responseUseCase := responseUCase.NewUseCase(responseRepository, userRepo)

	orderHandler := orderHandlers.NewHandler(orderUseCase)
	vacancyHandler := vacancyHandlers.NewHandler(vacancyUseCase)
	responseHandler := responseHandlers.NewHandler(responseUseCase)

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.Use(middleware.LoggingRequest)
	apiRoute.Use(sessionMiddleWare.CheckSession)

	order := apiRoute.PathPrefix("/order").Subrouter()
	order.Use(csrfMiddleware)
	order.HandleFunc("", orderHandler.CreateOrder).Methods(http.MethodPost)
	order.HandleFunc("", orderHandler.GetActualOrder).Methods(http.MethodGet)

	order.HandleFunc("/{id:[0-9]+}", orderHandler.GetOrder).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}", orderHandler.ChangeOrder).Methods(http.MethodPatch)
	order.HandleFunc("/{id:[0-9]+}", orderHandler.DeleteOrder).Methods(http.MethodDelete)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.CreatePostResponse).Methods(http.MethodPost)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.GetAllPostResponses).Methods(http.MethodGet)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.ChangePostResponse).Methods(http.MethodPatch)
	order.HandleFunc("/{id:[0-9]+}/response", responseHandler.DelPostResponse).Methods(http.MethodDelete)
	order.HandleFunc("/{id:[0-9]+}/select", orderHandler.SelectExecutor).Methods(http.MethodPost)
	order.HandleFunc("/{id:[0-9]+}/select", orderHandler.DeleteExecutor).Methods(http.MethodDelete)
	order.HandleFunc("/profile/{id:[0-9]+}", orderHandler.GetAllUserOrders).Methods(http.MethodGet)
	order.HandleFunc("/{id}/close", orderHandler.CloseOrder).Methods(http.MethodDelete)
	order.HandleFunc("/profile/{id:[0-9]+}/archive", orderHandler.GetAllArchiveUserOrders).Methods(http.MethodGet)
	order.HandleFunc("/search", orderHandler.SearchOrder).Methods(http.MethodPatch)
	order.HandleFunc("/suggest", orderHandler.SuggestOrderTitle).Methods(http.MethodGet)

	vacancy := apiRoute.PathPrefix("/vacancy").Subrouter()
	vacancy.Use(csrfMiddleware)
	vacancy.HandleFunc("", vacancyHandler.CreateVacancy).Methods(http.MethodPost)
	vacancy.HandleFunc("", vacancyHandler.GetActualVacancies).Methods(http.MethodGet)
	vacancy.HandleFunc("/{id:[0-9]+}", vacancyHandler.GetVacancy).Methods(http.MethodGet)
	vacancy.HandleFunc("/{id:[0-9]+}", vacancyHandler.ChangeVacancy).Methods(http.MethodPatch)
	vacancy.HandleFunc("/{id:[0-9]+}", vacancyHandler.DeleteVacancy).Methods(http.MethodDelete)
	vacancy.HandleFunc("/{id:[0-9]+}/response", responseHandler.CreatePostResponse).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}/response", responseHandler.GetAllPostResponses).Methods(http.MethodGet)
	vacancy.HandleFunc("/{id:[0-9]+}/response", responseHandler.ChangePostResponse).Methods(http.MethodPatch)
	vacancy.HandleFunc("/{id:[0-9]+}/response", responseHandler.DelPostResponse).Methods(http.MethodDelete)
	vacancy.HandleFunc("/profile/{id:[0-9]+}", vacancyHandler.GetAllUserVacancies).Methods(http.MethodGet)
	vacancy.HandleFunc("/{id:[0-9]+}/select", vacancyHandler.SelectExecutor).Methods(http.MethodPost)
	vacancy.HandleFunc("/{id:[0-9]+}/select", vacancyHandler.DeleteExecutor).Methods(http.MethodDelete)
	vacancy.HandleFunc("/{id}/close", vacancyHandler.CloseVacancy).Methods(http.MethodDelete)
	vacancy.HandleFunc("/profile/{id:[0-9]+}/archive", vacancyHandler.GetAllArchiveUserVacancies).Methods(http.MethodGet)
	vacancy.HandleFunc("/search", vacancyHandler.SearchVacancy).Methods(http.MethodPatch)

	c := middleware.CorsMiddleware(config.Origin)

	pMetric.New()
	server := &http.Server{
		Addr:    config.BindAddr,
		Handler: c.Handler(router),
	}

	go func() {
		if config.HTTPS {
			log.Println("TLS server starting at port: ", server.Addr)
			if err := server.ListenAndServeTLS(
				"/etc/letsencrypt/live/findfreelancer.ru/cert.pem",
				"/etc/letsencrypt/live/findfreelancer.ru/privkey.pem"); err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Server starting at port", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer(grpc.UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))
	srv := orderHandlers.NewGRPCServer(orderUseCase)
	api.RegisterOrderServer(s, srv)

	l, err := net.Listen("tcp", ":8086")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GRPC Server start on port :8086")
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
