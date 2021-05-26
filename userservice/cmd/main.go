package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"net"
	"user/api"
	repository5 "user/internal/app/order/repository"
	repository4 "user/internal/app/review/repository"
	usecase4 "user/internal/app/review/usecase"
	handlers3 "user/internal/app/session/handlers"
	repository3 "user/internal/app/session/repository"
	usecase3 "user/internal/app/session/usecase"
	repository2 "user/internal/app/specialize/repository"
	usecase2 "user/internal/app/specialize/usecase"
	"user/internal/app/user/handlers"
	userRepo "user/internal/app/user/repository"
	"user/internal/app/user/usecase"
	"user/pkg/middleware"

	traceutils "github.com/opentracing-contrib/go-grpc"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"log"
	"net/http"
	"user/configs"
	revHandler "user/internal/app/review/handlers"
	specHandler "user/internal/app/specialize/handler"
	"user/pkg/database/postgresql"
	"user/pkg/logger"
	pMetric "user/pkg/metric"
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

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "UserService",
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

	userRepository := userRepo.NewRepo(postgres.GetPostgres())
	specializeRepository := &repository2.Repository{
		Db: postgres.GetPostgres(),
	}
	reviewRepository := &repository4.Repository{
		Db: postgres.GetPostgres(),
	}

	//connect to auth service
	grpcConnAuth, err := grpc.Dial(
		":8085",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConnAuth.Close()

	//connect to order service
	grpcConnOrder, err := grpc.Dial(
		":8086",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConnOrder.Close()

	client := api.NewSessionClient(grpcConnAuth)
	sessionRepository := repository3.New(client)

	orderClient := api.NewOrderClient(grpcConnOrder)
	orderRepository := repository5.New(orderClient)

	userUseCase := usecase.New(userRepository, specializeRepository, reviewRepository)
	userHandler := handlers.New(userUseCase)

	specializeUseCase := usecase2.New(specializeRepository)
	specializeHandler := specHandler.New(specializeUseCase, userUseCase)

	sessionUseCase := usecase3.New(sessionRepository)
	sessionMiddleWare := handlers3.New(sessionUseCase)

	reviewUseCase := usecase4.New(reviewRepository, userRepository, orderRepository)
	reviewHandler := revHandler.New(reviewUseCase)
	csrfMiddleware := middleware.CSRFMiddleware(config.HTTPS)

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.Use(middleware.LoggingRequest)
	apiRoute.Use(sessionMiddleWare.CheckSession)
	apiRoute.Use(csrfMiddleware)
	apiRoute.HandleFunc("/profile/users", userHandler.GetUsers).Methods(http.MethodGet)
	apiRoute.HandleFunc("/profile/{id:[0-9]+}", userHandler.GetUserInfo).Methods(http.MethodGet)
	apiRoute.HandleFunc("/profile/{id:[0-9]+}", userHandler.ChangeProfile).Methods(http.MethodPatch)
	apiRoute.HandleFunc("/profile/{id:[0-9]+}/specialize", specializeHandler.Create).Methods(http.MethodPost)
	apiRoute.HandleFunc("/profile/{id:[0-9]+}/specialize", specializeHandler.Remove).Methods(http.MethodDelete)

	apiRoute.HandleFunc("/profile/review", reviewHandler.Create).Methods(http.MethodPost)
	apiRoute.HandleFunc("/profile/{id:[0-9]+}/review", reviewHandler.GetAllByUserId).Methods(http.MethodGet)
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
		log.Println("Server start on port", config.BindAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer(grpc.UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))
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
